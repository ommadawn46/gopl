package eval

import (
	"fmt"
	"strings"
)

func (v Var) Check(vars map[Var]bool) error {
	vars[v] = true
	return nil
}

func (literal) Check(vars map[Var]bool) error {
	return nil
}

func (u unary) Check(vars map[Var]bool) error {
	if !strings.ContainsRune("+-", u.op) {
		return fmt.Errorf("unexpected unary op %q", u.op)
	}
	return u.x.Check(vars)
}

func (b binary) Check(vars map[Var]bool) error {
	if !strings.ContainsRune("+-*/", b.op) {
		return fmt.Errorf("unexpected binary op %q", b.op)
	}
	if err := b.x.Check(vars); err != nil {
		return err
	}
	return b.y.Check(vars)
}

func (c call) Check(vars map[Var]bool) error {
	arity, ok := FuncNumParams[c.fn]
	if !ok {
		return fmt.Errorf("unknown function %q", c.fn)
	}
	if len(c.args) != arity {
		return fmt.Errorf("call to %s has %d args, want %d",
			c.fn, len(c.args), arity)
	}
	for _, arg := range c.args {
		if err := arg.Check(vars); err != nil {
			return err
		}
	}
	return nil
}

func (r reduce) Check(vars map[Var]bool) error {
	arity, ok := FuncNumParams[r.fn]
	if !ok {
		return fmt.Errorf("unknown function %q", r.fn)
	}
	if arity < 2 {
		return fmt.Errorf("reduce only allows 2 args function %q", r.fn)
	}
	if len(r.args) < arity {
		return fmt.Errorf("call to %s has %d args, want %d",
			r.fn, len(r.args), arity)
	}
	for _, arg := range r.args {
		if err := arg.Check(vars); err != nil {
			return err
		}
	}
	return nil
}
