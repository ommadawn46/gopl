package intset

import (
	"testing"
)

func TestIntSetAddAll(t *testing.T) {
	s := &IntSet{}
	s.AddAll(1, 2, 3, 11, 27, 37, 41, 73, 77, 116, 154)
	actual := s.String()
	expected := "{1 2 3 11 27 37 41 73 77 116 154}"
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}
}
