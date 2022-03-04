package trace

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"gitlab.com/mnm/bud/internal/tracer"
)

type Trace = tracer.Trace
type Span = tracer.Span

var traceKey struct{}

type traceData struct {
	trace Trace
	path  string
}

func Serve(ctx context.Context) (*Tracer, context.Context, error) {
	tempDir, err := ioutil.TempDir("", "bud-trace-*")
	if err != nil {
		return nil, nil, err
	}
	path := filepath.Join(tempDir, "trace.sock")
	server, err := tracer.Serve(path)
	if err != nil {
		return nil, nil, err
	}
	client, err := tracer.NewClient(path)
	if err != nil {
		return nil, nil, err
	}
	trace := tracer.New(client)
	ctx = context.WithValue(ctx, traceKey, &traceData{trace, path})
	return &Tracer{
		t: trace,
		s: server,
		c: client,
	}, ctx, nil
}

type Tracer struct {
	t tracer.Trace
	s *http.Server
	c *tracer.Client
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

var emptyData = &traceData{
	path:  "",
	trace: discard{},
}

func from(ctx context.Context) *traceData {
	traceData, ok := ctx.Value(traceKey).(*traceData)
	if !ok {
		return emptyData
	}
	return traceData
}

func Start(ctx context.Context, label string, attrs ...interface{}) (context.Context, Span) {
	return from(ctx).trace.Start(ctx, label, attrs...)
}

type traceState struct {
	Path string
	Data json.RawMessage
}

func Encode(ctx context.Context) ([]byte, error) {
	data, err := tracer.Encode(ctx)
	if err != nil {
		return nil, err
	}
	td := from(ctx)
	return json.Marshal(traceState{td.path, data})
}

func Decode(ctx context.Context, data []byte) (context.Context, error) {
	// If no data to decode, just pass the context through
	if len(data) == 0 {
		return ctx, nil
	}
	var ts traceState
	if err := json.Unmarshal(data, &ts); err != nil {
		return nil, err
	}
	client, err := tracer.NewClient(ts.Path)
	if err != nil {
		return nil, err
	}
	trace := tracer.New(client)
	ctx, err = tracer.Decode(ctx, ts.Data)
	if err != nil {
		return nil, err
	}
	ctx = context.WithValue(ctx, traceKey, &traceData{trace, ts.Path})
	return ctx, nil
}
