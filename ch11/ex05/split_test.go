package main

import (
	"strings"
	"testing"
)

func TestSplit(t *testing.T) {
	var tests = []struct {
		inputStr string
		inputSep string
		want     int
	}{
		{"", "", 0},
		{"", ":", 1},
		{"a:b:c", ":", 3},
		{"a:b:c", ",", 1},
		{"a:b:c", "", 5},
	}
	for _, test := range tests {
		words := strings.Split(test.inputStr, test.inputSep)
		if got := len(words); got != test.want {
			t.Errorf(
				"Split(%q, %q) returned %d words, want %d",
				test.inputStr, test.inputSep, got, test.want,
			)
		}
	}
}
