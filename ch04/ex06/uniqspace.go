package uniqspace

import (
	"unicode"
	"unicode/utf8"
)

func uniqspace(bytes []byte) []byte {
	if len(bytes) <= 0 {
		return bytes
	}
	blank := 0
	isPrevSpace := false
	for i := 0; i < len(bytes); {
		r, s := utf8.DecodeRune(bytes[i:])
		isSpace := unicode.IsSpace(r)
		if isPrevSpace && isSpace {
			blank += s
		} else if isSpace {
			bytes[i-blank] = byte(' ')
			blank += s - 1
		} else {
			copy(bytes[i-blank:], bytes[i:i+s])
		}
		isPrevSpace = isSpace
		i += s
	}
	return bytes[:len(bytes)-blank]
}
