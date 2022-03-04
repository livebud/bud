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

// testExporter is an otel exporter for traces
type testExporter struct {
	m map[string][]*trace.SpanData // key is parent SpanID
}

var _ sdktrace.SpanExporter = (*testExporter)(nil)

func newTestExporter() *testExporter {
	return &testExporter{m: map[string][]*trace.SpanData{}}
}

func (e *testExporter) ExportSpans(ctx context.Context, ss []sdktrace.ReadOnlySpan) error {
	for _, s := range ss {
		span := trace.ToSpanData(s)
		sid := span.ParentID
		e.m[sid] = append(e.m[sid], span)
	}
	return nil
}

func (e *testExporter) Shutdown(ctx context.Context) error {
	return nil
}

func (e *testExporter) Print() string {
	root := e.m["0000000000000000"][0]
	var buf bytes.Buffer
	e.print(&buf, root)
	return buf.String()
}

func (e *testExporter) print(w io.Writer, span *trace.SpanData) {
	fmt.Fprintf(w, "%s", span.Name)
	if span.Error != "" {
		fmt.Fprintf(w, " error=%q", span.Error)
	}
	for _, field := range span.Fields() {
		fmt.Fprintf(w, " %s=%s", field.Key, field.Value)
	}
	children := e.m[span.ID]
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

func exporter() *testExporter {
	return newTestExporter()
}

func TestTrace(t *testing.T) {
	// Setup functions
	d := func(tracer trace.Tracer, ctx context.Context) (err error) {
		_, span := tracer.Start(ctx, "d")
		defer span.End(&err)
		return nil
	}
	b := func(tracer trace.Tracer, ctx context.Context) (err error) {
		_, span := tracer.Start(ctx, "b")
		defer span.End(&err)
		return nil
	}
	c := func(tracer trace.Tracer, ctx context.Context) (err error) {
		ctx, span := tracer.Start(ctx, "c")
		defer span.End(&err)
		if err := d(tracer, ctx); err != nil {
			return err
		}
		return nil
	}
	a := func(tracer trace.Tracer, ctx context.Context) (err error) {
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
	e := func(tracer trace.Tracer, ctx context.Context) (err error) {
		_, span := tracer.Start(ctx, "e")
		defer span.End(&err)
		return nil
	}
	d := func(tracer trace.Tracer, ctx context.Context) (err error) {
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
	b := func(tracer trace.Tracer, ctx context.Context) (err error) {
		_, span := tracer.Start(ctx, "b")
		defer span.End(&err)
		return nil
	}
	c := func(tracer trace.Tracer, ctx context.Context) (err error) {
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
	a := func(tracer trace.Tracer, ctx context.Context) (err error) {
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

func TestError(t *testing.T) {
	// Setup functions
	d := func(tracer trace.Tracer, ctx context.Context) (err error) {
		_, span := tracer.Start(ctx, "d")
		defer span.End(&err)
		return fmt.Errorf("oh noz")
	}
	b := func(tracer trace.Tracer, ctx context.Context) (err error) {
		_, span := tracer.Start(ctx, "b")
		defer span.End(&err)
		return nil
	}
	c := func(tracer trace.Tracer, ctx context.Context) (err error) {
		ctx, span := tracer.Start(ctx, "c")
		defer span.End(&err)
		if err := d(tracer, ctx); err != nil {
			return err
		}
		return nil
	}
	a := func(tracer trace.Tracer, ctx context.Context) (err error) {
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
	is.Equal(err.Error(), "oh noz")
	actual := exporter.Print()
	is.Equal(actual, `a error="oh noz" (b c error="oh noz" (d error="oh noz"))`)
}

func TestAttributes(t *testing.T) {
	// Setup functions
	d := func(tracer trace.Tracer, ctx context.Context) (err error) {
		_, span := tracer.Start(ctx, "d", "path", "/")
		defer span.End(&err)
		return nil
	}
	b := func(tracer trace.Tracer, ctx context.Context) (err error) {
		_, span := tracer.Start(ctx, "b")
		defer span.End(&err)
		return nil
	}
	c := func(tracer trace.Tracer, ctx context.Context) (err error) {
		ctx, span := tracer.Start(ctx, "c")
		defer span.End(&err)
		if err := d(tracer, ctx); err != nil {
			return err
		}
		return nil
	}
	a := func(tracer trace.Tracer, ctx context.Context) (err error) {
		ctx, span := tracer.Start(ctx, "a", "port", 3000, "id", "10")
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
	is.Equal(actual, `a id=10 port=3000 (b c (d path=/))`)
}
