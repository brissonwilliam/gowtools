package gowslice

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPaginate(t *testing.T) {
	type sampleObject struct {
		Number int
	}

	sampleObjects := []sampleObject{
		sampleObject{1},
		sampleObject{2},
		sampleObject{3},
		sampleObject{4},
	}

	t.Run("Test Paginate succeeds when limit and offset are set", func(t *testing.T) {
		limit := uint64(1)
		result := Paginate(sampleObjects, &limit, 2)

		assert.Equal(t, sampleObjects[2:3], result)
	})

	t.Run("Test Paginate succeeds when limit is set", func(t *testing.T) {
		limit := uint64(2)
		result := Paginate(sampleObjects, &limit, 0)

		assert.Equal(t, sampleObjects[0:2], result)
	})

	t.Run("Test Paginate succeeds when offset is set", func(t *testing.T) {
		limit := uint64(250)
		result := Paginate(sampleObjects, &limit, 2)

		assert.Equal(t, sampleObjects[2:4], result)
	})

	t.Run("Test Paginate succeeds when offset overflows slice length", func(t *testing.T) {
		limit := uint64(250)
		offset := uint64(len(sampleObjects) + 1)
		result := Paginate(sampleObjects, &limit, offset)

		assert.Equal(t, []sampleObject{}, result)
	})

	t.Run("Test Paginate succeeds when offset is under but limit is over slice length", func(t *testing.T) {
		limit := uint64(len(sampleObjects) + 2)
		offset := uint64(len(sampleObjects) - 2)
		result := Paginate(sampleObjects, &limit, offset)

		assert.Equal(t, sampleObjects[2:4], result)
	})

	t.Run("Test Paginate succeeds with offset and limit under slice length", func(t *testing.T) {
		limit := uint64(2)
		offset := uint64(2)
		result := Paginate(sampleObjects, &limit, offset)

		assert.Equal(t, sampleObjects[2:4], result)
	})

	t.Run("Test Paginate returns the entire slice when limit is nil", func(t *testing.T) {
		result := Paginate(sampleObjects, nil, 0)

		assert.Equal(t, sampleObjects, result)
	})
}
