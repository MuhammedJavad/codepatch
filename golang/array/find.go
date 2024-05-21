package array

func Find[TIn any](arr []TIn, predicate func(val TIn) bool) *TIn {
	for i := range arr {
		if predicate(arr[i]) {
			return &arr[i]
		}
	}
	return nil
}
