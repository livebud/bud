package di

import (
	"context"
	"errors"
)

// ErrNotInContext is returned when the injector is not in the context
var ErrNotInContext = errors.New("di: injector is not in the context")

type contextKey string

var key contextKey = "di"

// WithInjector returns a new context that contains the injector
func WithInjector(parent context.Context, in Injector) context.Context {
	return context.WithValue(parent, key, in)
}

// FromContext returns the injector from a context if present
func FromContext(ctx context.Context) (Injector, error) {
	in, ok := ctx.Value(key).(Injector)
	if !ok {
		return nil, ErrNotInContext
	}
	return in, nil
}

// FromContext returns the injector from a context if present
func LoadFrom[Dep any](ctx context.Context) (d Dep, err error) {
	in, ok := ctx.Value(key).(Injector)
	if !ok {
		return d, ErrNotInContext
	}
	return Load[Dep](in)
}
