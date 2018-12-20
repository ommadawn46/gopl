package sexpr

import (
	"fmt"
	"io"
	"strconv"
	"text/scanner"
)

type Token interface{}

type Symbol string

type String string

type Int int

type Nil struct{}

func (_ Nil) String() string {
	return "nil"
}

type StartList struct{}

func (_ StartList) String() string {
	return "StartList"
}

type EndList struct{}

func (_ EndList) String() string {
	return "EndList"
}

type Decoder struct {
	r   io.Reader
	lex *lexer
}

func NewDecoder(r io.Reader) *Decoder {
	lex := &lexer{scan: scanner.Scanner{Mode: scanner.GoTokens}}
	lex.scan.Init(r)
	lex.next() // get the first token
	return &Decoder{r: r, lex: lex}
}

type lexer struct {
	scan  scanner.Scanner
	token rune // the current token
}

func (lex *lexer) next()        { lex.token = lex.scan.Scan() }
func (lex *lexer) text() string { return lex.scan.TokenText() }

func (d *Decoder) Token() (token Token, err error) {
	defer func() {
		// NOTE: this is not an example of ideal error handling.
		if x := recover(); x != nil {
			token = nil
			err = fmt.Errorf("error at %s: %v", d.lex.scan.Position, x)
		}
	}()
	switch d.lex.token {
	case scanner.Ident:
		// The only valid identifiers are
		// "nil" and struct field names.
		if d.lex.text() == "nil" {
			d.lex.next()
			return Nil{}, nil
		}
	case scanner.String:
		s, _ := strconv.Unquote(d.lex.text()) // NOTE: ignoring errors
		d.lex.next()
		return String(s), nil
	case scanner.Int:
		i, _ := strconv.Atoi(d.lex.text()) // NOTE: ignoring errors
		d.lex.next()
		return Int(i), nil
	case '(':
		d.lex.next()
		return StartList{}, nil
	case ')':
		d.lex.next()
		return EndList{}, nil
	case scanner.EOF:
		return nil, io.EOF
	}
	panic(fmt.Sprintf("unexpected token %q", d.lex.text()))
}
