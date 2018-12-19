package display

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var DEPTH_LIMIT = 8

func Display(name string, x interface{}) {
	fmt.Printf("Display %s (%T):\n", name, x)
	display(name, reflect.ValueOf(x), 0)
}

func formatAtom(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Invalid:
		return "invalid"
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Bool:
		if v.Bool() {
			return "true"
		}
		return "false"
	case reflect.String:
		return strconv.Quote(v.String())
	case reflect.Chan, reflect.Func, reflect.Ptr,
		reflect.Slice, reflect.Map:
		return v.Type().String() + " 0x" +
			strconv.FormatUint(uint64(v.Pointer()), 16)
	case reflect.Array:
		elems := []string{}
		for i := 0; i < v.Len(); i++ {
			elems = append(elems, formatAtom(v.Index(i)))
		}
		return fmt.Sprintf(
			"%s{%s}", v.Type().String(), strings.Join(elems, ", "),
		)
	case reflect.Struct:
		fields := []string{}
		for i := 0; i < v.NumField(); i++ {
			fieldName := v.Type().Field(i).Name
			fields = append(
				fields,
				fmt.Sprintf("%s: %s", fieldName, formatAtom(v.FieldByName(fieldName))),
			)
		}
		return fmt.Sprintf(
			"%s{%s}", v.Type().String(), strings.Join(fields, ", "),
		)
	default:
		return v.Type().String() + " value"
	}
}

func display(path string, v reflect.Value, depth int) {
	if depth >= DEPTH_LIMIT {
		fmt.Printf("%s = stuck\n", path)
		return
	}
	switch v.Kind() {
	case reflect.Invalid:
		fmt.Printf("%s = invalid\n", path)
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			display(fmt.Sprintf("%s[%d]", path, i), v.Index(i), depth+1)
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			fieldPath := fmt.Sprintf("%s.%s", path, v.Type().Field(i).Name)
			display(fieldPath, v.Field(i), depth+1)
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			display(
				fmt.Sprintf(
					"%s[%s]", path, formatAtom(key),
				),
				v.MapIndex(key), depth+1,
			)
		}
	case reflect.Ptr:
		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
		} else {
			display(fmt.Sprintf("(*%s)", path), v.Elem(), depth+1)
		}
	case reflect.Interface:
		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
		} else {
			fmt.Printf("%s.type = %s\n", path, v.Elem().Type())
			display(path+".value", v.Elem(), depth+1)
		}
	default:
		fmt.Printf("%s = %s\n", path, formatAtom(v))
	}
}
