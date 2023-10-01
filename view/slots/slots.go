package slots

import (
	"bytes"
	"io"
	"sync"

	"github.com/livebud/bud/view"
)

func New(w io.Writer) *Slot {
	return &Slot{
		write: w.Write,
		read: func(p []byte) (n int, err error) {
			return 0, io.EOF
		},
		close: func(error) error {
			return nil
		},
	}
}

type Slot struct {
	mu    sync.RWMutex
	read  func(p []byte) (n int, err error)
	write func(p []byte) (n int, err error)
	close func(error) error
}

var _ view.Slot = (*Slot)(nil)

func (s *Slot) New() *Slot {
	pipe := newPipe()
	s.setReader(pipe)
	return &Slot{
		write: pipe.Write,
		read: func(p []byte) (n int, err error) {
			return 0, io.EOF
		},
		close: pipe.CloseWithError,
	}
}

func (s *Slot) setReader(r io.Reader) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.read = r.Read
}

func (s *Slot) Read(p []byte) (n int, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.read(p)
}

func (s *Slot) Write(p []byte) (n int, err error) {
	return s.write(p)
}

func (s *Slot) Close(err *error) error {
	if err != nil {
		return s.close(*err)
	}
	return s.close(nil)
}

func newPipe() *pipe {
	return &pipe{
		done: make(chan error, 1),
	}
}

type pipe struct {
	mu   sync.Mutex
	b    bytes.Buffer // Written data
	done chan error   // Writes are done
}

var _ io.ReadWriter = (*pipe)(nil)

func (p *pipe) Write(b []byte) (int, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.b.Write(b)
}

func (p *pipe) Read(b []byte) (int, error) {
	err := <-p.done
	if err != nil {
		return 0, err
	}
	return p.b.Read(b)
}

func (p *pipe) CloseWithError(err error) error {
	p.done <- err
	close(p.done)
	return nil
}
