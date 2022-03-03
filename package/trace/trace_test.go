package trace_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/package/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func exporter() *testExporter {
	return newTestExporter()
}

func TestTrace(t *testing.T) {
	// Setup functions
	d := func(tracer *trace.Tracer, ctx context.Context) (err error) {
		_, span := tracer.Start(ctx, "d")
		defer span.End(&err)
		return nil
	}
	b := func(tracer *trace.Tracer, ctx context.Context) (err error) {
		_, span := tracer.Start(ctx, "b")
		defer span.End(&err)
		return nil
	}
	c := func(tracer *trace.Tracer, ctx context.Context) (err error) {
		ctx, span := tracer.Start(ctx, "c")
		defer span.End(&err)
		if err := d(tracer, ctx); err != nil {
			return err
		}
		return nil
	}
	a := func(tracer *trace.Tracer, ctx context.Context) (err error) {
		ctx, span := tracer.Start(ctx, "a")
		defer span.End(&err)
		if err := b(tracer, ctx); err != nil {
			return err
		}
		if err := c(tracer, ctx); err != nil {
			return err
		}
		return nil
	}
	// Test
	is := is.New(t)
	ctx := context.Background()
	exporter := exporter()
	tracer := trace.New(exporter)
	err := a(tracer, ctx)
	is.NoErr(err)
	actual := exporter.Print()
	is.Equal(actual, `a (b c (d))`)
}

func TestPropagation(t *testing.T) {
	exporter := exporter()
	// Setup functions
	e := func(tracer *trace.Tracer, ctx context.Context) (err error) {
		_, span := tracer.Start(ctx, "e")
		defer span.End(&err)
		return nil
	}
	d := func(tracer *trace.Tracer, ctx context.Context) (err error) {
		ctx, span := tracer.Start(ctx, "d")
		defer span.End(&err)
		if err := e(tracer, ctx); err != nil {
			return err
		}
		return nil
	}
	subprocess := func(data []byte) (err error) {
		tracer := trace.New(exporter)
		ctx, err := trace.Decode(context.Background(), data)
		if err != nil {
			return err
		}
		ctx, span := tracer.Start(ctx, "subprocess")
		defer span.End(&err)
		if err := d(tracer, ctx); err != nil {
			return err
		}
		return nil
	}
	b := func(tracer *trace.Tracer, ctx context.Context) (err error) {
		_, span := tracer.Start(ctx, "b")
		defer span.End(&err)
		return nil
	}
	c := func(tracer *trace.Tracer, ctx context.Context) (err error) {
		ctx, span := tracer.Start(ctx, "c")
		defer span.End(&err)
		data, err := trace.Encode(ctx)
		if err != nil {
			return err
		}
		if err := subprocess(data); err != nil {
			return err
		}
		return nil
	}
	a := func(tracer *trace.Tracer, ctx context.Context) (err error) {
		ctx, span := tracer.Start(ctx, "a")
		defer span.End(&err)
		if err := b(tracer, ctx); err != nil {
			return err
		}
		if err := c(tracer, ctx); err != nil {
			return err
		}
		return nil
	}
	// Test
	is := is.New(t)
	ctx := context.Background()
	tracer := trace.New(exporter)
	err := a(tracer, ctx)
	is.NoErr(err)
	actual := exporter.Print()
	is.Equal(actual, `a (b c (subprocess (d (e))))`)
}

// testExporter is an otel exporter for traces
type testExporter struct {
	m map[trace.SpanID][]sdktrace.ReadOnlySpan // key is parent SpanID
}

var _ sdktrace.SpanExporter = (*testExporter)(nil)

func newTestExporter() *testExporter {
	return &testExporter{m: map[trace.SpanID][]sdktrace.ReadOnlySpan{}}
}

func (e *testExporter) ExportSpans(ctx context.Context, ss []sdktrace.ReadOnlySpan) error {
	for _, s := range ss {
		sid := s.Parent().SpanID()
		e.m[sid] = append(e.m[sid], s)
	}
	return nil
}

func (e *testExporter) Shutdown(ctx context.Context) error {
	return nil
}

func (e *testExporter) Print() string {
	root := e.m[trace.SpanID{}][0]
	var buf bytes.Buffer
	e.print(&buf, root)
	return buf.String()
}

func (e *testExporter) print(w io.Writer, ss sdktrace.ReadOnlySpan) {
	fmt.Fprintf(w, "%s", ss.Name())
	children := e.m[ss.SpanContext().SpanID()]
	if len(children) > 0 {
		fmt.Fprint(w, " (")
		for i, ss := range children {
			if i != 0 {
				fmt.Fprint(w, " ")
			}
			e.print(w, ss)
		}
		fmt.Fprint(w, ")")
	}
}
