package rev

import (
  "testing"
)

func TestReverseZero(t *testing.T) {
  actual := [8]int{}
  reverse(&actual)
  expected := [8]int{}
  if actual != expected {
    t.Errorf("actual %v want %v", actual, expected)
  }
}

func TestReverseArray(t *testing.T) {
  actual := [8]int{0, 1, 2, 3, 4, 5, 6, 7}
  reverse(&actual)
  expected := [8]int{7, 6, 5, 4, 3, 2, 1, 0}
  if actual != expected {
    t.Errorf("actual %v want %v", actual, expected)
  }
}
