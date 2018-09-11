package reader

import (
	"bytes"
	"io"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestLineRead(t *testing.T) {
	r := NewReader("ABCDEFGH")
	buf := make([]byte, 16)
	r.Read(buf)
	actual := buf
	expected := []byte{65, 66, 67, 68, 69, 70, 71, 72, 0, 0, 0, 0, 0, 0, 0, 0}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("actual %v want %v", actual, expected)
	}
}

func TestLineReadAll(t *testing.T) {
	r := NewReader("ABCDEFGH")
	actual, _ := ioutil.ReadAll(r)
	expected := []byte{65, 66, 67, 68, 69, 70, 71, 72}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("actual %v want %v", actual, expected)
	}
}

func TestLineCopy(t *testing.T) {
	r := NewReader("ABCDEFGH")
	buf := new(bytes.Buffer)
	io.Copy(buf, r)
	actual := buf.String()
	expected := "ABCDEFGH"
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}
}
