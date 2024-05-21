package array

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

func Sum[T any, N Numbers](arr []T, selector func(val T) N) N {
	var summed N
	for i := 0; i < len(arr); i++ {
		r := selector(arr[i])
		summed += r
	}
	return summed
}
