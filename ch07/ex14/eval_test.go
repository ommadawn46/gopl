package eval

import (
	"fmt"
	"testing"
)

func TestEval(t *testing.T) {
	tests := []struct {
		expr string
		env  Env
		want string
	}{
		{"reduce(max, A, -100, B, 500)", Env{"A": -50, "B": 1000}, "1000"},
		{"reduce(min, A, -100, B, 500, C)", Env{"A": -50, "B": 1000, "C": -600}, "-600"},
		{"reduce(add, A, -100, B, 500, C, 200)", Env{"A": -50, "B": 1000, "C": -600}, "950"},
		{"reduce(mul, A, -1, B, 5, C, 2)", Env{"A": -5, "B": 10, "C": -6}, "-3000"},
		{"reduce(pow, 2, 2, 2, 2, 2)", nil, "65536"},
	}
	var prevExpr string
	for _, test := range tests {
		if test.expr != prevExpr {
			fmt.Printf("\n%s\n", test.expr)
			prevExpr = test.expr
		}
		expr, err := Parse(test.expr)
		if err != nil {
			t.Error(err)
			continue
		}
		got := fmt.Sprintf("%.6g", expr.Eval(test.env))
		fmt.Printf("\t%v => %s\n", test.env, got)
		if got != test.want {
			t.Errorf("%s.Eval() in %v = %q, want %q\n",
				test.expr, test.env, got, test.want)
		}
	}
}

func TestErrors(t *testing.T) {
	var tests = []struct {
		input string
		env   Env
		want  string
	}{
		{"recude()", nil, "unknown function \"recude\""},
		{"recude(1, 2, 3, 4, 5)", nil, "unknown function \"recude\""},
	}

	for _, test := range tests {
		expr, err := Parse(test.input)
		if err == nil {
			err = expr.Check(map[Var]bool{})
		}
		if err != nil {
			if err.Error() != test.want {
				t.Errorf("%s: got %q, want %q", test.input, err, test.want)
			}
			continue
		}

		got := fmt.Sprintf("%.6g", expr.Eval(test.env))
		if got != test.want {
			t.Errorf("%s: %v => %s, want %s",
				test.input, test.env, got, test.want)
		}
	}
}
