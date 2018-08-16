package anagram

import (
	"testing"
)

func TestAnagramEmpty(t *testing.T) {
	actual := isAnagram("", "")
	expected := true
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}
}

func TestAnagramTrue(t *testing.T) {
	actual := isAnagram("Statue of Liberty", "Built to stay free")
	expected := true
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}
}

func TestAnagramFalse(t *testing.T) {
	actual := isAnagram("ABCDEFGH", "HIJKLMNO")
	expected := false
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}
}

func TestAnagramDuplicateFalse(t *testing.T) {
	actual := isAnagram("ABCDEFG", "ABCDEFGG")
	expected := false
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}
}
