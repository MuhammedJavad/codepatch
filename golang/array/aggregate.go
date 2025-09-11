package array

// Numbers is a constraint for integer and float numeric types.
type Numbers interface {
	int | int8 | int16 | int32 | int64 | float32 | float64
}

// Chunk returns a new slice of slices, splitting arr into chunks of size chunkSize.
func Chunk[T interface{}](arr []T, chunkSize int) [][]T {
	var chunkedArray [][]T

	for i := 0; i < len(arr); i += chunkSize {
		end := i + chunkSize

		if end > len(arr) {
			end = len(arr)
		}

		chunkedArray = append(chunkedArray, arr[i:end])
	}

	return chunkedArray
}

// Sum aggregates arr by summing selector(val) into a zero value of N.
func Sum[T any, N Numbers](arr []T, selector func(val T) N) N {
	var summed N
	for i := 0; i < len(arr); i++ {
		r := selector(arr[i])
		summed += r
	}
	return summed
}
