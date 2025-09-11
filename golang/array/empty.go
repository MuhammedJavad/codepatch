package array

// IsEmpty reports whether arr has zero elements. It is safe for nil slices.
func IsEmpty[TIn any](arr []TIn) bool {
	return len(arr) == 0
}
