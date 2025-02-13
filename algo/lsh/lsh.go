package algo

// References:
// https://github.com/pinecone-io/examples/blob/master/learn/search/faiss-ebook/locality-sensitive-hashing-traditional/sparse_implementation.ipynb
// https://www.pinecone.io/learn/series/faiss/locality-sensitive-hashing/

import (
	"fmt"
	"gowtools/algo"
	"math/rand/v2"
	"time"
)

var Verbose = false

type hashVal = int

type LSH struct {
	Vocab     []string
	HashFuncs [][]int
	Entries   []LshEntry

	nbBands           int
	shingleWindowSize int
}

func (l LSH) Find(key string, limit int) []any {
	shingles := shingle(l.shingleWindowSize, key)
	searchSignature := getHashSignature(shingles, l.HashFuncs, l.Vocab)
	searchBands := splitHashSignatureIntoBands(l.nbBands, searchSignature)

	result := []any{}

	// first, evaluate candidates by matching a pair of single band
	candidates := []*LshEntry{}
	for i := range l.Entries {
		e := l.Entries[i]
		for _, band := range e.SignatureBands {
			// stop processing the entry as soon as we have 1 band match
			found := false
			for _, searchBand := range searchBands {
				if isEqual(searchBand, band) {
					found = true
					break
				}
			}
			if found {
				candidates = append(candidates, &l.Entries[i])

				// early exit if we can :)
				if len(result) >= limit {
					return result
				}

				// otherwise breal frag bands, go to next frag iteration
				break
			}
		}
	}
	log(fmt.Sprintf("Found %d candidates", len(candidates)))

	// then, run jacard similarity for each entry
	results := []any{}
	for i, c := range candidates {
		similarity := algo.Jaccard(c.Singature, searchSignature)
		if similarity >= 0.6 {
			results = append(results, candidates[i].Value)
		}
	}
	log(fmt.Sprintf("Found %d results with good hash similarity", len(candidates)))

	return result
}

func isEqual(a []int, b []int) bool {
	max := len(b)
	if len(a) > len(b) {
		max = len(a)
	}

	for i := 0; i < max; i++ {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

type LshEntry struct {
	OriginalKey    string
	Shingles       map[string]uint8 // nil most of the time, dicarded after processing
	Singature      []int
	SignatureBands [][]int
	Value          any
}

type KeyValue struct {
	Key   string
	Value any
}

// signatureLength is the hash size. The bigger, the more precision each entry will have
//
// nBands is the number of subvectors that will be generated for each hash.
// signatureLength must be divisable by this size otherwise the code will panic.
// panic will also occur if nBands is under 1
// Searching will be considered a hit if a subvector of input matches at least 1 band, for efficiency.
// Increase nb of bands means less precision is required to match, but vastly widens the search (increases false positives)
// Decrease nb of bands means more precision, but thinner search (few words won't do it)
//
// shingleWindow size determines the size the raw values observed to build the global vocabularity.  Increasing shingle vastly improves uniqueness of hashes, but is costly for indexing time and reduces fuzzyness
//
// data is the data which defines the key to hash along its value (or pointer value preferable) to store in indexes
func BuildLSH(signatureLength int, nBands int, shingleWindowSize int, data []KeyValue) LSH {
	if signatureLength%nBands != 0 {
		panic("lsh: signature length must be divisible by nb of bands")
	}
	if nBands < 1 {
		panic("lsh: nBand must be at least 1")
	}

	start := time.Now()

	entries := make([]LshEntry, len(data)) // the index entries (hashed vals + actual values)
	vocabMap := map[string]uint8{}         // vocab holds all the unique shingles

	for i := range data {
		d := data[i]

		shingles := shingle(int(shingleWindowSize), d.Key)

		// add shingling to global vocab
		for s := range shingles {
			vocabMap[s] = 0
		}

		// create entry in the index
		entries[i] = LshEntry{
			OriginalKey:    d.Key,
			Shingles:       shingles,
			Singature:      nil,
			SignatureBands: nil,
			Value:          d.Value,
		}
	}

	vectorSize := len(vocabMap)

	// we must have vocab addressable by index / position for vectors
	vocabSlc := make([]string, vectorSize)
	i := 0
	for shingle, _ := range vocabMap {
		vocabSlc[i] = shingle
		i++
	}
	vocabMap = nil // free memory (or at least let GC know)

	log(fmt.Sprintf("vocab size is %d", len(vocabSlc)))

	// prepare the hash functions
	// Each hash function is ran based on the signature / hash length, with a randomized slice of vocab positions
	// See other comment in function below for more explanations

	hashFuncs := make([][]int, signatureLength)
	for i := range hashFuncs {
		hashFuncs[i] = getNewHashVectorRandomized(vocabSlc)
	}
	log(fmt.Sprintf("Prepared all random hash funcs for signature length of %d", signatureLength))

	log(fmt.Sprintf("Hashing %d elements... This can take some time", len(entries)))
	for i := range entries {
		e := &entries[i]
		signature := getHashSignature(e.Shingles, hashFuncs, vocabSlc)

		// Lastly, we create subvectors (nbBands)
		// This lowers accuracy but vastly improves comparison speed
		bands := splitHashSignatureIntoBands(nBands, signature)

		e.SignatureBands = bands
	}

	log(fmt.Sprintf("loaded lsh index in %s", time.Since(start).String()))

	return LSH{
		Vocab:             vocabSlc,
		HashFuncs:         hashFuncs,
		Entries:           entries,
		nbBands:           nBands,
		shingleWindowSize: shingleWindowSize,
	}
}

func log(msg string) {
	if Verbose {
		fmt.Println("[LSH index]: " + msg)
	}
}

// shingle moves a sliding window of size k across val and returns all unique values (shingles)
// Values of the map can be discarded
func shingle(k int, val string) map[string]uint8 {
	// Move a sliding window of size k to gather all values
	shingles := make(map[string]uint8)
	for i := range len(val) - k + 1 {
		subslice := val[i : i+k]
		shingles[subslice] = 0
	}

	return shingles
}

// getNewHashVectorRandomized creates a randomized vector, whose values contain every possible position / index in vocab. But 1-indexed
// A hash function / vector is meant to be used to determine a single value in a signature
func getNewHashVectorRandomized(vocabSlc []string) []int {
	shuffledHashValues := make([]int, len(vocabSlc))
	for idxVocab := range vocabSlc {
		shuffledHashValues[idxVocab] = idxVocab + 1
	}

	rand.Shuffle(len(shuffledHashValues), func(i, j int) {
		iVal := shuffledHashValues[i]
		shuffledHashValues[i] = shuffledHashValues[j]
		shuffledHashValues[j] = iVal
	})

	return shuffledHashValues
}

func getHashSignature(entryShingles map[string]uint8, hashFuncs [][]int, vocab []string) []int {
	// once we have all shingle, we create a sparse (wide/long) vectors
	// which is filled with 0 and set to 1 when a shingling is present in the global vocab
	// This is also known as "one-hot encoding"
	sparseVector := make([]uint8, len(vocab))
	for i, vocabVal := range vocab {
		if _, ok := entryShingles[vocabVal]; ok {
			sparseVector[i] = 1
		}
	}

	// Then, we compress the sparse vector by MinHashing
	// Second, we iterate each value (some random vocab shignling) of the randomized hash vector, check if the current vector has it, if not loop until we hit a random value that is in our entry's shingles
	// We repeat this n amount of times (signatureLength) to build a minhash signature, aka our dense vector
	signature := make([]int, len(hashFuncs))
	for idxSignature, hashFuncVals := range hashFuncs {
		for _, vectorPosValue := range hashFuncVals {
			isValueInVector := sparseVector[vectorPosValue-1]
			if isValueInVector == 1 {
				signature[idxSignature] = vectorPosValue // add the 1-indexed value, not 0 indexed
				break
			}
		}
	}

	return signature
}

func splitHashSignatureIntoBands(nbBands int, fullSignature []int) [][]int {
	bandSize := len(fullSignature) / nbBands // a validation should be done before this

	bands := make([][]int, nbBands)
	for i := range nbBands {
		start := i * bandSize
		end := start + bandSize // always excluded bound
		bands[i] = fullSignature[start:end]
	}

	return bands
}
