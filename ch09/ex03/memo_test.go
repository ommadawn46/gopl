package memo

import (
	"fmt"
	"testing"
)

func TestDone(t *testing.T) {
	var actual, expected interface{}
	var actualErr, expectedErr error

	cancelErr := fmt.Errorf("function cancelled")

	echo := func(s string, done <-chan struct{}) (interface{}, error) {
		select {
		case <-done:
			return "", cancelErr
		default:
			return s, nil
		}
	}

	done := make(chan struct{})
	memo := New(Func(echo))

	v, err := memo.Get("hoge", done)
	actual = v
	actualErr = err
	expected = "hoge"
	expectedErr = nil
	if actual != expected || actualErr != expectedErr {
		t.Errorf("actual %v, %v want %v, %v", actual, actualErr, expected, expectedErr)
	}

	close(done)
	v, err = memo.Get("hoge", done)
	actual = v
	actualErr = err
	expected = "hoge"
	expectedErr = nil
	if actual != expected || actualErr != expectedErr {
		t.Errorf("actual %v, %v want %v, %v", actual, actualErr, expected, expectedErr)
	}

	v, err = memo.Get("fuga", done)
	actual = v
	actualErr = err
	expected = ""
	expectedErr = cancelErr
	if actual != expected || actualErr != expectedErr {
		t.Errorf("actual %v, %v want %v, %v", actual, actualErr, expected, expectedErr)
	}
}
