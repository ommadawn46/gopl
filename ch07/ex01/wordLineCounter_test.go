package counter

import (
	"testing"
)

func TestLineCount(t *testing.T) {
	for _, test := range []struct {
		input    string
		expected int
	}{
		{"", 0},
		{"ABCD", 1},
		{"ABCD\nEFGH", 2},
		{"ABCD\nEFGH\nIJKL", 3},
		{"ABCD\nEFGH\nIJKL\n\n\n\n\n\n\nMNOP", 10},
	} {
		var c LineCounter
		c.Write([]byte(test.input))
		actual := int(c)
		if actual != test.expected {
			t.Errorf("actual %v want %v", actual, test.expected)
		}
	}
}

func TestWordCount(t *testing.T) {
	for _, test := range []struct {
		input    string
		expected int
	}{
		{"", 0},
		{"ABCD", 1},
		{"ABCD EFGH", 2},
		{"ABCD EFGH IJKL", 3},
		{"ABCD EFGH IJKL          MNOP", 4},
	} {
		var c WordCounter
		c.Write([]byte(test.input))
		actual := int(c)
		if actual != test.expected {
			t.Errorf("actual %v want %v", actual, test.expected)
		}
	}
}
