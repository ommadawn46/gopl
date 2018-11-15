package intset

type IntSet interface {
	Has(x int) bool
	Add(x int)
	AddAll(xs ...int)
	UnionWith(t IntSet)
	IntersectWith(t IntSet)
	DifferenceWith(t IntSet)
	SymmetricDifference(t IntSet)
	Len() int
	Remove(x int)
	Clear()
	Copy() IntSet
	String() string
	Elems() (elems []int)
}
