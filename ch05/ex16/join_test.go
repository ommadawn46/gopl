package join

import (
	"testing"
)

func TestJoin(t *testing.T) {
	for _, test := range []struct {
		sep      string
		s        []string
		expected string
	}{
		{"", []string{}, ""},
		{"abc", []string{}, ""},
		{"", []string{"a", "b", "c"}, "abc"},
		{" ", []string{"foo", "bar", "baz"}, "foo bar baz"},
		{", ", []string{"hoge", "fuga", "piyo", "hogera"}, "hoge, fuga, piyo, hogera"},
	} {
		actual := join(test.sep, test.s...)
		if actual != test.expected {
			t.Errorf("actual %v want %v", actual, test.expected)
		}
	}
}
