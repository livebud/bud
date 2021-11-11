package bail

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

func (s *Struct) Bail(err error) {
	s.err = err
	panic(bail{})
}
