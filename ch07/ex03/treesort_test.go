package treesort

import (
	"testing"
)

func TestTreeString(test *testing.T) {
	var t *tree
	t = add(t, 10)
	t = add(t, -10)
	t = add(t, 100)
	t = add(t, 5)
	t = add(t, -1000)
	t = add(t, 10000)
	t = add(t, -50)

	actual := t.String()
	expected := "[-1000 -50 -10 5 10 100 10000]"
	if actual != expected {
		test.Errorf("actual %v want %v", actual, expected)
	}
}
