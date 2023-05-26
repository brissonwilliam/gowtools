package gowmath

import "math"

func RoundToFixedDecimals(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return math.Round(num*output) / output
}

// Mod does a modulus operation like most languages when handling negatives values.
// It returns the absolute integers number that are always positive
// In go -5 % 24 = -5
// In others -5 % 24 = 19
func Mod(a, b int) int {
	return (a%b + b) % b
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
