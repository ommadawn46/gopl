package intset

import (
	"bytes"
	"fmt"
	"sort"
	"testing"
)

func mapToString(set map[int]struct{}) string {
	var list []int
	for v := range set {
		list = append(list, v)
	}
	sort.IntSlice(list).Sort()

	var buf bytes.Buffer
	buf.WriteByte('{')
	for _, n := range list {
		if buf.Len() > len("{") {
			buf.WriteByte(' ')
		}
		fmt.Fprintf(&buf, "%d", n)
	}
	buf.WriteByte('}')
	return buf.String()
}

func TestMapToString(t *testing.T) {
	m := make(map[int]struct{})
	actual := mapToString(m)
	expected := "{}"
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}

	for _, n := range []int{256, 512, 1024, 2048, 4096} {
		m[n] = struct{}{}
	}
	actual = mapToString(m)
	expected = "{256 512 1024 2048 4096}"
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}
}

func TestIntSetEmpty(t *testing.T) {
	s := &IntSet{}
	m := make(map[int]struct{})

	actualStr := s.String()
	expectedStr := mapToString(m)
	if actualStr != expectedStr {
		t.Errorf("string: actual %v want %v", actualStr, expectedStr)
	}

	actualLen := s.Len()
	expectedLen := len(m)
	if actualLen != expectedLen {
		t.Errorf("length: actual %v want %v", actualLen, expectedLen)
	}
}

func TestIntSetAdd(t *testing.T) {
	s := &IntSet{}
	m := make(map[int]struct{})
	for _, n := range []int{256, 512, 1024, 2048, 4096} {
		s.Add(n)
		m[n] = struct{}{}
	}

	actual := s.String()
	expected := mapToString(m)
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}
}

func TestIntSetHas(t *testing.T) {
	s := &IntSet{}
	m := make(map[int]struct{})
	for _, n := range []int{256, 512, 1024, 2048, 4096} {
		s.Add(n)
		m[n] = struct{}{}
	}

	actual := s.Has(1024)
	_, expected := m[1024]
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}

	actual = s.Has(1 << 13)
	_, expected = m[1<<13]
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}
}

func TestIntSetLen(t *testing.T) {
	s := &IntSet{}
	m := make(map[int]struct{})
	for _, n := range []int{256, 512, 1024, 2048, 4096} {
		s.Add(n)
		m[n] = struct{}{}
	}

	actual := s.Len()
	expected := len(m)
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}
}

func TestIntSetRemove(t *testing.T) {
	s := &IntSet{}
	m := make(map[int]struct{})
	for _, n := range []int{256, 512, 1024, 2048, 4096} {
		s.Add(n)
		m[n] = struct{}{}
	}

	s.Remove(1024)
	delete(m, 1024)
	actual := s.String()
	expected := mapToString(m)
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}

	s.Remove(512)
	delete(m, 512)
	actual = s.String()
	expected = mapToString(m)
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}

	s.Remove(256)
	delete(m, 256)
	actual = s.String()
	expected = mapToString(m)
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}

	s.Remove(1024)
	delete(m, 1024)
	actual = s.String()
	expected = mapToString(m)
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}
}
