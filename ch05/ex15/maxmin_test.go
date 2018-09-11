package maxmin

import (
	"testing"
)

func TestMax(t *testing.T) {
	for _, test := range []struct {
		vals     []int
		expected int
	}{
		{[]int{}, 0},
		{[]int{1, 1, 1}, 1},
		{[]int{1, 2, 3, 4, 5}, 5},
		{[]int{100, -50, 1000, -5000, 10000, 3000, -8000}, 10000},
	} {
		actual := max(test.vals...)
		if actual != test.expected {
			t.Errorf("actual %v want %v", actual, test.expected)
		}
	}
}

func TestMax_(t *testing.T) {
	for _, test := range []struct {
		first    int
		vals     []int
		expected int
	}{
		{-10, []int{}, -10},
		{1, []int{1, 1, 1}, 1},
		{0, []int{1, 2, 3, 4, 5}, 5},
		{-100, []int{100, -50, 1000, -5000, 10000, 3000, -8000}, 10000},
	} {
		actual := max_(test.first, test.vals...)
		if actual != test.expected {
			t.Errorf("actual %v want %v", actual, test.expected)
		}
	}
}

func TestMin(t *testing.T) {
	for _, test := range []struct {
		vals     []int
		expected int
	}{
		{[]int{}, 0},
		{[]int{1, 1, 1}, 1},
		{[]int{1, 2, 3, 4, 5}, 1},
		{[]int{100, -50, 1000, -5000, 10000, 3000, -8000}, -8000},
	} {
		actual := min(test.vals...)
		if actual != test.expected {
			t.Errorf("actual %v want %v", actual, test.expected)
		}
	}
}

func TestMin_(t *testing.T) {
	for _, test := range []struct {
		first    int
		vals     []int
		expected int
	}{
		{-10, []int{}, -10},
		{1, []int{1, 1, 1}, 1},
		{0, []int{1, 2, 3, 4, 5}, 0},
		{-100, []int{100, -50, 1000, -5000, 10000, 3000, -8000}, -8000},
	} {
		actual := min_(test.first, test.vals...)
		if actual != test.expected {
			t.Errorf("actual %v want %v", actual, test.expected)
		}
	}
}
