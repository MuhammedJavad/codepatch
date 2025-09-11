package array

import (
	"errors"
	"testing"
)

func TestAll(t *testing.T) {
	if !All([]int{2, 4, 6}, func(v int) bool { return v%2 == 0 }) {
		t.Fatal("expected true")
	}
	if All([]int{2, 3, 6}, func(v int) bool { return v%2 == 0 }) {
		t.Fatal("expected false")
	}
}

func TestAny(t *testing.T) {
	if !Any([]int{1, 3, 4}, func(v int) bool { return v%2 == 0 }) {
		t.Fatal("expected true")
	}
	if Any([]int{1, 3, 5}, func(v int) bool { return v%2 == 0 }) {
		t.Fatal("expected false")
	}
}

func TestAnyErr(t *testing.T) {
	err := AnyErr([]int{1, 2, 3}, func(v int) error {
		if v == 2 {
			return errors.New("boom")
		}
		return nil
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if e := AnyErr([]int{1, 3, 5}, func(v int) error { return nil }); e != nil {
		t.Fatalf("unexpected error: %v", e)
	}
}

func TestContains(t *testing.T) {
	if !Contains([]string{"a", "b"}, "a") {
		t.Fatal("expected true")
	}
	if Contains([]string{"a", "b"}, "c") {
		t.Fatal("expected false")
	}
}
