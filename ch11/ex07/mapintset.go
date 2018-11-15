package intset

import (
	"bytes"
	"fmt"
	"sort"
)

type MapIntSet struct {
	m map[int]struct{}
}

func NewMapIntSet() *MapIntSet {
	return &MapIntSet{make(map[int]struct{})}
}

func (s *MapIntSet) Has(x int) bool {
	_, ok := s.m[x]
	return ok
}

func (s *MapIntSet) Add(x int) {
	s.m[x] = struct{}{}
}

func (s *MapIntSet) AddAll(xs ...int) {
	for _, x := range xs {
		s.Add(x)
	}
}

func (s *MapIntSet) UnionWith(t IntSet) {
	for _, x := range t.Elems() {
		s.Add(x)
	}
}

func (s *MapIntSet) IntersectWith(t IntSet) {
	for _, x := range s.Elems() {
		if !t.Has(x) {
			s.Remove(x)
		}
	}
}

func (s *MapIntSet) DifferenceWith(t IntSet) {
	for _, x := range s.Elems() {
		if t.Has(x) {
			s.Remove(x)
		}
	}
}

func (s *MapIntSet) SymmetricDifference(t IntSet) {
	for _, x := range t.Elems() {
		if s.Has(x) {
			s.Remove(x)
		} else {
			s.Add(x)
		}
	}
}

func (s *MapIntSet) Len() int {
	return len(s.m)
}

func (s *MapIntSet) Remove(x int) {
	delete(s.m, x)
}

func (s *MapIntSet) Clear() {
	s.m = make(map[int]struct{})
}

func (s *MapIntSet) Copy() IntSet {
	new_s := NewMapIntSet()
	for key, value := range s.m {
		new_s.m[key] = value
	}
	return new_s
}

func (s *MapIntSet) String() string {
	elems := s.Elems()
	var buf bytes.Buffer
	buf.WriteByte('{')
	for _, n := range elems {
		if buf.Len() > len("{") {
			buf.WriteByte(' ')
		}
		fmt.Fprintf(&buf, "%d", n)
	}
	buf.WriteByte('}')
	return buf.String()
}

func (s *MapIntSet) Elems() (elems []int) {
	for v := range s.m {
		elems = append(elems, v)
	}
	sort.IntSlice(elems).Sort()
	return
}
