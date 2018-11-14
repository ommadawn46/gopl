package zip

import (
	"archive/zip"
	"bytes"
	"io"
	"io/ioutil"
	"os"

	"github.com/ommadawn46/the_go_programming_language-training/ch10/ex02/unarchive"
)

type zipFile struct {
	*zip.File
}

func (f *zipFile) Path() string {
	return f.Name
}

func Unarchive(reader io.Reader) ([]unarchive.ArchivedFile, error) {
	var readerAt io.ReaderAt
	var size int64
	switch v := reader.(type) {
	case *os.File:
		readerAt = v
		fileInfo, err := v.Stat()
		if err != nil {
			return nil, err
		}
		size = fileInfo.Size()
	default:
		buf, err := ioutil.ReadAll(reader)
		if err != nil {
			return nil, err
		}
		readerAt = bytes.NewReader(buf)
		size = int64(len(buf))
	}

	zipReader, err := zip.NewReader(readerAt, size)
	if err != nil {
		return nil, err
	}

	var files []unarchive.ArchivedFile
	for _, file := range zipReader.File {
		files = append(files, &zipFile{file})
	}
	return files, nil
}

func init() {
	magic := "PK\x03\x04"
	magicOffset := 0
	unarchive.RegisterFormat("zip", magic, magicOffset, Unarchive)
}
