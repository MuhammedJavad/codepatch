package array

// FlatMap projects values from a map to a flat slice using collectionSelector and resultSelector.
func FlatMap[TKey comparable, TIn any, TC any, TOut any](m map[TKey]TIn, collectionSelector func(TIn) []TC, resultSelector func(TC) TOut) []TOut {
	var output []TOut
	for i := range m {
		out := collectionSelector(m[i])
		for _, v := range out {
			output = append(output, resultSelector(v))
		}
	}
	return output
}

// ToMap constructs a map keyed by keySelector from a slice.
func ToMap[TKey comparable, TIn any](in []TIn, keySelector func(TIn) TKey) map[TKey]TIn {
	m := make(map[TKey]TIn)
	for _, val := range in {
		m[keySelector(val)] = val
	}
	return m
}
