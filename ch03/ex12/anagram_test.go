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

func TestAnagramFalse1(t *testing.T) {
	actual := isAnagram("ABCDEFGH", "HIJKLMNO")
	expected := false
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}
}

func TestAnagramFalse2(t *testing.T) {
	actual := isAnagram("A", "ABBB")
	expected := false
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}
}

func TestAnagramDuplicateFalse1(t *testing.T) {
	actual := isAnagram("ABCDEFG", "ABCDEFGG")
	expected := false
	if actual != expected {
		t.Errorf("actual %v want %v", actual, expected)
	}
}
