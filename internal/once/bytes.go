package once

import "sync"

type Bytes struct {
	o sync.Once
	v []byte
	e error
}

func (b *Bytes) Do(fn func() ([]byte, error)) ([]byte, error) {
	b.o.Do(func() { b.v, b.e = fn() })
	return b.v, b.e
}
