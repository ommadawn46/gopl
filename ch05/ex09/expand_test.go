package expand

import (
	"strings"
	"testing"
)

func TestExpand(t *testing.T) {
	for _, test := range []struct {
		s        string
		f        func(string) string
		expected string
	}{
		{"", strings.ToUpper, ""},
		{"abcd efgh ijkl", strings.ToUpper, "abcd efgh ijkl"},
		{"$abcd efgh ijkl", strings.ToUpper, "ABCD efgh ijkl"},
		{"abcd $efgh ijkl", strings.ToUpper, "abcd EFGH ijkl"},
		{"$abcd efgh $ijkl", strings.ToUpper, "ABCD efgh IJKL"},
		{"$ABCD EFGH IJKL", strings.ToLower, "abcd EFGH IJKL"},
		{"ABCD $EFGH IJKL", strings.ToLower, "ABCD efgh IJKL"},
		{"ABCD $EFGH $IJKL", strings.ToLower, "ABCD efgh ijkl"},
		{"$abcd efgh ijkl", strings.Title, "Abcd efgh ijkl"},
		{"abcd $efgh ijkl", strings.Title, "abcd Efgh ijkl"},
		{"$abcd $efgh $ijkl", strings.Title, "Abcd Efgh Ijkl"},
	} {
		actual := expand(test.s, test.f)
		if actual != test.expected {
			t.Errorf("actual %v want %v", actual, test.expected)
		}
	}
}
