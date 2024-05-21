package array

type Numbers interface {
	int | int8 | int16 | int32 | int64 | float32 | float64
}

func BubbleSort[T any, N Numbers](arr []T, selector func(val T) N) {
	n := len(arr)

	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			c := selector(arr[j])
			n := selector(arr[j+1])

			if c > n {
				// swap arr[j] and arr[j+1]
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
	}
}

func BubbleSortDesc[T any](arr []T, selector func(val T) float64) {
	n := len(arr)

	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			c := selector(arr[j])
			n := selector(arr[j+1])

			if c < n { // Change comparison operator to less than
				// swap arr[j] and arr[j+1]
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
	}
}
