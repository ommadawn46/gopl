package palindrome

import (
	"testing"
)

type ByteSlice []byte

func (x ByteSlice) Len() int {
	return len(x)
}
func (x ByteSlice) Less(i, j int) bool {
	return x[i] < x[j]
}
func (x ByteSlice) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

func TestPalindrome(t *testing.T) {
	for _, test := range []struct {
		input    string
		expected bool
	}{
		{"", true},
		{"ABCCBA", true},
		{"ABCDCBA", true},
		{"ABCDEF", false},
		{"ABCDABC", false},
		{"ABCDECBA", false},
	} {
		actual := IsPalindrome(ByteSlice([]byte(test.input)))
		if actual != test.expected {
			t.Errorf("actual %v want %v", actual, test.expected)
		}
	}
}
