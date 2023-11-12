// Package u provides utility functions.
package u

import "sync"

func Must[V any](v V, err error) V {
	if err != nil {
		panic(err)
	}
	return v
}

func Once[V any](fn func() (V, error)) func() (V, error) {
	o := new(once[V])
	return func() (V, error) {
		return o.Do(fn)
	}
}

// Once can be used to only call a function once. It is safe for concurrent use.
// Subsequent calls to Do will return the same value and error as the first
// call.
type once[V any] struct {
	o sync.Once
	v V
	e error
}

func (o *once[V]) Do(fn func() (V, error)) (V, error) {
	o.o.Do(func() { o.v, o.e = fn() })
	return o.v, o.e
}
