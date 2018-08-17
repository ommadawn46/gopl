package uniqspace

import (
	"unicode"
	"unicode/utf8"
)

func uniqspace(bytes []byte) []byte {
	if len(bytes) <= 0 {
		return bytes
	}
	size := 0
	prevIsSpace := false
	for i := 0; i < len(bytes); {
		r, s := utf8.DecodeRune(bytes[i:])
		isSpace := unicode.IsSpace(r)
		if !prevIsSpace && isSpace {
			bytes[size] = byte(' ')
			size += 1
		} else if !isSpace {
			copy(bytes[size:], bytes[i:i+s])
			size += s
		}
		prevIsSpace = isSpace
		i += s
	}
	return bytes[:size]
}
