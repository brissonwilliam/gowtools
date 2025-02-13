package helpers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsAlpha3Code(t *testing.T) {
	t.Run("Succeed", func(t *testing.T) {
		result := IsAlpha3Code("CAN")
		assert.True(t, result)
	})

	t.Run("Succeed with lower and upper case", func(t *testing.T) {
		result := IsAlpha3Code("cAn")
		assert.True(t, result)
	})

	t.Run("Fails because too short", func(t *testing.T) {
		result := IsAlpha3Code("CA")
		assert.False(t, result)
	})

	t.Run("Fails because too long", func(t *testing.T) {
		result := IsAlpha3Code("CANADA")
		assert.False(t, result)
	})

	t.Run("Fails because not in list", func(t *testing.T) {
		result := IsAlpha3Code("LST")
		assert.False(t, result)
	})

}
