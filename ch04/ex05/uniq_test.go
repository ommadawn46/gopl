package uniq

import (
	"reflect"
	"testing"
)

func TestUniqEmpty(t *testing.T) {
	actual := uniq([]string{})
	expected := []string{}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("actual %v want %v", actual, expected)
	}
}

func TestUniqSlice1(t *testing.T) {
	actual := uniq([]string{"A", "A", "A", "B", "B", "C", "D", "D"})
	expected := []string{"A", "B", "C", "D"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("actual %v want %v", actual, expected)
	}
}

func TestUniqSlice2(t *testing.T) {
	actual := uniq([]string{"A", "A", "A", "B", "B", "A", "B", "B", "C", "C", "C", "B", "B"})
	expected := []string{"A", "B", "A", "B", "C", "B"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("actual %v want %v", actual, expected)
	}
}
