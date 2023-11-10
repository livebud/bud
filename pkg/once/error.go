package once

import "sync"

// Error can be used to only call a function once. Unlike Func, it only returns
// an error, no results. It is safe for concurrent use. Subsequent calls to Do
// will return the same error as the first call.
type Error struct {
	o   sync.Once
	err error
}

func (e *Error) Do(fn func() error) (err error) {
	e.o.Do(func() { e.err = fn() })
	return e.err
}
