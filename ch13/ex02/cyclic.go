package cyclic

import (
	"reflect"
	"unsafe"
)

type typePointer struct {
	p unsafe.Pointer
	t reflect.Type
}

func cyclic(v reflect.Value, seen map[typePointer]bool) bool {
	if v.CanAddr() {
		vptr := unsafe.Pointer(v.UnsafeAddr())
		c := typePointer{vptr, v.Type()}
		if seen[c] {
			return true // already seen
		}
		seen[c] = true
	}

	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		if cyclic(v.Elem(), seen) {
			return true
		}
	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			if cyclic(v.Index(i), seen) {
				return true
			}
		}
	case reflect.Struct:
		for i, n := 0, v.NumField(); i < n; i++ {
			if cyclic(v.Field(i), seen) {
				return true
			}
		}
	case reflect.Map:
		for _, k := range v.MapKeys() {
			if cyclic(v.MapIndex(k), seen) {
				return true
			}
		}
	}
	return false
}

func Cyclic(v interface{}) bool {
	seen := make(map[typePointer]bool)
	return cyclic(reflect.ValueOf(v), seen)
}
