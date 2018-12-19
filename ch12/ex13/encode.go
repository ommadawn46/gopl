package sexpr

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
)

type Encoder struct {
	w io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

func (enc *Encoder) Encode(v interface{}) error {
	data, err := Marshal(v)
	if err != nil {
		return err
	}
	_, err = enc.w.Write(data)
	return err
}

func Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := encode(&buf, reflect.ValueOf(v), 0); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func encode(buf *bytes.Buffer, v reflect.Value, indent int) error {
	switch v.Kind() {
	case reflect.Invalid:
		buf.WriteString("nil")

	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		fmt.Fprintf(buf, "%d", v.Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		fmt.Fprintf(buf, "%d", v.Uint())

	case reflect.String:
		fmt.Fprintf(buf, "%q", v.String())

	case reflect.Ptr:
		return encode(buf, v.Elem(), indent)

	case reflect.Array, reflect.Slice: // (value ...)
		buf.WriteByte('(')
		indent += 1
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				fmt.Fprintf(buf, "\n%*s", indent, "")
			}
			if err := encode(buf, v.Index(i), indent); err != nil {
				return err
			}
		}
		buf.WriteByte(')')

	case reflect.Struct: // ((name value) ...)
		buf.WriteByte('(')
		indent += 1
		firstField := true
		for i := 0; i < v.NumField(); i++ {
			fieldInfo := v.Type().Field(i) // a reflect.StructField
			tag := fieldInfo.Tag           // a reflect.StructTag
			name := tag.Get("sexpr")
			if name == "" {
				name = fieldInfo.Name
			}

			field := v.Field(i)
			if reflect.DeepEqual(reflect.Zero(field.Type()).Interface(), field.Interface()) {
				continue
			}
			if !firstField {
				fmt.Fprintf(buf, "\n%*s", indent, "")
			}
			outLen := buf.Len()
			fmt.Fprintf(buf, "(%s ", name)
			outLen = buf.Len() - outLen
			if err := encode(buf, field, indent+outLen); err != nil {
				return err
			}
			buf.WriteByte(')')
			firstField = false
		}
		buf.WriteByte(')')

	case reflect.Map: // ((key value) ...)
		buf.WriteByte('(')
		indent += 1
		for i, key := range v.MapKeys() {
			if i > 0 {
				fmt.Fprintf(buf, "\n%*s", indent, "")
			}
			buf.WriteByte('(')
			if err := encode(buf, key, indent); err != nil {
				return err
			}
			buf.WriteByte(' ')
			if err := encode(buf, v.MapIndex(key), indent); err != nil {
				return err
			}
			buf.WriteByte(')')
		}
		buf.WriteByte(')')

	case reflect.Bool:
		if v.Bool() {
			buf.WriteString("t")
		} else {
			buf.WriteString("nil")
		}

	case reflect.Float32, reflect.Float64:
		fmt.Fprintf(buf, "%g", v.Float())

	case reflect.Complex64, reflect.Complex128:
		buf.WriteString("#C(")
		cmplx := v.Complex()
		fmt.Fprintf(buf, "%g %g", real(cmplx), imag(cmplx))
		buf.WriteByte(')')

	case reflect.Interface:
		outLen := buf.Len()
		if !v.IsNil() {
			fmt.Fprintf(buf, "(%q ", v.Elem().Type().String())
			outLen = buf.Len() - outLen
			if err := encode(buf, v.Elem(), indent+outLen); err != nil {
				return err
			}
			buf.WriteByte(')')
		} else {
			buf.WriteString("nil")
		}

	default: // chan, func
		return fmt.Errorf("unsupported type: %s", v.Type())
	}
	return nil
}
