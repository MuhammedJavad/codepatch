package array

// All returns true when every element in arr satisfies predicate.
func All[T any](arr []T, predicate func(val T) bool) bool {
	for i := 0; i < len(arr); i++ {
		if !predicate(arr[i]) {
			return false
		}
	}
	return true
}

// Any returns true when any element in arr satisfies predicate.
func Any[TIn any](arr []TIn, predicate func(val TIn) bool) bool {
	for i := 0; i < len(arr); i++ {
		if predicate(arr[i]) {
			return true
		}
	}
	return false
}

// AnyErr returns the first error from predicate over elements in arr, or nil if none.
func AnyErr[TIn any](arr []TIn, predicate func(val TIn) error) error {
	for i := 0; i < len(arr); i++ {
		if err := predicate(arr[i]); err != nil {
			return err
		}
	}
	return nil
}

// Contains reports whether slice contains element using ==.
func Contains[T comparable](slice []T, element T) bool {
	for i := 0; i < len(slice); i++ {
		if slice[i] == element {
			return true
		}
	}
	return false
}
