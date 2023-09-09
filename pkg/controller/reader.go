package controller

import (
	"net/http"
	"reflect"

	"github.com/livebud/bud/pkg/controller/internal/request"
	"github.com/livebud/bud/pkg/di"
)

type defaultReader struct {
}

var _ reader = (*defaultReader)(nil)

func (defaultReader) ReadContext(r *http.Request, typ reflect.Type) (reflect.Value, error) {
	val := reflect.New(typ)
	in := val.Interface()
	injector, err := di.FromContext(r.Context())
	if err != nil {
		return zeroValue, err
	}
	if err := di.Unmarshal(injector, in); err != nil {
		return zeroValue, err
	}
	return val.Elem(), nil
}

func (defaultReader) ReadInput(r *http.Request, t reflect.Type) (reflect.Value, error) {
	value := reflect.New(t)
	in := value.Interface()
	if err := request.Unmarshal(r, in); err != nil {
		return zeroValue, err
	}
	return value.Elem(), nil
}
