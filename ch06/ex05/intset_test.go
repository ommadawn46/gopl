package intset

import (
	"reflect"
	"testing"
)

func TestIntSetEmpty(t *testing.T) {
	s := &IntSet{}
	actual := s.String()
	expected := "{}"
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}
}

func TestIntSetAdd(t *testing.T) {
	s := &IntSet{}
	s.Add(1 << 8)
	s.Add(1 << 9)
	s.Add(1 << 10)
	s.Add(1 << 11)
	s.Add(1 << 12)
	actual := s.String()
	expected := "{256 512 1024 2048 4096}"
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}
}

func TestIntSetHas(t *testing.T) {
	s := &IntSet{}
	s.Add(1 << 8)
	s.Add(1 << 9)
	s.Add(1 << 10)
	s.Add(1 << 11)
	s.Add(1 << 12)

	actual := s.Has(1 << 10)
	expected := true
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}

	actual = s.Has(1 << 13)
	expected = false
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}
}

func TestIntSetUnionWith(t *testing.T) {
	s1 := &IntSet{}
	s1.Add(1 << 8)
	s1.Add(1 << 9)
	s1.Add(1 << 10)
	s1.Add(1 << 11)
	s1.Add(1 << 12)

	s2 := &IntSet{}
	s2.Add(1 << 10)
	s2.Add(1 << 11)
	s2.Add(1 << 12)
	s2.Add(1 << 13)
	s2.Add(1 << 14)
	s1.UnionWith(s2)

	actual := s1.String()
	expected := "{256 512 1024 2048 4096 8192 16384}"
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}
}

func TestIntSetLen(t *testing.T) {
	s := &IntSet{}
	s.Add(1 << 8)
	s.Add(1 << 9)
	s.Add(1 << 10)
	s.Add(1 << 11)
	s.Add(1 << 12)

	actual := s.Len()
	expected := 5
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}
}

func TestIntSetRemove(t *testing.T) {
	s := &IntSet{}
	s.Add(1 << 8)
	s.Add(1 << 9)
	s.Add(1 << 10)
	s.Add(1 << 11)
	s.Add(1 << 12)

	s.Remove(1 << 10)
	actual := s.String()
	expected := "{256 512 2048 4096}"
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}

	s.Remove(1 << 9)
	s.Remove(1 << 11)
	actual = s.String()
	expected = "{256 4096}"
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}

	s.Remove(1 << 8)
	s.Remove(1 << 12)
	actual = s.String()
	expected = "{}"
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}

	s.Remove(1 << 10)
	actual = s.String()
	expected = "{}"
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}
}

func TestIntSetClear(t *testing.T) {
	s := &IntSet{}
	s.Add(1 << 8)
	s.Add(1 << 9)
	s.Add(1 << 10)
	s.Add(1 << 11)
	s.Add(1 << 12)
	s.Clear()
	actual := s.String()
	expected := "{}"
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}
}

func TestIntSetCopy(t *testing.T) {
	s1 := &IntSet{}
	s1.Add(1 << 8)
	s1.Add(1 << 9)
	s1.Add(1 << 10)
	s1.Add(1 << 11)
	s1.Add(1 << 12)
	s2 := s1.Copy()
	actual := s1.String() == s2.String()
	expected := true
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}
}

func TestIntSetAddAll(t *testing.T) {
	s := &IntSet{}
	s.AddAll(1, 2, 3, 11, 27, 37, 41, 73, 77, 116, 154)
	actual := s.String()
	expected := "{1 2 3 11 27 37 41 73 77 116 154}"
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}
}

func TestIntSetIntersectWith(t *testing.T) {
	s1 := &IntSet{}
	s1.AddAll(1, 3, 11, 37, 41, 77, 116, 154)
	s2 := &IntSet{}
	s2.AddAll(1, 2, 3, 11, 27, 37, 41, 73, 154, 1024)
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
