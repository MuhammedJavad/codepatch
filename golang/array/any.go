package array

func All[T any](arr []T, predicate func(val T) bool) bool {
	for i := 0; i < len(arr); i++ {
		if !predicate(arr[i]) {
			return false
		}
	}
	return true
}

// Any returns true if any element in the array satisfies the predicate; otherwise, it returns false.
func Any[TIn any](arr []TIn, predicate func(val TIn) bool) bool {
	for i := 0; i < len(arr); i++ {
		if predicate(arr[i]) {
			return true
		}
	}
	return false
}

// Contains checks if a slice contains a specific element.
func Contains[T comparable](slice []T, element T) bool {
	for i := 0; i < len(slice); i++ {
		if slice[i] == element {
			return true
		}
	}
	return false
}
