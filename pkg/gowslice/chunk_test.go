package gowslice

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChunkSlice(t *testing.T) {
	var tests = []struct {
		in       []int
		nChunks  uint
		expected [][]int
		testName string
	}{
		{
			in:      []int{1, 2, 3},
			nChunks: 0,
			expected: [][]int{
				{1, 2, 3},
			},
			testName: "nChunks 0",
		},
		{
			in:      []int{1, 2, 3},
			nChunks: 1,
			expected: [][]int{
				{1, 2, 3},
			},
			testName: "with 1 chunk",
		},
		{
			in:      []int{1, 2, 3, 4, 5, 6, 7},
			nChunks: 2,
			expected: [][]int{
				{1, 2, 3, 4},
				{5, 6, 7},
			},
			testName: "with nChunks 2 not perfectly dividable with input",
		},
		{
			in:      []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
			nChunks: 3,
			expected: [][]int{
				{1, 2, 3, 4},
				{5, 6, 7, 8},
				{9, 10, 11},
			},
			testName: "with nChunks 3 not perfectly dividable with input",
		},
		{
			in:      []int{1, 2},
			nChunks: 3,
			expected: [][]int{
				{1},
				{2},
				{},
			},
			testName: "with nChunks bigger than input",
		},
	}

	for _, ut := range tests {
		chunks := ChunkSlice(ut.in, ut.nChunks)

		msg := fmt.Sprintf("Test %s \nExpected %#v, got %#v", ut.testName, ut.expected, chunks)
		assert.Equal(t, ut.expected, chunks, msg)
	}
}
