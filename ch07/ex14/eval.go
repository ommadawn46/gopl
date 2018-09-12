package eval

import (
	"fmt"
	"math"
)

type Env map[Var]float64

func (v Var) Eval(env Env) float64 {
	return env[v]
}

func (l literal) Eval(_ Env) float64 {
	return float64(l)
}

func (u unary) Eval(env Env) float64 {
	switch u.op {
	case '+':
		return +u.x.Eval(env)
	case '-':
		return -u.x.Eval(env)
	}
	panic(fmt.Sprintf("unsupported unary operator: %q", u.op))
}

func (b binary) Eval(env Env) float64 {
	switch b.op {
	case '+':
		return b.x.Eval(env) + b.y.Eval(env)
	case '-':
		return b.x.Eval(env) - b.y.Eval(env)
	case '*':
		return b.x.Eval(env) * b.y.Eval(env)
	case '/':
		return b.x.Eval(env) / b.y.Eval(env)
	}
	panic(fmt.Sprintf("unsupported binary operator: %q", b.op))
}

func (c call) Eval(env Env) float64 {
	switch c.fn {
	case "max":
		x, y := c.args[0].Eval(env), c.args[1].Eval(env)
		if x > y {
			return x
		} else {
			return y
		}
	case "min":
		x, y := c.args[0].Eval(env), c.args[1].Eval(env)
		if x < y {
			return x
		} else {
			return y
		}
	case "add":
		return c.args[0].Eval(env) + c.args[1].Eval(env)
	case "mul":
		return c.args[0].Eval(env) * c.args[1].Eval(env)
	case "pow":
		return math.Pow(c.args[0].Eval(env), c.args[1].Eval(env))
	case "sin":
		return math.Sin(c.args[0].Eval(env))
	case "sqrt":
		return math.Sqrt(c.args[0].Eval(env))
	}
	panic(fmt.Sprintf("unsupported function call: %s", c.fn))
}

func (r reduce) Eval(env Env) float64 {
	result := r.args[0]
	for i := 1; i < len(r.args); i++ {
		x := r.args[i]
		c := call{r.fn, []Expr{result, x}}
		result = literal(c.Eval(env))
	}
	return result.Eval(env)
}
