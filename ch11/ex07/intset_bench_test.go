package intset

import (
	"math/rand"
	"testing"
)

func benchmarkIntSetAdd(b *testing.B, s IntSet, size int) {
	for i := uint64(0); i < uint64(b.N); i++ {
		t := s.Copy()
		t.Add(rand.Intn(size))
	}
}

func benchmarkIntSetRemove(b *testing.B, s IntSet, size int) {
	for i := uint64(0); i < uint64(b.N); i++ {
		t := s.Copy()
		n := rand.Intn(size)
		t.Add(n)
		t.Remove(n)
	}
}

func benchmarkIntSetLen(b *testing.B, s IntSet, size int) {
	for i := uint64(0); i < uint64(b.N); i++ {
		t := s.Copy()
		for j := 0; j < 10; j++ {
			t.Add(rand.Intn(size))
		}
		t.Len()
	}
}

func benchmarkIntSetHas(b *testing.B, s IntSet, size int) {
	s.Add(size)
	for i := uint64(0); i < uint64(b.N); i++ {
		s.Has(rand.Intn(size))
	}
}

func benchmarkIntSetUnionWith(b *testing.B, s1 IntSet, s2 IntSet, size int) {
	for i := uint64(0); i < uint64(b.N); i++ {
		t1 := s1.Copy()
		t2 := s2.Copy()
		t1.Add(rand.Intn(size))
		t2.Add(rand.Intn(size))
		t1.UnionWith(t2)
	}
}

func benchmarkIntSetIntersectWith(b *testing.B, s1 IntSet, s2 IntSet, size int) {
	for i := uint64(0); i < uint64(b.N); i++ {
		t1 := s1.Copy()
		t2 := s2.Copy()
		t1.Add(rand.Intn(size))
		t2.Add(rand.Intn(size))
		t1.IntersectWith(t2)
	}
}

func benchmarkIntSetDifferenceWith(b *testing.B, s1 IntSet, s2 IntSet, size int) {
	for i := uint64(0); i < uint64(b.N); i++ {
		t1 := s1.Copy()
		t2 := s2.Copy()
		t1.Add(rand.Intn(size))
		t2.Add(rand.Intn(size))
		t1.DifferenceWith(t2)
	}
}

func benchmarkIntSetSymmetricDifference(b *testing.B, s1 IntSet, s2 IntSet, size int) {
	for i := uint64(0); i < uint64(b.N); i++ {
		t1 := s1.Copy()
		t2 := s2.Copy()
		t1.Add(rand.Intn(size))
		t2.Add(rand.Intn(size))
		t1.SymmetricDifference(t2)
	}
}

func BenchmarkBitIntSetAdd16(b *testing.B) { benchmarkIntSetAdd(b, &BitIntSet{}, 1<<16) }
func BenchmarkBitIntSetAdd28(b *testing.B) { benchmarkIntSetAdd(b, &BitIntSet{}, 1<<28) }
func BenchmarkBitIntSetAdd31(b *testing.B) { benchmarkIntSetAdd(b, &BitIntSet{}, 1<<31) }
func BenchmarkMapIntSetAdd16(b *testing.B) { benchmarkIntSetAdd(b, NewMapIntSet(), 1<<16) }
func BenchmarkMapIntSetAdd28(b *testing.B) { benchmarkIntSetAdd(b, NewMapIntSet(), 1<<28) }
func BenchmarkMapIntSetAdd31(b *testing.B) { benchmarkIntSetAdd(b, NewMapIntSet(), 1<<31) }

func BenchmarkBitIntSetRemove16(b *testing.B) { benchmarkIntSetRemove(b, &BitIntSet{}, 1<<16) }
func BenchmarkBitIntSetRemove28(b *testing.B) { benchmarkIntSetRemove(b, &BitIntSet{}, 1<<28) }
func BenchmarkBitIntSetRemove31(b *testing.B) { benchmarkIntSetRemove(b, &BitIntSet{}, 1<<31) }
func BenchmarkMapIntSetRemove16(b *testing.B) { benchmarkIntSetRemove(b, NewMapIntSet(), 1<<16) }
func BenchmarkMapIntSetRemove28(b *testing.B) { benchmarkIntSetRemove(b, NewMapIntSet(), 1<<28) }
func BenchmarkMapIntSetRemove31(b *testing.B) { benchmarkIntSetRemove(b, NewMapIntSet(), 1<<31) }

func BenchmarkBitIntSetLen16(b *testing.B) { benchmarkIntSetLen(b, &BitIntSet{}, 1<<16) }
func BenchmarkBitIntSetLen28(b *testing.B) { benchmarkIntSetLen(b, &BitIntSet{}, 1<<28) }
func BenchmarkBitIntSetLen31(b *testing.B) { benchmarkIntSetLen(b, &BitIntSet{}, 1<<31) }
func BenchmarkMapIntSetLen16(b *testing.B) { benchmarkIntSetLen(b, NewMapIntSet(), 1<<16) }
func BenchmarkMapIntSetLen28(b *testing.B) { benchmarkIntSetLen(b, NewMapIntSet(), 1<<28) }
func BenchmarkMapIntSetLen31(b *testing.B) { benchmarkIntSetLen(b, NewMapIntSet(), 1<<31) }

func BenchmarkBitIntSetHas16(b *testing.B) { benchmarkIntSetHas(b, &BitIntSet{}, 1<<16) }
func BenchmarkBitIntSetHas28(b *testing.B) { benchmarkIntSetHas(b, &BitIntSet{}, 1<<28) }
func BenchmarkBitIntSetHas31(b *testing.B) { benchmarkIntSetHas(b, &BitIntSet{}, 1<<31) }
func BenchmarkMapIntSetHas16(b *testing.B) { benchmarkIntSetHas(b, NewMapIntSet(), 1<<16) }
func BenchmarkMapIntSetHas28(b *testing.B) { benchmarkIntSetHas(b, NewMapIntSet(), 1<<28) }
func BenchmarkMapIntSetHas31(b *testing.B) { benchmarkIntSetHas(b, NewMapIntSet(), 1<<31) }

func BenchmarkBitIntSetUnionWith16(b *testing.B) {
	benchmarkIntSetUnionWith(b, &BitIntSet{}, &BitIntSet{}, 1<<16)
}
func BenchmarkBitIntSetUnionWith28(b *testing.B) {
	benchmarkIntSetUnionWith(b, &BitIntSet{}, &BitIntSet{}, 1<<28)
}
func BenchmarkBitIntSetUnionWith31(b *testing.B) {
	benchmarkIntSetUnionWith(b, &BitIntSet{}, &BitIntSet{}, 1<<31)
}
func BenchmarkMapIntSetUnionWith16(b *testing.B) {
	benchmarkIntSetUnionWith(b, NewMapIntSet(), NewMapIntSet(), 1<<16)
}
func BenchmarkMapIntSetUnionWith28(b *testing.B) {
	benchmarkIntSetUnionWith(b, NewMapIntSet(), NewMapIntSet(), 1<<28)
}
func BenchmarkMapIntSetUnionWith31(b *testing.B) {
	benchmarkIntSetUnionWith(b, NewMapIntSet(), NewMapIntSet(), 1<<31)
}

func BenchmarkBitIntSetIntersectWith16(b *testing.B) {
	benchmarkIntSetIntersectWith(b, &BitIntSet{}, &BitIntSet{}, 1<<16)
}
func BenchmarkBitIntSetIntersectWith28(b *testing.B) {
	benchmarkIntSetIntersectWith(b, &BitIntSet{}, &BitIntSet{}, 1<<28)
}
func BenchmarkBitIntSetIntersectWith31(b *testing.B) {
	benchmarkIntSetIntersectWith(b, &BitIntSet{}, &BitIntSet{}, 1<<31)
}
func BenchmarkMapIntSetIntersectWith16(b *testing.B) {
	benchmarkIntSetIntersectWith(b, NewMapIntSet(), NewMapIntSet(), 1<<16)
}
func BenchmarkMapIntSetIntersectWith28(b *testing.B) {
	benchmarkIntSetIntersectWith(b, NewMapIntSet(), NewMapIntSet(), 1<<28)
}
func BenchmarkMapIntSetIntersectWith31(b *testing.B) {
	benchmarkIntSetIntersectWith(b, NewMapIntSet(), NewMapIntSet(), 1<<31)
}

func BenchmarkBitIntSetDifferenceWith16(b *testing.B) {
	benchmarkIntSetDifferenceWith(b, &BitIntSet{}, &BitIntSet{}, 1<<16)
}
func BenchmarkBitIntSetDifferenceWith28(b *testing.B) {
	benchmarkIntSetDifferenceWith(b, &BitIntSet{}, &BitIntSet{}, 1<<28)
}
func BenchmarkBitIntSetDifferenceWith31(b *testing.B) {
	benchmarkIntSetDifferenceWith(b, &BitIntSet{}, &BitIntSet{}, 1<<31)
}
func BenchmarkMapIntSetDifferenceWith16(b *testing.B) {
	benchmarkIntSetDifferenceWith(b, NewMapIntSet(), NewMapIntSet(), 1<<16)
}
func BenchmarkMapIntSetDifferenceWith28(b *testing.B) {
	benchmarkIntSetDifferenceWith(b, NewMapIntSet(), NewMapIntSet(), 1<<28)
}
func BenchmarkMapIntSetDifferenceWith31(b *testing.B) {
	benchmarkIntSetDifferenceWith(b, NewMapIntSet(), NewMapIntSet(), 1<<31)
}

func BenchmarkBitIntSetSymmetricDifference16(b *testing.B) {
	benchmarkIntSetSymmetricDifference(b, &BitIntSet{}, &BitIntSet{}, 1<<16)
}
func BenchmarkBitIntSetSymmetricDifference28(b *testing.B) {
	benchmarkIntSetSymmetricDifference(b, &BitIntSet{}, &BitIntSet{}, 1<<28)
}
func BenchmarkBitIntSetSymmetricDifference31(b *testing.B) {
	benchmarkIntSetSymmetricDifference(b, &BitIntSet{}, &BitIntSet{}, 1<<31)
}
func BenchmarkMapIntSetSymmetricDifference16(b *testing.B) {
	benchmarkIntSetSymmetricDifference(b, NewMapIntSet(), NewMapIntSet(), 1<<16)
}
func BenchmarkMapIntSetSymmetricDifference28(b *testing.B) {
	benchmarkIntSetSymmetricDifference(b, NewMapIntSet(), NewMapIntSet(), 1<<28)
}
func BenchmarkMapIntSetSymmetricDifference31(b *testing.B) {
	benchmarkIntSetSymmetricDifference(b, NewMapIntSet(), NewMapIntSet(), 1<<31)
}
