package once

import "sync"

type Error struct {
	o   sync.Once
	err error
}

func (e *Error) Do(fn func() error) (err error) {
	e.o.Do(func() { e.err = fn() })
	return e.err
}
