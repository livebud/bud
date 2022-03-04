package tracer

import (
	"context"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

func Start() (*Tracer, error) {
	tempDir, err := ioutil.TempDir("", "bud-tracer-*")
	if err != nil {
		return nil, err
	}
	socketPath := filepath.Join(tempDir, "trace.sock")
	server, err := Serve(socketPath)
	if err != nil {
		return nil, err
	}
	client, err := NewClient(socketPath)
	if err != nil {
		return nil, err
	}
	tracer := New(client)
	return &Tracer{
		t: tracer,
		s: server,
		c: client,
	}, nil
}

type Tracer struct {
	t Trace
	s *http.Server
	c *Client
}

func (t *Tracer) Start(ctx context.Context, label string, attrs ...interface{}) (context.Context, Span) {
	return t.t.Start(ctx, label, attrs...)
}

func (t *Tracer) Print(ctx context.Context) (string, error) {
	return t.c.Print(ctx)
}

func (t *Tracer) Shutdown(ctx context.Context) error {
	if err := t.c.Shutdown(ctx); err != nil {
		return err
	}
	if err := t.s.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

// Resume from subprocess
func Resume(ctx context.Context, addr string, data []byte) (context.Context, Trace, error) {
	client, err := NewClient(addr)
	if err != nil {
		return nil, nil, err
	}
	ctx, err = Decode(ctx, data)
	if err != nil {
		return nil, nil, err
	}
	tracer := New(client)
	return ctx, tracer, nil
}
