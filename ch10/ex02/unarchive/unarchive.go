package unarchive

import (
	"bufio"
	"errors"
	"io"
	"os"
)

var ErrFormat = errors.New("image: unknown format")

type ArchivedFile interface {
	Open() (io.ReadCloser, error)
	Path() string
	FileInfo() os.FileInfo
}

type reader interface {
	io.Reader
	Peek(int) ([]byte, error)
}

func asReader(r io.Reader) reader {
	if rr, ok := r.(reader); ok {
		return rr
	}
	return bufio.NewReader(r)
}

var formats []format

type format struct {
	name, magic string
	magicOffset int
	unarchive   func(io.Reader) ([]ArchivedFile, error)
}

func RegisterFormat(name, magic string, magicOffset int, unarchive func(io.Reader) ([]ArchivedFile, error)) {
	formats = append(formats, format{name, magic, magicOffset, unarchive})
}

func match(magic string, b []byte) bool {
	if len(magic) != len(b) {
		return false
	}
	for i, c := range b {
		if magic[i] != c && magic[i] != '?' {
			return false
		}
	}
	return true
}

func sniff(r reader) format {
	for _, f := range formats {
		b, err := r.Peek(f.magicOffset + len(f.magic))
		if err == nil && match(f.magic, b[f.magicOffset:]) {
			return f
		}
	}
	return format{}
}

func Unarchive(r io.Reader) ([]ArchivedFile, string, error) {
	rr := asReader(r)
	format := sniff(rr)
	if format.unarchive == nil {
		return nil, "", ErrFormat
	}
	files, err := format.unarchive(rr)
	return files, format.name, err
}
