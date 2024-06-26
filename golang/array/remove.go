package array

func Remove(slice *[]string, value any) {
	for i := 0; i < len(*slice); i++ {
		if (*slice)[i] == value {
			// Remove the element at index i
			*slice = append((*slice)[:i], (*slice)[i+1:]...)
			return
		}
	}
}
