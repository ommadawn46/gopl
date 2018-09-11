package limitreader

import (
	"io"
)

type limitReader struct {
	r        io.Reader
	n, limit int
}

func (r *limitReader) Read(p []byte) (n int, err error) {
	lim := r.limit - r.n
	if len(p) < lim {
		lim = len(p)
	}
	n, err = r.r.Read(p[:lim])
	r.n += n
	if r.n >= r.limit {
		err = io.EOF
	}
	return
}

func LimitReader(r io.Reader, limit int) io.Reader {
	return &limitReader{r: r, limit: limit}
}
