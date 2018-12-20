package cyclic

import (
	"testing"
)

func TestCyclic(t *testing.T) {
	type CyclePtr *CyclePtr
	var cyclePtr1, cyclePtr2 CyclePtr
	cyclePtr1 = &cyclePtr1
	cyclePtr2 = &cyclePtr2

	var nonCyclePtr1, nonCyclePtr2 CyclePtr
	nonCyclePtr1 = &nonCyclePtr2

	type link struct {
		value string
		tail  *link
	}
	cycleLink1, cycleLink2, cycleLink3 := &link{value: "a"}, &link{value: "b"}, &link{value: "c"}
	cycleLink1.tail, cycleLink2.tail, cycleLink3.tail = cycleLink2, cycleLink1, cycleLink3

	nonCycleLink1, nonCycleLink2, nonCycleLink3 := &link{value: "a"}, &link{value: "b"}, &link{value: "c"}
	nonCycleLink1.tail, nonCycleLink2.tail = nonCycleLink2, nonCycleLink3

	type CycleSlice []CycleSlice
	var cycleSlice = make(CycleSlice, 1)
	cycleSlice[0] = make(CycleSlice, 1)
	cycleSlice[0][0] = cycleSlice

	var nonCycleSlice = make(CycleSlice, 1)
	nonCycleSlice[0] = make(CycleSlice, 1)

	for _, test := range []struct {
		v    interface{}
		want bool
	}{
		{nil, false}, {0, false}, {"ABC", false},
		{cyclePtr1, true}, {cyclePtr2, true},
		{nonCyclePtr1, false}, {nonCyclePtr2, false},
		{cycleLink1, true}, {cycleLink2, true}, {cycleLink3, true},
		{nonCycleLink1, false}, {nonCycleLink2, false}, {nonCycleLink3, false},
		{cycleSlice, true}, {nonCycleSlice, false},
	} {
		if Cyclic(test.v) != test.want {
			t.Errorf("Cyclic(%v) = %t", test.v, !test.want)
		}
	}
}
