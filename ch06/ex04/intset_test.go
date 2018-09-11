package intset

import (
	"testing"
	"reflect"
)

func TestIntSetEmptyElems(t *testing.T) {
	s := &IntSet{}
	actual := len(s.Elems())
	expected := 0
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}
}

func TestIntSetElems(t *testing.T) {
	s := &IntSet{}
	s.AddAll(1, 2, 3, 11, 27, 37, 41, 73, 77, 116, 154)
	actual := s.Elems()
	expected := []int{1, 2, 3, 11, 27, 37, 41, 73, 77, 116, 154}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("actual %v want %v", actual, expected)
	}
}
