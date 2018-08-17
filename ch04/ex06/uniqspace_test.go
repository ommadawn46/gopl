package uniqspace

import (
	"reflect"
	"testing"
)

func TestUniqspaceEmpty(t *testing.T) {
	actual := uniqspace([]byte(""))
	expected := []byte("")
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("actual %v want %v", actual, expected)
	}
}

func TestUniqspaceAscii(t *testing.T) {
	actual := uniqspace([]byte("A B  C   D    E     F      G       H"))
	expected := []byte("A B C D E F G H")
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("actual %s want %s", actual, expected)
	}
}

func TestUniqspaceUtf8(t *testing.T) {
	actual := uniqspace([]byte("あ  い　　う 　 え　 　お"))
	expected := []byte("あ い う え お")
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("actual %s want %s", actual, expected)
	}
}
