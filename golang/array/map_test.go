package array

import (
	"testing"
)

func TestFlatMap(t *testing.T) {
	type user struct{ roles []string }
	m := map[int]user{1: {roles: []string{"a", "b"}}, 2: {roles: []string{"c"}}}
	out := FlatMap(m, func(u user) []string { return u.roles }, func(r string) string { return r })
	if len(out) != 3 {
		t.Fatalf("got %v", out)
	}
	// order is not guaranteed; check membership
	got := map[string]bool{}
	for _, v := range out {
		got[v] = true
	}
	if !(got["a"] && got["b"] && got["c"]) {
		t.Fatalf("got %v", out)
	}
}

func TestToMap(t *testing.T) {
	type u struct{ id int }
	arr := []u{{1}, {2}}
	m := ToMap(arr, func(x u) int { return x.id })
	if len(m) != 2 || m[2] != (u{2}) {
		t.Fatalf("got %v", m)
	}
}
