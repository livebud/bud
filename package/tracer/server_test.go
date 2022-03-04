package tracer_test

import (
	"context"
	"fmt"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/package/tracer"
)

func TestSingleClient(t *testing.T) {
	// Setup functions
	d := func(trace tracer.Trace, ctx context.Context) (err error) {
		_, span := trace.Start(ctx, "d")
		defer span.End(&err)
		return nil
	}
	b := func(trace tracer.Trace, ctx context.Context) (err error) {
		_, span := trace.Start(ctx, "b")
		defer span.End(&err)
		return nil
	}
	c := func(trace tracer.Trace, ctx context.Context) (err error) {
		ctx, span := trace.Start(ctx, "c")
		defer span.End(&err)
		if err := d(trace, ctx); err != nil {
			return err
		}
		return nil
	}
	a := func(trace tracer.Trace, ctx context.Context) (err error) {
		ctx, span := trace.Start(ctx, "a")
		defer span.End(&err)
		if err := b(trace, ctx); err != nil {
			return err
		}
		if err := c(trace, ctx); err != nil {
			return err
		}
		return nil
	}
	// Test
	is := is.New(t)
	ctx := context.Background()
	server := httptest.NewServer(tracer.Handler())
	defer server.Close()
	client, err := tracer.NewClient(server.URL)
	is.NoErr(err)
	trace := tracer.New(client)
	err = a(trace, ctx)
	is.NoErr(err)
	tree, err := client.Print(ctx)
	is.NoErr(err)
	is.True(strings.Contains(tree, "a ("))
	is.True(strings.Contains(tree, "├── b ("))
	is.True(strings.Contains(tree, "└── c ("))
	is.True(strings.Contains(tree, "    └── d ("))
}

func TestMultiClient(t *testing.T) {
	// Setup server
	is := is.New(t)
	ctx := context.Background()
	server := httptest.NewServer(tracer.Handler())
	defer server.Close()
	// Setup functions
	e := func(trace tracer.Trace, ctx context.Context) (err error) {
		_, span := trace.Start(ctx, "e")
		defer span.End(&err)
		return nil
	}
	d := func(trace tracer.Trace, ctx context.Context) (err error) {
		ctx, span := trace.Start(ctx, "d")
		defer span.End(&err)
		if err := e(trace, ctx); err != nil {
			return err
		}
		return nil
	}
	subprocess := func(data []byte) (err error) {
		ctx, trace, err := tracer.Resume(context.Background(), server.URL, data)
		if err != nil {
			return err
		}
		ctx, span := trace.Start(ctx, "subprocess")
		defer span.End(&err)
		if err := d(trace, ctx); err != nil {
			return err
		}
		return nil
	}
	b := func(trace tracer.Trace, ctx context.Context) (err error) {
		_, span := trace.Start(ctx, "b")
		defer span.End(&err)
		return nil
	}
	c := func(t tracer.Trace, ctx context.Context) (err error) {
		ctx, span := t.Start(ctx, "c")
		defer span.End(&err)
		data, err := tracer.Encode(ctx)
		if err != nil {
			return err
		}
		if err := subprocess(data); err != nil {
			return err
		}
		return nil
	}
	a := func(trace tracer.Trace, ctx context.Context) (err error) {
		ctx, span := trace.Start(ctx, "a")
		defer span.End(&err)
		if err := b(trace, ctx); err != nil {
			return err
		}
		if err := c(trace, ctx); err != nil {
			return err
		}
		return nil
	}
	// Start Test
	client, err := tracer.NewClient(server.URL)
	is.NoErr(err)
	trace := tracer.New(client)
	err = a(trace, ctx)
	is.NoErr(err)
	// Print results
	tree, err := client.Print(ctx)
	is.NoErr(err)
	is.True(strings.Contains(tree, "a ("))
	is.True(strings.Contains(tree, "├── b ("))
	is.True(strings.Contains(tree, "└── c ("))
	is.True(strings.Contains(tree, "    └── subprocess ("))
	is.True(strings.Contains(tree, "        └── d ("))
	is.True(strings.Contains(tree, "            └── e ("))
}

func TestServer(t *testing.T) {
	// Setup functions
	d := func(trace tracer.Trace, ctx context.Context) (err error) {
		_, span := trace.Start(ctx, "d")
		defer span.End(&err)
		return nil
	}
	b := func(trace tracer.Trace, ctx context.Context) (err error) {
		_, span := trace.Start(ctx, "b")
		defer span.End(&err)
		return nil
	}
	c := func(trace tracer.Trace, ctx context.Context) (err error) {
		ctx, span := trace.Start(ctx, "c")
		defer span.End(&err)
		if err := d(trace, ctx); err != nil {
			return err
		}
		return nil
	}
	a := func(trace tracer.Trace, ctx context.Context) (err error) {
		ctx, span := trace.Start(ctx, "a")
		defer span.End(&err)
		if err := b(trace, ctx); err != nil {
			return err
		}
		if err := c(trace, ctx); err != nil {
			return err
		}
		return nil
	}
	is := is.New(t)
	socketPath := filepath.Join(t.TempDir(), "trace.sock")
	server, err := tracer.Serve(socketPath)
	is.NoErr(err)
	ctx := context.Background()
	client, err := tracer.NewClient(socketPath)
	is.NoErr(err)
	trace := tracer.New(client)
	err = a(trace, ctx)
	is.NoErr(err)
	tree, err := client.Print(ctx)
	is.NoErr(err)
	is.True(strings.Contains(tree, "a ("))
	is.True(strings.Contains(tree, "├── b ("))
	is.True(strings.Contains(tree, "└── c ("))
	is.True(strings.Contains(tree, "    └── d ("))
	err = server.Shutdown(ctx)
	is.NoErr(err)
}

func TestStart(t *testing.T) {
	// Setup functions
	d := func(trace tracer.Trace, ctx context.Context) (err error) {
		_, span := trace.Start(ctx, "d")
		defer span.End(&err)
		return nil
	}
	b := func(trace tracer.Trace, ctx context.Context) (err error) {
		_, span := trace.Start(ctx, "b")
		defer span.End(&err)
		return nil
	}
	c := func(trace tracer.Trace, ctx context.Context) (err error) {
		ctx, span := trace.Start(ctx, "c")
		defer span.End(&err)
		if err := d(trace, ctx); err != nil {
			return err
		}
		return nil
	}
	a := func(trace tracer.Trace, ctx context.Context) (err error) {
		ctx, span := trace.Start(ctx, "a")
		defer span.End(&err)
		if err := b(trace, ctx); err != nil {
			return err
		}
		if err := c(trace, ctx); err != nil {
			return err
		}
		return nil
	}
	is := is.New(t)
	ctx := context.Background()
	tracer, err := tracer.Start()
	is.NoErr(err)
	err = a(tracer, ctx)
	is.NoErr(err)
	tree, err := tracer.Print(ctx)
	is.NoErr(err)
	is.True(strings.Contains(tree, "a ("))
	is.True(strings.Contains(tree, "├── b ("))
	is.True(strings.Contains(tree, "└── c ("))
	is.True(strings.Contains(tree, "    └── d ("))
	err = tracer.Shutdown(ctx)
	is.NoErr(err)
}

func TestServerError(t *testing.T) {
	// Setup functions
	d := func(trace tracer.Trace, ctx context.Context) (err error) {
		_, span := trace.Start(ctx, "d")
		defer span.End(&err)
		return fmt.Errorf("oh noz")
	}
	b := func(trace tracer.Trace, ctx context.Context) (err error) {
		_, span := trace.Start(ctx, "b")
		defer span.End(&err)
		return nil
	}
	c := func(trace tracer.Trace, ctx context.Context) (err error) {
		ctx, span := trace.Start(ctx, "c")
		defer span.End(&err)
		if err := d(trace, ctx); err != nil {
			return err
		}
		return nil
	}
	a := func(trace tracer.Trace, ctx context.Context) (err error) {
		ctx, span := trace.Start(ctx, "a")
		defer span.End(&err)
		if err := b(trace, ctx); err != nil {
			return err
		}
		if err := c(trace, ctx); err != nil {
			return err
		}
		return nil
	}
	// Test
	is := is.New(t)
	ctx := context.Background()
	server := httptest.NewServer(tracer.Handler())
	defer server.Close()
	client, err := tracer.NewClient(server.URL)
	is.NoErr(err)
	trace := tracer.New(client)
	err = a(trace, ctx)
	is.True(err != nil && err.Error() == "oh noz")
	tree, err := client.Print(ctx)
	is.NoErr(err)
	is.True(strings.Contains(tree, `) error="oh noz"`))
	is.True(strings.Contains(tree, "a ("))
	is.True(strings.Contains(tree, "├── b ("))
	is.True(strings.Contains(tree, "└── c ("))
	is.True(strings.Contains(tree, "    └── d ("))
}

func TestServerAttributes(t *testing.T) {
	// Setup functions
	d := func(trace tracer.Trace, ctx context.Context) (err error) {
		_, span := trace.Start(ctx, "d", "path", "/")
		defer span.End(&err)
		return nil
	}
	b := func(trace tracer.Trace, ctx context.Context) (err error) {
		_, span := trace.Start(ctx, "b")
		defer span.End(&err)
		return nil
	}
	c := func(trace tracer.Trace, ctx context.Context) (err error) {
		ctx, span := trace.Start(ctx, "c")
		defer span.End(&err)
		if err := d(trace, ctx); err != nil {
			return err
		}
		return nil
	}
	a := func(trace tracer.Trace, ctx context.Context) (err error) {
		ctx, span := trace.Start(ctx, "a", "port", 3000, "id", "10")
		defer span.End(&err)
		if err := b(trace, ctx); err != nil {
			return err
		}
		if err := c(trace, ctx); err != nil {
			return err
		}
		return nil
	}
	// Test
	is := is.New(t)
	ctx := context.Background()
	server := httptest.NewServer(tracer.Handler())
	defer server.Close()
	client, err := tracer.NewClient(server.URL)
	is.NoErr(err)
	trace := tracer.New(client)
	err = a(trace, ctx)
	is.NoErr(err)
	tree, err := client.Print(ctx)
	is.NoErr(err)
	is.True(strings.Contains(tree, "a ("))
	is.True(strings.Contains(tree, ") id=10 port=3000"))
	is.True(strings.Contains(tree, "├── b ("))
	is.True(strings.Contains(tree, "└── c ("))
	is.True(strings.Contains(tree, "    └── d ("))
	is.True(strings.Contains(tree, ") path=/"))
}
