package eval

import (
	"math"
	"testing"
)

func TestString(t *testing.T) {
	tests := []struct {
		exprStr string
		env     Env
	}{
		{"sqrt(A / pi)", Env{"A": 87616, "pi": math.Pi}},
		{"pow(x, 3) + pow(y, 3)", Env{"x": 12, "y": 1}},
		{"pow(x, 3) + pow(y, 3)", Env{"x": 9, "y": 10}},
		{"5 / 9 * (F - 32)", Env{"F": -40}},
		{"5 / 9 * (F - 32)", Env{"F": 32}},
		{"5 / 9 * (F - 32)", Env{"F": 212}},
		{"-1 + -x", Env{"x": 1}},
		{"-1 - x", Env{"x": 1}},
	}
	for _, test := range tests {
		expr1, err := Parse(test.exprStr)
		expected := expr1.Eval(test.env)
		if err != nil {
			t.Error(err)
			continue
		}

		expr2, err := Parse(expr1.String())
		actual := expr2.Eval(test.env)
		if err != nil {
			t.Error(err)
			continue
		}

		if actual != expected {
			t.Errorf("actual %v want %v", actual, test.exprStr)
		}
	}
}
