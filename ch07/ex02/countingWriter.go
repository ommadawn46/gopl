package countingWriter

import (
	"io"
)

type writeWrapper struct {
	count  int64
	writer io.Writer
}

func (w *writeWrapper) Write(p []byte) (n int, err error) {
	n, err = w.writer.Write(p)
	w.count += int64(n)
	return
}

func CountingWriter(w io.Writer) (io.Writer, *int64) {
	var wrapper writeWrapper
	wrapper.writer = w
	return &wrapper, &(wrapper.count)
}
