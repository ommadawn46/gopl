package intset

import (
	"bytes"
	"fmt"
)

const uintSize = 32 << (^uint(0) >> 63)

type BitIntSet struct {
	words []uint
}

func (s *BitIntSet) Has(x int) bool {
	word, bit := x/uintSize, uint(x%uintSize)
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}

func (s *BitIntSet) Add(x int) {
	word, bit := x/uintSize, uint(x%uintSize)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

func (s *BitIntSet) AddAll(xs ...int) {
	for _, x := range xs {
		s.Add(x)
	}
}

func (s *BitIntSet) UnionWith(t IntSet) {
	switch v := t.(type) {
	case *BitIntSet:
		for i, tword := range v.words {
			if i < len(s.words) {
				s.words[i] |= tword
			} else {
				s.words = append(s.words, tword)
			}
		}
	default:
		for _, x := range t.Elems() {
			s.Add(x)
		}
	}
}

func (s *BitIntSet) IntersectWith(t IntSet) {
	switch v := t.(type) {
	case *BitIntSet:
		for i, tword := range v.words {
			if i < len(s.words) {
				s.words[i] &= tword
			} else {
				break
			}
		}
	default:
		for _, x := range s.Elems() {
			if !t.Has(x) {
				s.Remove(x)
			}
		}
	}
}

func (s *BitIntSet) DifferenceWith(t IntSet) {
	switch v := t.(type) {
	case *BitIntSet:
		for i, tword := range v.words {
			if i < len(s.words) {
				s.words[i] |= tword
				s.words[i] ^= tword
			} else {
				s.words = append(s.words, tword)
			}
		}
	default:
		for _, x := range s.Elems() {
			if t.Has(x) {
				s.Remove(x)
			}
		}
	}
}

func (s *BitIntSet) SymmetricDifference(t IntSet) {
	switch v := t.(type) {
	case *BitIntSet:
		for i, tword := range v.words {
			if i < len(s.words) {
				s.words[i] ^= tword
			} else {
				s.words = append(s.words, tword)
			}
		}
	default:
		for _, x := range t.Elems() {
			if s.Has(x) {
				s.Remove(x)
			} else {
				s.Add(x)
			}
		}
	}
}

func popcount(x uint) int {
	var count int
	for i := 0; i < uintSize; i++ {
		count += int(x & 1)
		x >>= 1
	}
	return count
}

func (s *BitIntSet) Len() int {
	count := 0
	for _, word := range s.words {
		count += popcount(word)
	}
	return count
}

func (s *BitIntSet) Remove(x int) {
	word, bit := x/uintSize, uint(x%uintSize)
	if word < len(s.words) {
		s.words[word] |= 1 << bit
		s.words[word] ^= 1 << bit
	}
}

func (s *BitIntSet) Clear() {
	for i := range s.words {
		s.words[i] = 0
	}
}

func (s *BitIntSet) Copy() IntSet {
	new_s := &BitIntSet{}
	new_s.words = make([]uint, len(s.words))
	copy(new_s.words, s.words)
	return new_s
}

func (s *BitIntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < uintSize; j++ {
			if word&(1<<uint(j)) != 0 {
				if buf.Len() > len("{") {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%d", uintSize*i+j)
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}

func (s *BitIntSet) Elems() (elems []int) {
	for i, word := range s.words {
		bit := uint(0)
		for (word >> bit) > 0 {
			bit += 1
			if (word>>bit)&1 == 1 {
				elems = append(elems, int(bit)+i*uintSize)
			}
		}
	}
	return
}
