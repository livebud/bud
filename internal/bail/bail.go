package bail

import "fmt"

type Struct struct {
	err error
}

type bail struct{}

func (s *Struct) Recover(err *error) {
	if e := recover(); e != nil {
		// resume same panic if it's not bailing
		if _, ok := e.(bail); !ok {
			panic(e)
		}
		*err = s.err
	}
}

func (s *Struct) Recover2(err *error, prefix string) {
	if e := recover(); e != nil {
		// resume same panic if it's not bailing
		if _, ok := e.(bail); !ok {
			panic(e)
		}
		*err = fmt.Errorf(prefix+". %w", s.err)
	}
}

func (s *Struct) Bail(err error) {
	s.err = err
	panic(bail{})
}
