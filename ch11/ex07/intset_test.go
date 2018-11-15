package intset

import (
	"reflect"
	"testing"
)

func testIntSetEmpty(t *testing.T, s IntSet) {
	actual := s.String()
	expected := "{}"
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}
}

func testIntSetAdd(t *testing.T, s IntSet) {
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

func testIntSetHas(t *testing.T, s IntSet) {
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

func testIntSetUnionWith(t *testing.T, s1, s2 IntSet) {
	s1.Add(1 << 8)
	s1.Add(1 << 9)
	s1.Add(1 << 10)
	s1.Add(1 << 11)
	s1.Add(1 << 12)

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

func testIntSetLen(t *testing.T, s IntSet) {
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

func testIntSetRemove(t *testing.T, s IntSet) {
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

func testIntSetClear(t *testing.T, s IntSet) {
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

func testIntSetCopy(t *testing.T, s1 IntSet) {
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

func testIntSetAddAll(t *testing.T, s IntSet) {
	s.AddAll(1, 2, 3, 11, 27, 37, 41, 73, 77, 116, 154)
	actual := s.String()
	expected := "{1 2 3 11 27 37 41 73 77 116 154}"
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}
}

func testIntSetIntersectWith(t *testing.T, s1, s2 IntSet) {
	s1.AddAll(1, 3, 11, 37, 41, 77, 116, 154)
	s2.AddAll(1, 2, 3, 11, 27, 37, 41, 73, 154, 1024)
	s1.IntersectWith(s2)
	actual := s1.String()
	expected := "{1 3 11 37 41 154}"
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}
}

func testIntSetDifferenceWith(t *testing.T, s1, s2 IntSet) {
	s1.AddAll(1, 3, 11, 37, 41, 77, 116, 154)
	s2.AddAll(1, 2, 3, 11, 27, 37, 41, 73, 154)
	s1.DifferenceWith(s2)
	actual := s1.String()
	expected := "{77 116}"
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}
}

func testIntSetSymmetricDifference(t *testing.T, s1, s2 IntSet) {
	s1.AddAll(1, 3, 11, 37, 41, 77, 116, 154)
	s2.AddAll(1, 2, 3, 11, 27, 37, 41, 73, 154)
	s1.SymmetricDifference(s2)
	actual := s1.String()
	expected := "{2 27 73 77 116}"
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}
}

func testIntSetEmptyElems(t *testing.T, s IntSet) {
	actual := len(s.Elems())
	expected := 0
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}
}

func testIntSetElems(t *testing.T, s IntSet) {
	s.AddAll(1, 2, 3, 11, 27, 37, 41, 73, 77, 116, 154)
	actual := s.Elems()
	expected := []int{1, 2, 3, 11, 27, 37, 41, 73, 77, 116, 154}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("actual %v want %v", actual, expected)
	}
}

func TestBitIntSetEmpty(t *testing.T) {
	testIntSetEmpty(t, &BitIntSet{})
}

func TestMapIntSetEmpty(t *testing.T) {
	testIntSetEmpty(t, NewMapIntSet())
}

func TestBitIntSetAdd(t *testing.T) {
	testIntSetAdd(t, &BitIntSet{})
}

func TestMapIntSetAdd(t *testing.T) {
	testIntSetAdd(t, NewMapIntSet())
}

func TestBitIntSetHas(t *testing.T) {
	testIntSetHas(t, &BitIntSet{})
}

func TestMapIntSetHas(t *testing.T) {
	testIntSetHas(t, NewMapIntSet())
}

func TestBitIntSetLen(t *testing.T) {
	testIntSetLen(t, &BitIntSet{})
}

func TestMapIntSetLen(t *testing.T) {
	testIntSetLen(t, NewMapIntSet())
}

func TestBitIntSetRemove(t *testing.T) {
	testIntSetRemove(t, &BitIntSet{})
}

func TestMapIntSetRemove(t *testing.T) {
	testIntSetRemove(t, NewMapIntSet())
}

func TestBitIntSetClear(t *testing.T) {
	testIntSetClear(t, &BitIntSet{})
}

func TestMapIntSetClear(t *testing.T) {
	testIntSetClear(t, NewMapIntSet())
}

func TestBitIntSetCopy(t *testing.T) {
	testIntSetCopy(t, &BitIntSet{})
}

func TestMapIntSetCopy(t *testing.T) {
	testIntSetCopy(t, NewMapIntSet())
}

func TestBitIntSetAddAll(t *testing.T) {
	testIntSetAddAll(t, &BitIntSet{})
}

func TestMapIntSetAddAll(t *testing.T) {
	testIntSetAddAll(t, NewMapIntSet())
}

func TestBitIntSetUnionWith(t *testing.T) {
	testIntSetUnionWith(t, &BitIntSet{}, &BitIntSet{})
}

func TestMapIntSetUnionWith(t *testing.T) {
	testIntSetUnionWith(t, NewMapIntSet(), NewMapIntSet())
}

func TestBitIntSetIntersectWith(t *testing.T) {
	testIntSetIntersectWith(t, &BitIntSet{}, &BitIntSet{})
}

func TestMapIntSetIntersectWith(t *testing.T) {
	testIntSetIntersectWith(t, NewMapIntSet(), NewMapIntSet())
}

func TestBitIntSetDifferenceWith(t *testing.T) {
	testIntSetUnionWith(t, &BitIntSet{}, &BitIntSet{})
}

func TestMapIntSetDifferenceWith(t *testing.T) {
	testIntSetUnionWith(t, NewMapIntSet(), NewMapIntSet())
}

func TestBitIntSetSymmetricDifference(t *testing.T) {
	testIntSetUnionWith(t, &BitIntSet{}, &BitIntSet{})
}

func TestMapIntSetSymmetricDifference(t *testing.T) {
	testIntSetUnionWith(t, NewMapIntSet(), NewMapIntSet())
}

func TestBitIntSetEmptyElems(t *testing.T) {
	testIntSetEmptyElems(t, &BitIntSet{})
}

func TestMapIntSetEmptyElems(t *testing.T) {
	testIntSetEmptyElems(t, NewMapIntSet())
}

func TestBitIntSetElems(t *testing.T) {
	testIntSetElems(t, &BitIntSet{})
}

func TestMapIntSetElems(t *testing.T) {
	testIntSetElems(t, NewMapIntSet())
}
