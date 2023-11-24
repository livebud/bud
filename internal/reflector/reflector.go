package reflector

import (
	"fmt"
	"reflect"
)

var ErrNoName = fmt.Errorf("reflector: no name")

// TypeOf takes a generic and returns a string representation of its type.
func TypeOf[V any](v V) (string, error) {
	t := reflect.TypeOf(v)
	if t == nil {
		t = reflect.TypeOf(new(V)).Elem()
	}
	return Type(t)
}

// Type takes a reflect.Type and returns a string representation of it.
func Type(t reflect.Type) (string, error) {
	prefix := ""
	if t.Kind() == reflect.Ptr {
		prefix = "*"
		t = t.Elem()
	}
	pkgPath := t.PkgPath()
	name := t.Name()
	if pkgPath != "" && name != "" {
		return pkgPath + "." + prefix + name, nil
	}
	return "", fmt.Errorf("%w for %s", ErrNoName, t)
}
