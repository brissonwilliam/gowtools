package gowslice

import "reflect"

// Paginate returns a trimmed slice for the given offset and limit
func Paginate(s any, limit *uint64, offset uint64) any {
	if limit == nil {
		return s
	}

	if !isSlice(s) {
		panic("First parameter must be a slice")
	}

	sValue := reflect.ValueOf(s)
	sLen := uint64(sValue.Len())

	if offset >= sLen {
		return sValue.Slice(0, 0).Interface()
	}

	// endIndex is either max length of slice or the one given with offset
	endIndex := min(sLen, offset+(*limit))

	s = sValue.Slice(int(offset), int(endIndex)).Interface()

	return s
}

func min(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}

func isSlice(s interface{}) bool {
	arrType := reflect.TypeOf(s)
	return arrType.Kind() == reflect.Slice
}
