package unicoderev

import (
	"unicode/utf8"
)

func reverse(b []byte) {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
}

func unicodeReverse(b []byte) {
	for i := 0; i < len(b); {
		_, s := utf8.DecodeRune(b[i:])
		reverse(b[i : i+s])
		i += s
	}
	reverse(b)
}
