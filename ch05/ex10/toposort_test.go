package toposort

import (
	"testing"
)

func check(order []string) bool {
  contains := func(s []string, t string) bool {
    for _, a := range s {
      if a == t {
        return true
      }
    }
    return false
  }
  for i, item := range order {
    dependencies := prereqs[item]
    for _, depend_item := range dependencies {
      if !contains(order[:i], depend_item) {
        return false
      }
    }
  }
  return true
}

func TestCheck(t *testing.T) {
	for _, test := range []struct
	{
		order []string
		expected bool
	}{
		{
			[]string{
				"intro to programming",
				"discrete math",
				"data structures",
				"algorithms",
				"linear algebra",
				"calculus",
				"formal languages",
				"computer organiczation",
				"compilers",
				"operating systems",
				"networks",
				"databases",
				"programming languages",
			}, true,
		},
		{
			[]string{
				"intro to programming",
				"discrete math",
				"data structures",
				"linear algebra",
				"calculus",
				"formal languages",
				"computer organiczation",
				"operating systems",
				"programming languages",
				"compilers",
				"networks",
				"databases",
				"algorithms",
			}, true,
		},
		{
			[]string{
				"calculus",
				"linear algebra",
				"data structures",
				"discrete math",
				"intro to programming",
				"databases",
				"formal languages",
				"algorithms",
				"networks",
				"operating systems",
				"computer organiczation",
				"programming languages",
				"compilers",
			}, false,
		},
		{
			[]string{
				"intro to programming",
				"programming languages",
				"calculus",
				"linear algebra",
				"data structures",
				"discrete math",
				"formal languages",
				"operating systems",
				"computer organiczation",
				"compilers",
				"networks",
				"databases",
				"algorithms",
			}, false,
		},
	}{
		actual := check(test.order)
		if actual != test.expected {
			t.Errorf("actual %v want %v", actual, test.expected)
		}
	}
}

func TestTopoSort(t *testing.T) {
	actual := topoSort(prereqs)
	if !check(actual) {
		t.Errorf("actual %v test failue", actual)
	}
}
