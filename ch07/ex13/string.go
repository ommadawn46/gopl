package eval

import (
	"fmt"
)

func (v Var) String() string {
	return string(v)
}

func (l literal) String() string {
	return fmt.Sprintf("%g", l)
}

func (u unary) String() string {
	return fmt.Sprintf("%c%s", u.op, u.x.String())
}

func (b binary) String() string {
	return fmt.Sprintf("(%s %s %s)", b.x.String(), string(b.op), b.y.String())
}

func (c call) String() string {
	s := c.fn + "("
	for i, a := range c.args {
		if i != 0 {
			s += ", "
		}
		s += a.String()
	}
	return s + ")"
}
