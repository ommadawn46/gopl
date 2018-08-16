package comma

import (
	"testing"
)

func TestCommaLen0(t *testing.T) {
	actual := comma("")
	expected := ""
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}
}

func TestCommaLen3(t *testing.T) {
	actual := comma("184")
	expected := "184"
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}
}

func TestCommaLen10(t *testing.T) {
	actual := comma("9566027596")
	expected := "9,566,027,596"
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}
}

func TestCommaLen32(t *testing.T) {
	actual := comma("39570284651543759757968476938621")
	expected := "39,570,284,651,543,759,757,968,476,938,621"
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}
}
