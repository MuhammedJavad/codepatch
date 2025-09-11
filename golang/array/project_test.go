package array

import (
	"testing"
)

func TestProject(t *testing.T) {
	out := Project([]int{1, 2, 3}, func(v int, i int) int { return v * v })
	if len(out) != 3 || out[0] != 1 || out[1] != 4 || out[2] != 9 {
		t.Fatalf("got %v", out)
	}
}

func TestProjectErr(t *testing.T) {
	out, err := ProjectErr([]int{1}, func(v int, i int) (*int, error) { r := v + 1; return &r, nil })
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(out) != 1 || out[0] != 2 {
		t.Fatalf("got %v", out)
	}
}

func TestProjectMap(t *testing.T) {
	m := map[int]string{1: "a", 2: "b"}
	out := ProjectMap(m, func(v string, k int) string { return v + "x" })
	if len(out) != 2 {
		t.Fatalf("got %v", out)
	}
	// order not guaranteed; check membership
	got := map[string]bool{}
	for _, v := range out {
		got[v] = true
	}
	if !(got["ax"] && got["bx"]) {
		t.Fatalf("got %v", out)
	}
}

func TestFlat(t *testing.T) {
	type group struct{ items []int }
	arr := []group{{items: []int{1, 2}}, {items: []int{}}}
	out := Flat(arr, func(g group) []int { return g.items }, func(g group, it int) *int { return &it })
	if len(out) != 2 || out[0] != 1 || out[1] != 2 {
		t.Fatalf("got %v", out)
	}
}
