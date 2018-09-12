package limitreader

import (
	"bytes"
	"io"
	"io/ioutil"
	"reflect"
	"testing"

	"../ex04"
)

func TestReaderRead(t *testing.T) {
	r := reader.NewReader("ABCDEFGH")
	lr := LimitReader(r, 6)
	buf := make([]byte, 3)
	lr.Read(buf)
	actual := buf
	expected := []byte{65, 66, 67}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("actual %v want %v", actual, expected)
	}

	lr.Read(buf)
	actual = buf
	expected = []byte{68, 69, 70}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("actual %v want %v", actual, expected)
	}
}

func TestReaderReadAll(t *testing.T) {
	r := reader.NewReader("ABCDEFGH")
	lr := LimitReader(r, 4)
	actual, _ := ioutil.ReadAll(lr)
	expected := []byte{65, 66, 67, 68}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("actual %v want %v", actual, expected)
	}
}

func TestReaderCopy(t *testing.T) {
	r := reader.NewReader("ABCDEFGH")
	lr := LimitReader(r, 4)
	buf := new(bytes.Buffer)
	io.Copy(buf, lr)
	actual := buf.String()
	expected := "ABCD"
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}
}
