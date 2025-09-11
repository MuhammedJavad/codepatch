package array

// Find returns a pointer to the first element matching predicate, or nil if none.
func Find[TIn any](arr []TIn, predicate func(val TIn) bool) *TIn {
	for i := range arr {
		if predicate(arr[i]) {
			return &arr[i]
		}
	}
	return nil
}

// Filter returns a new slice of elements for which predicate returns true.
func Filter[TIn any](arr []TIn, predicate func(val *TIn) bool) []TIn {
	var r []TIn
	for i := range arr {
		if predicate(&arr[i]) {
			r = append(r, arr[i])
		}
	}
	return r
}

// FilterAndProject maps arr through selector and drops nil results.
func FilterAndProject[TIn any, TOut any](arr []TIn, selector func(val TIn, index int) *TOut) []TOut {
	var output []TOut
	for i := range arr {
		out := selector(arr[i], i)
		if out == nil {
			continue
		}
		output = append(output, *out)
	}
	return output
}
