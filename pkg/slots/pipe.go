package slots

import (
	"bytes"
)

func newPipe() *pipe {
	return &pipe{
		b: new(bytes.Buffer),
		c: make(chan struct{}),
	}
}

type pipe struct {
	b *bytes.Buffer
	c chan struct{}
}

func (p *pipe) Write(b []byte) (n int, err error) {
	return p.b.Write(b)
}

func (p *pipe) Close() error {
	close(p.c)
	return nil
}

func (p *pipe) Read(b []byte) (n int, err error) {
	<-p.c
	return p.b.Read(b)
}

func (p *pipe) Wait() {
	<-p.c
}
