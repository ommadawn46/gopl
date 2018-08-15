package unicoderev

import (
  "testing"
  "reflect"
)

func TestUnicoderevEmpty(t *testing.T) {
  actual := []byte("")
  unicodeReverse(actual)
  expected := []byte("")
  if !reflect.DeepEqual(actual, expected) {
    t.Errorf("actual %v want %v", actual, expected)
  }
}

func TestUnicoderevAscii(t *testing.T) {
  actual := []byte("ABCDEFGH")
  unicodeReverse(actual)
  expected := []byte("HGFEDCBA")
  if !reflect.DeepEqual(actual, expected) {
    t.Errorf("actual %v want %v", actual, expected)
  }
}

func TestUnicoderevUtf8(t *testing.T) {
  actual := []byte("あいうえお")
  unicodeReverse(actual)
  expected := []byte("おえういあ")
  if !reflect.DeepEqual(actual, expected) {
    t.Errorf("actual %v want %v", actual, expected)
  }
}
