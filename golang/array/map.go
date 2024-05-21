package array

func Map[TIn any, TOut any](arr []TIn, selector func(val TIn, index int) TOut) []TOut {
	var output []TOut
	for i := range arr {
		out := selector(arr[i], i)
		output = append(output, out)
	}
	return output
}

func ToMap[TKey comparable, T any](arr []T, keySelector func(T) TKey) map[TKey]T {
	result := make(map[TKey]T)

	for _, item := range arr {
		key := keySelector(item)
		result[key] = item
	}

	return result
}
