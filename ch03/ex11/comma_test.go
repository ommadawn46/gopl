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

func TestCommaPlus(t *testing.T) {
	actual := comma("+18467968320")
	expected := "+18,467,968,320"
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}
}

func TestCommaMinus(t *testing.T) {
	actual := comma("-18467968320")
	expected := "-18,467,968,320"
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}
}

func TestCommaDecimal(t *testing.T) {
	actual := comma("1234567890.0987654321")
	expected := "1,234,567,890.0987654321"
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}
}

func TestCommaPlusDecimal(t *testing.T) {
	actual := comma("+1234567890.0987654321")
	expected := "+1,234,567,890.0987654321"
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}
}

func TestCommaMinusDecimal(t *testing.T) {
	actual := comma("-1234567890.0987654321")
	expected := "-1,234,567,890.0987654321"
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}
}
