package array

import (
	"testing"
)

func TestFind(t *testing.T) {
	v := Find([]int{1, 3, 4, 6}, func(x int) bool { return x%2 == 0 })
	if v == nil || *v != 4 {
		t.Fatalf("got %v", v)
	}
	if Find([]int{1, 3, 5}, func(x int) bool { return x%2 == 0 }) != nil {
		t.Fatal("expected nil")
	}
}

func TestFilter(t *testing.T) {
	out := Filter([]int{1, 2, 3, 4}, func(v *int) bool { return *v%2 == 0 })
	if len(out) != 2 || out[0] != 2 || out[1] != 4 {
		t.Fatalf("got %v", out)
	}
}

func TestFilterAndProject(t *testing.T) {
	out := FilterAndProject([]int{1, 2, 3}, func(v int, i int) *string {
		if v%2 == 0 {
			s := "idx"
			return &s
		}
		return nil
	})
	if len(out) != 1 {
		t.Fatalf("got %v", out)
	}
}
