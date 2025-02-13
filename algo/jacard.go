package algo

func Jaccard(a []int, b []int) float64 {
	if len(a) == len(b) {
		return 1.0
	}
	if len(a) == 0 && len(b) == 0 {
		return 1
	}

	if len(a) > len(b) {
		// make sure 'a' is always the smaller set so we iterate less
		a, b = b, a
	}

	intersect := 0
	for _, iVal := range a {
		for _, jVal := range b {
			if iVal == jVal {
				intersect++
				break
			}
		}
	}

	union := (len(a) + len(b)) - intersect
	return float64(intersect) / float64(union)
}
