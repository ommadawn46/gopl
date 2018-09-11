package countingWriter

import (
	"io/ioutil"
	"testing"
)

func TestCountingWriter(t *testing.T) {
	w1, _ := CountingWriter(ioutil.Discard)
	w1.Write([]byte("ABCD"))
	w2, _ := CountingWriter(w1)
	w2.Write([]byte("EFGH"))
	w3, _ := CountingWriter(w2)
	w3.Write([]byte("IJKL"))
	w4, _ := CountingWriter(w3)
	w4.Write([]byte("MNOP"))
	w5, _ := CountingWriter(w4)
	w5.Write([]byte(""))

	actual := w1.(*writeWrapper).count
	expected := int64(16)
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}

	actual = w2.(*writeWrapper).count
	expected = int64(12)
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}

	actual = w3.(*writeWrapper).count
	expected = int64(8)
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}

	actual = w4.(*writeWrapper).count
	expected = int64(4)
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}

	actual = w5.(*writeWrapper).count
	expected = int64(0)
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}
}
