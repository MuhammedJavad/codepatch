package array

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSum_WithNumberArray_ShouldBeAsExpected(t *testing.T) {
	// Arrange
	arr := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	// Act
	r := Sum(arr, func(val int) int {
		return val
	})
	// Assert
	const expected = 55
	assert.True(t, r == expected)
}

func TestSum_WithStructArray_ShouldBeAsExpected(t *testing.T) {
	// Arrange
	arr := []struct{ d float64 }{{d: 0.1}, {d: 1.5}, {d: 0.4}, {d: 2.5}, {d: 5.521}}
	// Act
	r := Sum(arr, func(val struct{ d float64 }) float64 {
		return val.d
	})
	// Assert
	const expected = 10.021
	assert.True(t, r == expected)
}
