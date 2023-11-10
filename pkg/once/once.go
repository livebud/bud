package once

import "sync"

// Func creates a function that will only be called once. It is safe for
// concurrent use. Subsequent calls to Do will return the same value and error
// as the first call.
func Func[V any](fn func() (V, error)) func() (V, error) {
	o := new(Once[V])
	return func() (V, error) {
		return o.Do(fn)
	}
}

// Once can be used to only call a function once. It is safe for concurrent use.
// Subsequent calls to Do will return the same value and error as the first
// call.
type Once[V any] struct {
	o sync.Once
	v V
	e error
}

func (o *Once[V]) Do(fn func() (V, error)) (V, error) {
	o.o.Do(func() { o.v, o.e = fn() })
	return o.v, o.e
}

// Bytes is a helper function for Once[[]byte].
type Bytes = Once[[]byte]

// String is a helper function for Once[string].
type String = Once[string]
