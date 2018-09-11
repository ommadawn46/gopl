package intset

import (
	"testing"
)

func TestIntSetIntersectWith(t *testing.T) {
	s1 := &IntSet{}
	s1.AddAll(1, 3, 11, 37, 41, 77, 116, 154)
	s2 := &IntSet{}
	s2.AddAll(1, 2, 3, 11, 27, 37, 41, 73, 154)
	s1.IntersectWith(s2)
	actual := s1.String()
	expected := "{1 3 11 37 41 154}"
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}
}

func TestIntSetDifferenceWith(t *testing.T) {
	s1 := &IntSet{}
	s1.AddAll(1, 3, 11, 37, 41, 77, 116, 154)
	s2 := &IntSet{}
	s2.AddAll(1, 2, 3, 11, 27, 37, 41, 73, 154)
	s1.DifferenceWith(s2)
	actual := s1.String()
	expected := "{77 116}"
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}
}

func TestIntSetSymmetricDifference(t *testing.T) {
	s1 := &IntSet{}
	s1.AddAll(1, 3, 11, 37, 41, 77, 116, 154)
	s2 := &IntSet{}
	s2.AddAll(1, 2, 3, 11, 27, 37, 41, 73, 154)
	s1.SymmetricDifference(s2)
	actual := s1.String()
	expected := "{2 27 73 77 116}"
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}
}
