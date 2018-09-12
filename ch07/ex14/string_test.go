package eval

import (
	"testing"
)

func TestString(t *testing.T) {
	tests := []struct {
		exprStr string
		env     Env
	}{
		{"reduce(max, A, -100, B, 500)", Env{"A": -50, "B": 1000}},
		{"reduce(min, A, -100, B, 500, C)", Env{"A": -50, "B": 1000, "C": -600}},
		{"reduce(add, A, -100, B, 500, C, 200)", Env{"A": -50, "B": 1000, "C": -600}},
		{"reduce(mul, A, -1, B, 5, C, 2)", Env{"A": -5, "B": 10, "C": -6}},
		{"reduce(pow, 2, 2, 2, 2, 2)", nil},
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
