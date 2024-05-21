package array

func IsEmpty[TIn any](arr []TIn) bool {
	return arr == nil || len(arr) <= 0
}
