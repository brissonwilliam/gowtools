package algo

import "math"

func CosineSimilarityUint32(a []uint32, b []uint32) float64 {
	if len(a) != len(b) {
		panic("a and b must be of equal length")
	}

	dotProduct := uint64(0)
	sumA := uint64(0)
	sumB := uint64(0)
	for i := range a {
		dotProduct += uint64(a[i]) * uint64(b[i])
		sumA += (uint64(a[i]) * uint64(a[i]))
		sumB += (uint64(b[i]) * uint64(b[i]))
	}

	magnitudeA := math.Sqrt(float64(sumA))
	magnitudeB := math.Sqrt(float64(sumB))

	return float64(dotProduct) / (magnitudeA * magnitudeB)
}

func CosineSimilarityUint16(a []uint16, b []uint16) float64 {
	if len(a) != len(b) {
		panic("a and b must be of equal length")
	}

	dotProduct := uint64(0)
	sumA := uint64(0)
	sumB := uint64(0)
	for i := range a {
		dotProduct += uint64(a[i]) * uint64(b[i])
		sumA += (uint64(a[i]) * uint64(a[i]))
		sumB += (uint64(b[i]) * uint64(b[i]))
	}

	magnitudeA := math.Sqrt(float64(sumA))
	magnitudeB := math.Sqrt(float64(sumB))

	return float64(dotProduct) / (magnitudeA * magnitudeB)
}

func CosineSimilarity(a []int, b []int) float64 {
	if len(a) != len(b) {
		panic("a and b must be of equal length")
	}

	dotProduct := 0
	sumA := uint64(0)
	sumB := uint64(0)
	for i := range a {
		dotProduct += a[i] * b[i]
		sumA += (uint64(a[i]) * uint64(a[i]))
		sumB += (uint64(b[i]) * uint64(b[i]))
	}

	magnitudeA := math.Sqrt(float64(sumA))
	magnitudeB := math.Sqrt(float64(sumB))

	return float64(dotProduct) / (magnitudeA * magnitudeB)
}
