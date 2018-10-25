package bank

import (
	"testing"
)

func TestBank(t *testing.T) {
	actual := Balance()
	expected := 0
	if actual != expected {
		t.Errorf("Not zero: actual %v want %v", actual, expected)
	}

	Deposit(1000)
	actual = Balance()
	expected = 1000
	if actual != expected {
		t.Errorf("Wrong deposit: actual %v want %v", actual, expected)
	}

	actualSuccess := Withdraw(250)
	expectedSuccess := true
	if actualSuccess != expectedSuccess {
		t.Errorf("Must succeed: actual %v want %v", actual, expected)
	}

	actual = Balance()
	expected = 750
	if actual != expected {
		t.Errorf("Wrong withdraw: actual %v want %v", actual, expected)
	}

	actualSuccess = Withdraw(1000)
	expectedSuccess = false
	if actualSuccess != expectedSuccess {
		t.Errorf("Must fail: actual %v want %v", actual, expected)
	}

	actual = Balance()
	expected = 750
	if actual != expected {
		t.Errorf("Must not reduce: actual %v want %v", actual, expected)
	}
}
