package array

// Project maps each element of arr through selector into a new slice.
func Project[TIn any, TOut any](arr []TIn, selector func(val TIn, index int) TOut) []TOut {
	output := make([]TOut, 0, len(arr))
	for i := range arr {
		out := selector(arr[i], i)
		output = append(output, out)
	}
	return output
}

// ProjectErr maps arr through selector and returns either results or an error.
func ProjectErr[TIn any, TOut any](arr []TIn, selector func(val TIn, index int) (*TOut, error)) ([]TOut, error) {
	var output []TOut
	for i := range arr {
		out, err := selector(arr[i], i)
		if err != nil {
			return nil, err
		}
		output = append(output, *out)
	}
	return output, nil
}

// ProjectMap maps entries of m to a slice using selector over (value, key).
func ProjectMap[TKey comparable, TIn any, TOut any](m map[TKey]TIn, selector func(val TIn, key TKey) TOut) []TOut {
	output := make([]TOut, 0, len(m))
	for i := range m {
		out := selector(m[i], i)
		output = append(output, out)
	}
	return output
}

// Flat flattens the projection of each element into a single slice, skipping nil results.
func Flat[TIn any, TC any, TOut any](m []TIn, collectionSelector func(TIn) []TC, resultSelector func(TIn, TC) *TOut) []TOut {
	var output []TOut
	for i := range m {
		out := collectionSelector(m[i])
		for _, v := range out {
			result := resultSelector(m[i], v)
			if result == nil {
				continue
			}
			output = append(output, *result)
		}
	}
	return output
}
