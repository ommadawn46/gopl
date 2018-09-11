package reader

import (
	"io"
)

type stringReader string

func (r *stringReader) Read(p []byte) (n int, err error) {
	n = copy(p, *r)
	*r = (*r)[n:]
	if len(*r) == 0 {
		err = io.EOF
	}
	return
}

func NewReader(s string) io.Reader {
	new_r := stringReader(s)
	return &new_r
}
