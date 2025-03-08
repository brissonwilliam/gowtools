package algo

func Jaccard(a []int, b []int) float64 {
	if len(a) == 0 && len(b) == 0 {
		return 1
	}

	if len(a) > len(b) {
		// make sure 'a' is always the smaller set so we iterate less
		a, b = b, a
	}

	intersect := map[int]uint8{}
	for _, iVal := range a {
		for _, jVal := range b {
			if iVal == jVal {
				intersect[iVal] = 1
				break
			}
		}
	}

	union := (len(a) + len(b)) - len(intersect)
	return float64(len(intersect)) / float64(union)
}

// JaccardFast uses a map to find intersecting elements, reducing the complexity of the olgarithm to O(n) instead of O(n^2)
func JaccardFast(a []uint16, b map[uint16]uint8) float64 {
	if len(a) == 0 && len(b) == 0 {
		return 1
	}

	intersect := map[uint16]uint8{}
	for _, iVal := range a {
		if _, ok := b[iVal]; ok {
			intersect[iVal] = 1
		}
	}

	union := (len(a) + len(b)) - len(intersect)
	return float64(len(intersect)) / float64(union)
}

