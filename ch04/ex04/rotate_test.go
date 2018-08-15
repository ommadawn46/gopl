package rotate

import (
  "testing"
  "reflect"
)

func TestRotateEmpty(t *testing.T) {
  actual := []int{}
  rotate(actual, 3)
  expected := []int{}
  if !reflect.DeepEqual(actual, expected) {
    t.Errorf("actual %v want %v", actual, expected)
  }
}

func TestRotateZero(t *testing.T) {
  actual := []int{0, 1, 2, 3, 4, 5}
  rotate(actual, 0)
  expected := []int{0, 1, 2, 3, 4, 5}
  if !reflect.DeepEqual(actual, expected) {
    t.Errorf("actual %v want %v", actual, expected)
  }
}

func TestRotateSlice(t *testing.T) {
  actual := []int{0, 1, 2, 3, 4, 5}
  rotate(actual, 3)
  expected := []int{3, 4, 5, 0, 1, 2}
  if !reflect.DeepEqual(actual, expected) {
    t.Errorf("actual %v want %v", actual, expected)
  }
}

func TestRotateOverLen(t *testing.T) {
  actual := []int{0, 1, 2, 3, 4, 5}
  rotate(actual, 16)
  expected := []int{4, 5, 0, 1, 2, 3}
  if !reflect.DeepEqual(actual, expected) {
    t.Errorf("actual %v want %v", actual, expected)
  }
}
