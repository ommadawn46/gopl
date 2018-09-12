package eval

type Expr interface {
	Eval(env Env) float64
	Check(vars map[Var]bool) error
	String() string
}

type Var string

type literal float64

type unary struct {
	op rune // + -
	x  Expr
}

type binary struct {
	op   rune // + - * /
	x, y Expr
}

var FuncNumParams = map[string]int{"max": 2, "min": 2, "add": 2, "mul": 2, "pow": 2, "sin": 1, "sqrt": 1}

// 2 args funcs: max, min, add, mul, pow
// 1 arg funcs: sin, sqrt
type call struct {
	fn   string // any func
	args []Expr
}

type reduce struct {
	fn   string // 2 args func
	args []Expr
}

func ExtractAllVar(expr Expr) (vars []Var) {
	switch e := expr.(type) {
	case Var:
		vars = append(vars, e)
	case unary:
		vars = append(vars, ExtractAllVar(e.x)...)
	case binary:
		vars = append(vars, ExtractAllVar(e.x)...)
		vars = append(vars, ExtractAllVar(e.y)...)
	case call:
		for _, e_ := range e.args {
			vars = append(vars, ExtractAllVar(e_)...)
		}
	case reduce:
		for _, e_ := range e.args {
			vars = append(vars, ExtractAllVar(e_)...)
		}
	}
	return
}
