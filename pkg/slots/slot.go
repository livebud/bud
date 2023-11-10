package slots

import (
	"bytes"
	"io"
	"sync"
)

func New() *Slots {
	return &Slots{
		reader: bytes.NewBuffer(nil),
		inners: make(map[string]*bytes.Buffer),
		pipe:   newPipe(),
	}
}

type Slots struct {
	reader io.Reader
	pipe   *pipe
	mu     sync.Mutex
	inners map[string]*bytes.Buffer
}

var _ io.ReadWriteCloser = (*Slots)(nil)

func (s *Slots) Read(p []byte) (n int, err error) {
	return s.reader.Read(p)
}

func (s *Slots) Write(p []byte) (n int, err error) {
	return s.pipe.Write(p)
}

func (s *Slots) Close() error {
	return s.pipe.Close()
}

func (s *Slots) Next() *Slots {
	pipe := newPipe()
	slots := &Slots{
		reader: s.pipe,
		pipe:   pipe,
		inners: s.inners,
	}
	return slots
}

type Slot interface {
	io.ReadWriter
}

func (s *Slots) Slot(name string) Slot {
	s.mu.Lock()
	defer s.mu.Unlock()
	inner, ok := s.inners[name]
	if !ok {
		inner := new(bytes.Buffer)
		s.inners[name] = inner
		return inner
	}
	return inner
}
