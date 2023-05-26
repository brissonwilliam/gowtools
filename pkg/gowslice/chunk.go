package gowslice

// ChunkSlice breaks down a slice of data into multiple chunks
// The function will always return nChunks, even if they are empty
func ChunkSlice[T any](data []T, nChunks uint) [][]T {
	dl := uint(len(data))

	if nChunks == 0 || dl == 0 {
		return [][]T{data}
	}

	chunkSize := dl / nChunks
	remainder := dl % nChunks
	chunks := make([][]T, 0, nChunks)

	start := uint(0)
	for i := uint(0); i < nChunks; i++ {
		end := start + chunkSize
		if i < remainder {
			end++
		}
		if end > dl {
			end = dl
		}
		chunks = append(chunks, data[start:end])
		start = end
	}

	return chunks
}
