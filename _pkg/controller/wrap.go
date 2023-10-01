package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
)

var zeroValue reflect.Value

var ErrInvalidHandler = errors.New("invalid handler")

func invalidHandler(t reflect.Type) error {
	return fmt.Errorf(`rpc: "%v" is an %w type`, t, ErrInvalidHandler)
}

var contextContext = reflect.TypeOf((*context.Context)(nil)).Elem()

func isContext(t reflect.Type) bool {
	return t.Implements(contextContext)
}

func contextArgFunc(w http.ResponseWriter, r *http.Request) (reflect.Value, error) {
	return reflect.ValueOf(r.Context()), nil
}

var errorType = reflect.TypeOf((*error)(nil)).Elem()

func isError(t reflect.Type) bool {
	return t.Implements(errorType)
}

type statuser interface {
	Status() int
}

type badRequestError struct {
	err error
}

func (b *badRequestError) Error() string {
	return b.err.Error()
}

func (b *badRequestError) Status() int {
	return http.StatusBadRequest
}

func (b *badRequestError) Unwrap() error {
	return b.err
}

func isStruct(t reflect.Type) bool {
	return t.Kind() == reflect.Struct
}

func isPointer(t reflect.Type) bool {
	return t.Kind() == reflect.Pointer
}

var responseWriterType = reflect.TypeOf((*http.ResponseWriter)(nil)).Elem()

func isResponseWriter(t reflect.Type) bool {
	return t.Implements(responseWriterType)
}

var requestType = reflect.TypeOf((*http.Request)(nil))

func isRequest(t reflect.Type) bool {
	return t == requestType
}

// Wrap a handler for use as an http.Handler
func wrapValue(reader reader, writer writer, fn reflect.Value, preArgs ...reflect.Value) (h http.Handler, err error) {
	fnType := fn.Type()
	if fn.Kind() != reflect.Func {
		return nil, invalidHandler(fnType)
	}
	// Prebuild the list of arguments
	argFuncs, err := prebuildArgs(reader, fnType, len(preArgs))
	if err != nil {
		return nil, err
	}
	resultFunc, err := prebuildResult(writer, fnType)
	if err != nil {
		return nil, err
	}
	// Return the handler
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lenPreArgs := len(preArgs)
		args := make([]reflect.Value, lenPreArgs+len(argFuncs))
		copy(args, preArgs)
		for i, f := range argFuncs {
			arg, err := f(w, r)
			if err != nil {
				writer.WriteError(w, r, &badRequestError{err})
				return
			}
			args[lenPreArgs+i] = arg
		}
		resultFunc(w, r, fn.Call(args))
	}), nil
}

type argFunc = func(http.ResponseWriter, *http.Request) (reflect.Value, error)

func isContextStruct(t reflect.Type) bool {
	return isStruct(t) && isContext(t)
}

func prebuildArgs(reader reader, fnType reflect.Type, offset int) (fns []argFunc, err error) {
	numArgs := fnType.NumIn() - offset
	fns = make([]argFunc, numArgs)
	switch numArgs {
	case 0:
		return fns, nil
	case 1:
		fns[0], err = firstInput(reader, fnType, fnType.In(offset+0))
		if err != nil {
			return nil, err
		}
	case 2:
		firstArg := fnType.In(offset + 0)
		fns[0], err = firstInput(reader, fnType, firstArg)
		if err != nil {
			return nil, err
		}
		fns[1], err = secondInput(reader, fnType, firstArg, fnType.In(offset+1))
		if err != nil {
			return nil, err
		}
	default:
		return nil, invalidHandler(fnType)
	}
	return fns, nil
}

func firstInput(reader reader, fnType reflect.Type, firstType reflect.Type) (argFunc, error) {
	switch {
	case isPointer(firstType):
		argFunc, err := firstInput(reader, fnType, firstType.Elem())
		if err != nil {
			return nil, err
		}
		return func(w http.ResponseWriter, r *http.Request) (reflect.Value, error) {
			value, err := argFunc(w, r)
			if err != nil {
				return zeroValue, err
			}
			return value.Addr(), nil
		}, nil
	case isContextStruct(firstType):
		return func(w http.ResponseWriter, r *http.Request) (reflect.Value, error) {
			return reader.ReadContext(r, firstType)
		}, nil
	case isStruct(firstType):
		return func(w http.ResponseWriter, r *http.Request) (reflect.Value, error) {
			return reader.ReadInput(r, firstType)
		}, nil
	case isContext(firstType):
		return contextArgFunc, nil
	case isResponseWriter(firstType):
		return func(w http.ResponseWriter, r *http.Request) (reflect.Value, error) {
			return reflect.ValueOf(w), nil
		}, nil
	case isRequest(firstType):
		return func(w http.ResponseWriter, r *http.Request) (reflect.Value, error) {
			return reflect.ValueOf(r), nil
		}, nil
	default:
		return nil, invalidHandler(fnType)
	}
}

func secondInput(reader reader, fnType reflect.Type, firstType reflect.Type, secondType reflect.Type) (argFunc, error) {
	switch {
	case isRequest(secondType):
		if !isResponseWriter(firstType) && !isContext(firstType) {
			return nil, invalidHandler(fnType)
		}
		return func(w http.ResponseWriter, r *http.Request) (reflect.Value, error) {
			return reflect.ValueOf(r), nil
		}, nil
	case isPointer(secondType):
		argFunc, err := secondInput(reader, fnType, firstType, secondType.Elem())
		if err != nil {
			return nil, err
		}
		return func(w http.ResponseWriter, r *http.Request) (reflect.Value, error) {
			value, err := argFunc(w, r)
			if err != nil {
				return zeroValue, err
			}
			return value.Addr(), nil
		}, nil
	case isResponseWriter(secondType):
		if !isContext(firstType) {
			return nil, invalidHandler(fnType)
		}
		return func(w http.ResponseWriter, r *http.Request) (reflect.Value, error) {
			return reflect.ValueOf(w), nil
		}, nil
	case isContext(secondType): // Context should always be the first argument
		return nil, invalidHandler(fnType)
	case isStruct(secondType):
		if !isContext(firstType) {
			return nil, invalidHandler(fnType)
		}
		return func(w http.ResponseWriter, r *http.Request) (reflect.Value, error) {
			return reader.ReadInput(r, secondType)
		}, nil
	default:
		return nil, invalidHandler(fnType)
	}
}

type resultFunc = func(http.ResponseWriter, *http.Request, []reflect.Value)

func prebuildResult(writer writer, fnType reflect.Type) (resultFunc, error) {
	numResults := fnType.NumOut()
	switch numResults {
	case 0:
		// Write an empty response
		return func(w http.ResponseWriter, r *http.Request, values []reflect.Value) {
			writer.WriteEmpty(w, r)
		}, nil
	case 1:
		return firstResult(writer, fnType, fnType.Out(0))
	case 2:
		return secondResult(writer, fnType, fnType.Out(0), fnType.Out(1))
	default:
		return nil, invalidHandler(fnType)
	}
}

func firstResult(writer writer, fnType reflect.Type, firstType reflect.Type) (resultFunc, error) {
	switch {
	case isPointer(firstType):
		return firstResult(writer, fnType, firstType.Elem())
	case isStruct(firstType):
		return func(w http.ResponseWriter, r *http.Request, values []reflect.Value) {
			writer.WriteOutput(w, r, values[0].Interface())
		}, nil
	case isError(firstType):
		return func(w http.ResponseWriter, r *http.Request, values []reflect.Value) {
			firstValue := values[0].Interface()
			if err, ok := firstValue.(error); ok && err != nil {
				writer.WriteError(w, r, err)
				return
			}
			// Write an empty response
			writer.WriteEmpty(w, r)
		}, nil
	default:
		return nil, invalidHandler(fnType)
	}
}

func secondResult(writer writer, fnType reflect.Type, firstType, secondType reflect.Type) (resultFunc, error) {
	switch {
	case isError(secondType):
		if isError(firstType) {
			return nil, invalidHandler(fnType)
		}
		return func(w http.ResponseWriter, r *http.Request, values []reflect.Value) {
			secondValue := values[1].Interface()
			if err, ok := secondValue.(error); ok && err != nil {
				writer.WriteError(w, r, err)
				return
			}
			firstValue := values[0].Interface()
			writer.WriteOutput(w, r, firstValue)
		}, nil
	default:
		return nil, invalidHandler(fnType)
	}
}
