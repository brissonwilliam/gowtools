package lsh

// Locality Sensitive Hashing
// This algorithm is really cool. You can throw anything at it. It's used to compare DNA sequeneces, reverse image search, text search etc.
// It works by creating a big vector of the whole dataset (vocabulary), then hashing each entry with some randomness to a normalized / fixed-sized vectors.
// Hashes are split in bands, that are stored in bucket for faster lookup.
// On lookup, the input is hashed and vector similarity is ran against potential matches
// This specific implementation is more oriented towards text search.
//
// TODO: implement generics for storing concrete values instead of "any"
//
// References:
// https://github.com/pinecone-io/examples/blob/master/learn/search/faiss-ebook/locality-sensitive-hashing-traditional/sparse_implementation.ipynb
// https://www.pinecone.io/learn/series/faiss/locality-sensitive-hashing/

import (
	"fmt"
	"gowtools/algo"
	"math"
	"math/rand/v2"
	"sort"
	"strconv"
	"time"
)

var Verbose = false

// Value holds the real data to index and retreive
type Value interface {
	GetID() any
}

// TODO: allow different hash sizes
type hashVal = uint32

const maxHashVal = math.MaxUint32

type LSH struct {
	Vocab     []string
	HashFuncs [][]hashVal
	Entries   []LshEntry
	Buckets   []LSHBucket

	signatureLength   int
	nbBands           int
	shingleWindowSize int
}

func isEqual(a []hashVal, b []hashVal) bool {
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
	OriginalKey string
	Shingles    map[string]uint8 // nil most of the time, dicarded after processing
	Singature   []hashVal
	Value       Value
}

type KeyValue struct {
	Key   string
	Value Value
}

type LSHBucket struct {
	Bands map[string]LSHBucketBand
}

type LSHBucketBand struct {
	Band     []hashVal
	Elements []*LshEntry
}

// signatureLength is the hash size. The bigger, the more precision each entry will have
//
// nBands is the number of subvectors that will be generated for each hash.
// signatureLength must be divisable by this size otherwise the code will panic.
// panic will also occur if nBands is under 1
//
// shingleWindow size determines the size the raw values observed to build the global vocabularity. Increasing shingle vastly improves uniqueness of values (increased sparseness), but is costly for indexing time and reduces fuzzyness
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
			OriginalKey: d.Key,
			Shingles:    shingles,
			Singature:   nil,
			Value:       d.Value,
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

	hashFuncs := make([][]hashVal, nBands)
	for i := range hashFuncs {
		hashFuncs[i] = getNewHashVectorRandomized(vocabSlc)
	}
	log(fmt.Sprintf("Prepared %d random hash funcs for signature length of %d", len(hashFuncs), signatureLength))

	// prepare band buckets
	// Each band bucket is to increase search speed and allow not having to iterate
	// and compare all the data against an input, which can be expensive
	// An input will be hashed on search, and we will try to only look into each bucket if there are entries to compare (candidates)
	// This is the "locality" part of the algorithm
	buckets := make([]LSHBucket, nBands)

	log(fmt.Sprintf("Hashing %d elements... This can take some time", len(entries)))
	for i := range entries {
		e := &entries[i]
		signature := getHashSignature(e.Shingles, signatureLength, hashFuncs, vocabSlc)
		e.Shingles = nil // shignles not needed anymore, free some memory
		e.Singature = signature

		// Lastly, we create subvectors (nbBands)
		// And assign it to the right bucket for increased search speed
		bands := splitHashSignatureIntoSubvectors(nBands, signature)
		if len(bands) != len(buckets) {
			panic("[lsh] signature nb of bands does not match nb of buckets allocated")
		}
		for i := range bands {
			bandHash := hashBandForBucketAccess(bands[i])
			// Check if the band already exists. If so append
			// If not, create it

			if buckets[i].Bands == nil {
				buckets[i].Bands = map[string]LSHBucketBand{
					bandHash: {
						Band:     bands[i],
						Elements: []*LshEntry{e},
					},
				}
			}

			bucketBand, bandExistsInBucket := buckets[i].Bands[bandHash]
			if bandExistsInBucket {
				bucketBand.Elements = append(bucketBand.Elements, e)
				buckets[i].Bands[bandHash] = bucketBand
			} else {
				buckets[i].Bands[bandHash] = LSHBucketBand{
					Band:     bands[i],
					Elements: []*LshEntry{e},
				}

			}
		}
	}
	if Verbose {
		totalBucketElements := 0
		for i := range buckets {
			totalBucketElements += len(buckets[i].Bands)
		}
		avgBucketSize := totalBucketElements / len(buckets)
		log(fmt.Sprintf("made %d buckets of avg size %d", len(buckets), avgBucketSize))
	}

	log(fmt.Sprintf("loaded lsh index in %s", time.Since(start).String()))

	return LSH{
		Vocab:             vocabSlc,
		HashFuncs:         hashFuncs,
		Entries:           entries,
		Buckets:           buckets,
		signatureLength:   signatureLength,
		nbBands:           nBands,
		shingleWindowSize: shingleWindowSize,
	}
}

func hashBandForBucketAccess(band []hashVal) string {
	pseudoHash := ""
	for _, bandVal := range band {
		pseudoHash += strconv.Itoa(int(bandVal)) + "-"
	}
	return pseudoHash
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
func getNewHashVectorRandomized(vocabSlc []string) []hashVal {
	shuffledHashValues := make([]hashVal, len(vocabSlc))
	for idxVocab := range vocabSlc {
		if idxVocab+1 >= maxHashVal {
			panic("cannot assign hash value: vocab rand position index exceeds max allowed hash value. Consider reducing vocab, or changing hashVal type")
		}
		shuffledHashValues[idxVocab] = hashVal(idxVocab + 1)
	}

	rand.Shuffle(len(shuffledHashValues), func(i, j int) {
		iVal := shuffledHashValues[i]
		shuffledHashValues[i] = shuffledHashValues[j]
		shuffledHashValues[j] = iVal
	})

	return shuffledHashValues
}

func getHashSignature(entryShingles map[string]uint8, signatureLength int, hashFuncs [][]hashVal, vocab []string) []hashVal {
	if signatureLength%len(hashFuncs) != 0 {
		panic("signature length and nBands (hashfuncs) must be divisable")
	}
	bandLength := signatureLength / len(hashFuncs)

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
	signature := make([]hashVal, signatureLength)
	idxHashFunc := 0
	for i := 0; i < signatureLength; i += bandLength {
		hashFunc := hashFuncs[idxHashFunc]
		for j := range bandLength {
			// find the first matching element in hashFunc values (which are random positions, 1-indexed)
			for _, randomPosHashVal := range hashFunc {
				isValueInVector := sparseVector[randomPosHashVal-1]
				if isValueInVector == 1 {
					signature[i+j] = randomPosHashVal // add the 1-indexed value, not 0 indexed
					break
				}
			}
		}
		idxHashFunc += 1
	}

	return signature
}

func splitHashSignatureIntoSubvectors(nbBands int, fullSignature []hashVal) [][]hashVal {
	bandSize := len(fullSignature) / nbBands // a validation should be done before this

	bands := make([][]hashVal, nbBands)
	for i := range nbBands {
		start := i * bandSize
		end := start + bandSize // always excluded bound
		bands[i] = fullSignature[start:end]
	}

	return bands
}

type LSHResult struct {
	Score float64
	Value any
}

func (l LSH) Find(key string, hashSimilarity float64) []LSHResult {
	shingles := shingle(l.shingleWindowSize, key)
	fmt.Printf("search shingles: %#v\n", shingles)

	searchSignature := getHashSignature(shingles, l.signatureLength, l.HashFuncs, l.Vocab)

	searchSignatureMap := make(map[hashVal]uint8, len(searchSignature))
	for _, hashVal := range searchSignature {
		searchSignatureMap[hashVal] = 1
	}

	searchBands := splitHashSignatureIntoSubvectors(l.nbBands, searchSignature)

	// first, evaluate candidates by looking into buckets if we have a match
	// to not have to compare against entire data set
	candidatesDeduped := map[any]*LshEntry{}
	bucketMatchCount := 0
	for i, searchBand := range searchBands {
		bucket := l.Buckets[i]

		searchBandHash := hashBandForBucketAccess(searchBand)

		if bucketBand, existsInBucket := bucket.Bands[searchBandHash]; existsInBucket {
			bucketMatchCount++
			for j, elem := range bucketBand.Elements {
				candidatesDeduped[elem.Value.GetID()] = bucketBand.Elements[j]
			}
		}

	}
	log(fmt.Sprintf("Found %d candidates in %d buckets. Comparing", len(candidatesDeduped), bucketMatchCount))

	// then, check vector similarity for each entry
	results := []LSHResult{}
	for _, c := range candidatesDeduped {
		similarity := algo.CosineSimilarityUint32(c.Singature, searchSignature)
		if similarity >= hashSimilarity {
			valueCpy := c.Value
			results = append(results, LSHResult{Score: similarity, Value: valueCpy})
		}
	}
	log(fmt.Sprintf("Found %d results with good hash similarity, pruned %d", len(results), len(candidatesDeduped)-len(results)))

	// order the results by score
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results
}

type LSHSearchResult struct {
	Score float64
	Val   any
}

