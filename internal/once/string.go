package once

import "sync"

type String struct {
	o sync.Once
	s string
	e error
}

func (s *String) Do(fn func() (string, error)) (string, error) {
	s.o.Do(func() { s.s, s.e = fn() })
	return s.s, s.e
}
