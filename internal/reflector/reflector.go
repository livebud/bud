package reflector

import (
	"fmt"
	"path"
	"reflect"
	"strings"
)

var ErrNoName = fmt.Errorf("reflector: no name")

func TypeOf[V any](v V) (string, error) {
	t := reflect.TypeOf(v)
	if t == nil {
		t = reflect.TypeOf(new(V)).Elem()
	}
	return Type(t)
}

func Type(t reflect.Type) (string, error) {
	prefix, inner, err := innermost("", t)
	if err != nil {
		return "", err
	}
	if inner.Name() == "" {
		return "", fmt.Errorf("%w for %s", ErrNoName, t)
	}
	key, err := toString(prefix, inner)
	if err != nil {
		return "", err
	}
	return key, nil
}

func innermost(prefix string, t reflect.Type) (string, reflect.Type, error) {
	switch t.Kind() {
	case reflect.Ptr:
		return innermost(prefix+"*", t.Elem())
	case reflect.Slice:
		return innermost(prefix+"[]", t.Elem())
	case reflect.Map:
		key, innerKey, err := innermost("", t.Key())
		if err != nil {
			return "", nil, err
		}
		key, err = toString(key, innerKey)
		if err != nil {
			return "", nil, err
		}
		return innermost(prefix+"map["+key+"]", t.Elem())
	default:
		return prefix, t, nil
	}
}

func toString(prefix string, t reflect.Type) (string, error) {
	typeName := t.String()
	typeParts := strings.SplitN(typeName, ".", 2)
	dir := ""
	if len(typeParts) == 2 {
		dir = typeParts[0]
		typeName = typeParts[1]
	}
	pkgPath := t.PkgPath()
	if pkgPath == "" {
		return "", fmt.Errorf("%w for %s", ErrNoName, t)
	} else if strings.HasPrefix(pkgPath, "/") {
		return path.Dir(pkgPath) + "/" + dir + "." + prefix + typeName, nil
	}
	return pkgPath + "." + prefix + typeName, nil
}
