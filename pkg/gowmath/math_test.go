package gowmath

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRoundToFixedDecimals(t *testing.T) {
	testValues := []float64{1.23456789, 456.25815, 456.0, 785.1}
	expectedResultsOn2dec := []float64{1.23, 456.26, 456.0, 785.1}
	expectedResultsOn3dec := []float64{1.235, 456.258, 456.0, 785.1}

	t.Run("Rounding on 2 decimals", func(t *testing.T) {
		for index, value := range testValues {
			assert.Equal(t, expectedResultsOn2dec[index], RoundToFixedDecimals(value, 2))
		}
	})
	t.Run("Rounding on 3 decimals", func(t *testing.T) {
		for index, value := range testValues {
			assert.Equal(t, expectedResultsOn3dec[index], RoundToFixedDecimals(value, 3))
		}
	})
}

func TestMod(t *testing.T) {
	assert.Equal(t, Mod(-5, 24), 19)
	assert.Equal(t, Mod(5, 24), 5)
	assert.Equal(t, Mod(27, 24), 3)
}

func TestMax(t *testing.T) {
	assert.Equal(t, 2, Max(1, 2))
	assert.Equal(t, 2, Max(2, 1))
	assert.Equal(t, 1, Max(1, 1))
}
