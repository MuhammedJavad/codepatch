package array

import (
	"testing"
)

func TestChunk(t *testing.T) {
	out := Chunk([]int{1, 2, 3, 4, 5}, 2)
	want := [][]int{{1, 2}, {3, 4}, {5}}
	if len(out) != len(want) {
		t.Fatalf("unexpected length: got %d want %d", len(out), len(want))
	}
	for i := range want {
		if len(out[i]) != len(want[i]) {
			t.Fatalf("chunk %d: got %v want %v", i, out[i], want[i])
		}
		for j := range want[i] {
			if out[i][j] != want[i][j] {
				t.Fatalf("elem (%d,%d): got %d want %d", i, j, out[i][j], want[i][j])
			}
		}
	}
}

func TestSum_WithNumberArray_ShouldBeAsExpected(t *testing.T) {
	// Arrange
	arr := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	// Act
	r := Sum(arr, func(val int) int {
		return val
	})
	// Assert
	const expected = 55
	if r != expected {
		t.Fatalf("got %v want %v", r, expected)
	}
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
	if r != expected {
		t.Fatalf("got %v want %v", r, expected)
	}
}
