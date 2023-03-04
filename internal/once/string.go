package once

import "sync"

type onceString struct {
	o sync.Once
	s string
	e error
}

func (s *onceString) Do(fn func() (string, error)) (string, error) {
	s.o.Do(func() { s.s, s.e = fn() })
	return s.s, s.e
}

func String(fn func() (string, error)) func() (string, error) {
	s := new(onceString)
	return func() (string, error) {
		return s.Do(fn)
	}
}
