package sexpr

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"text/scanner"
)

var INTERFACES map[string]reflect.Type

type Decoder struct {
	r io.Reader
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r}
}

func (d *Decoder) Decode(v interface{}) (err error) {
	lex := &lexer{scan: scanner.Scanner{Mode: scanner.GoTokens}}
	lex.scan.Init(d.r)
	lex.next() // get the first token
	defer func() {
		// NOTE: this is not an example of ideal error handling.
		//if x := recover(); x != nil {
		//	err = fmt.Errorf("error at %s: %v", lex.scan.Position, x)
		//}
	}()
	read(lex, reflect.ValueOf(v).Elem())
	return nil
}

func Unmarshal(data []byte, out interface{}) (err error) {
	r := bytes.NewBuffer(data)
	decoder := NewDecoder(r)
	return decoder.Decode(out)
}

type lexer struct {
	scan  scanner.Scanner
	token rune // the current token
}

func (lex *lexer) next()        { lex.token = lex.scan.Scan() }
func (lex *lexer) text() string { return lex.scan.TokenText() }

func (lex *lexer) consume(want rune) {
	if lex.token != want { // NOTE: Not an example of good error handling.
		panic(fmt.Sprintf("got %q, want %q", lex.text(), want))
	}
	lex.next()
}

func read(lex *lexer, v reflect.Value) {
	switch lex.token {
	case scanner.Ident:
		// The only valid identifiers are
		// "nil" and struct field names.
		switch lex.text() {
		case "nil":
			v.Set(reflect.Zero(v.Type()))
			lex.next()
			return
		case "t":
			v.SetBool(true)
			lex.next()
			return
		}
	case scanner.String:
		s, _ := strconv.Unquote(lex.text()) // NOTE: ignoring errors
		v.SetString(s)
		lex.next()
		return
	case scanner.Int:
		i, _ := strconv.Atoi(lex.text()) // NOTE: ignoring errors
		v.SetInt(int64(i))
		lex.next()
		return
	case scanner.Float:
		f, _ := strconv.ParseFloat(lex.text(), 64) // NOTE: ignoring errors
		v.SetFloat(float64(f))
		lex.next()
		return
	case '#':
		lex.next()
		if lex.text() != "C" {
			break
		}
		lex.next()
		if lex.text() != "(" {
			break
		}
		lex.next()
		f1, _ := strconv.ParseFloat(lex.text(), 64)
		lex.next()
		f2, _ := strconv.ParseFloat(lex.text(), 64)
		v.SetComplex(complex(f1, f2))
		lex.next()
		if lex.text() != ")" {
			break
		}
		lex.next()
		return

	case '(':
		lex.next()
		readList(lex, v)
		lex.next() // consume ')'
		return
	}
	panic(fmt.Sprintf("unexpected token %q", lex.text()))
}

func readList(lex *lexer, v reflect.Value) {
	switch v.Kind() {
	case reflect.Array: // (item ...)
		for i := 0; !endList(lex); i++ {
			read(lex, v.Index(i))
		}

	case reflect.Slice: // (item ...)
		for !endList(lex) {
			item := reflect.New(v.Type().Elem()).Elem()
			read(lex, item)
			v.Set(reflect.Append(v, item))
		}

	case reflect.Struct: // ((name value) ...)
		fieldNames := make(map[string]string)
		for i := 0; i < v.NumField(); i++ {
			fieldInfo := v.Type().Field(i) // a reflect.StructField
			tag := fieldInfo.Tag           // a reflect.StructTag
			name := tag.Get("sexpr")
			if name == "" {
				name = fieldInfo.Name
			}
			fieldNames[name] = fieldInfo.Name
		}
		for !endList(lex) {
			lex.consume('(')
			if lex.token != scanner.Ident {
				panic(fmt.Sprintf("got token %q, want field name", lex.text()))
			}
			name := fieldNames[lex.text()]
			lex.next()

			read(lex, v.FieldByName(name))
			lex.consume(')')
		}

	case reflect.Map: // ((key value) ...)
		v.Set(reflect.MakeMap(v.Type()))
		for !endList(lex) {
			lex.consume('(')
			key := reflect.New(v.Type().Key()).Elem()
			read(lex, key)
			value := reflect.New(v.Type().Elem()).Elem()
			read(lex, value)
			v.SetMapIndex(key, value)
			lex.consume(')')
		}

	case reflect.Interface:
		typeStr, _ := strconv.Unquote(lex.text())
		lex.next()
		reflectType, ok := INTERFACES[typeStr]
		if !ok {
			panic(fmt.Sprintf("unknown type interface %v", typeStr))
		}
		value := reflect.New(reflectType).Elem()
		read(lex, value)
		v.Set(value)

	default:
		panic(fmt.Sprintf("cannot decode list into %v", v.Type()))
	}
}

func endList(lex *lexer) bool {
	switch lex.token {
	case scanner.EOF:
		panic("end of file")
	case ')':
		return true
	}
	return false
}

func init() {
	INTERFACES = map[string]reflect.Type{
		"[]int":        reflect.TypeOf([]int{}),
		"[]string":     reflect.TypeOf([]string{}),
		"[]float64":    reflect.TypeOf([]float64{}),
		"[]complex128": reflect.TypeOf([]complex128{}),
	}
}
