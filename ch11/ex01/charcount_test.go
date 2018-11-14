package main

import (
	"bufio"
	"reflect"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestCharCount(t *testing.T) {
	var tests = []struct {
		input       string
		wantCounts  map[rune]int
		wantUtflen  [utf8.UTFMax + 1]int
		wantInvalid int
	}{
		{
			"",
			map[rune]int{},
			[utf8.UTFMax + 1]int{0, 0, 0, 0, 0},
			0,
		},
		{
			"\xff\xfc\xf8",
			map[rune]int{},
			[utf8.UTFMax + 1]int{0, 0, 0, 0, 0},
			3,
		},
		{
			"Hello, 世界.",
			map[rune]int{'H': 1, 'e': 1, 'l': 2, 'o': 1, ',': 1, ' ': 1, '世': 1, '界': 1, '.': 1},
			[utf8.UTFMax + 1]int{0, 8, 0, 2, 0},
			0,
		},
		{
			"H\xf4ell\xf0o, 世\xec界.",
			map[rune]int{'H': 1, 'e': 1, 'l': 2, 'o': 1, ',': 1, ' ': 1, '世': 1, '界': 1, '.': 1},
			[utf8.UTFMax + 1]int{0, 8, 0, 2, 0},
			3,
		},
	}

	for _, test := range tests {
		reader := bufio.NewReader(strings.NewReader(test.input))
		actualCounts, actualUtflen, actualInvalid := charcount(reader)
		if !reflect.DeepEqual(actualCounts, test.wantCounts) {
			t.Errorf("counts\ninput: %q\nactual: %v\n  want: %v", test.input, actualCounts, test.wantCounts)
		}
		if actualUtflen != test.wantUtflen {
			t.Errorf("utflen\ninput: %q\nactual: %v\n  want: %v", test.input, actualUtflen, test.wantUtflen)
		}
		if actualInvalid != test.wantInvalid {
			t.Errorf("invalid\ninput: %q\nactual: %v\n  want: %v", test.input, actualInvalid, test.wantInvalid)
		}
	}
}
