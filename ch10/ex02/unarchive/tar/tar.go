package tar

import (
	"archive/tar"
	"bytes"
	"io"
	"io/ioutil"
	"os"

	"github.com/ommadawn46/gopl/ch10/ex02/unarchive"
)

type tarFile struct {
	header *tar.Header
	reader *tarFileReader
}

func (f *tarFile) Open() (io.ReadCloser, error) {
	return f.reader, nil
}

func (f *tarFile) Path() string {
	return f.header.Name
}

func (f *tarFile) FileInfo() os.FileInfo {
	return f.header.FileInfo()
}

type tarFileReader struct {
	readSeeker io.ReadSeeker
	offset     int64
	size       int64
	cur        int64
}

func (fr *tarFileReader) Read(b []byte) (n int, err error) {
	fr.readSeeker.Seek(fr.offset+fr.cur, 0)
	if int64(len(b)) > fr.size-fr.cur {
		b = b[:fr.size-fr.cur]
	}
	if len(b) > 0 {
		n, err = fr.readSeeker.Read(b)
		fr.cur += int64(n)
	}
	switch {
	case err == io.EOF && fr.cur < fr.size:
		return n, io.ErrUnexpectedEOF
	case err == nil && fr.cur == fr.size:
		return n, io.EOF
	default:
		return n, err
	}
}

func (t *tarFileReader) Close() error {
	return nil
}

func Unarchive(reader io.Reader) ([]unarchive.ArchivedFile, error) {
	var readSeeker io.ReadSeeker
	switch v := reader.(type) {
	case *os.File:
		readSeeker = v
	default:
		// ファイル以外の場合はバッファに一度読み込む
		buf, err := ioutil.ReadAll(reader)
		if err != nil {
			return nil, err
		}
		readSeeker = io.NewSectionReader(bytes.NewReader(buf), 0, int64(len(buf)))
	}
	tarReader := tar.NewReader(readSeeker)

	var files []unarchive.ArchivedFile
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		// 現在のファイルオフセットを取得
		offset, err := readSeeker.Seek(0, 1)
		if err != nil {
			return nil, err
		}
		fr := &tarFileReader{
			readSeeker,
			offset,
			header.Size,
			0,
		}
		files = append(files, &tarFile{header, fr})
	}
	return files, nil
}

func init() {
	magic := "ustar"
	magicOffset := 257
	unarchive.RegisterFormat("tar", magic, magicOffset, Unarchive)
}
