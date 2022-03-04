package trace_test

import (
	"context"
	"fmt"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/package/trace"
)

func TestSingleClient(t *testing.T) {
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
	server := httptest.NewServer(trace.Handler())
	defer server.Close()
	client, err := trace.NewClient(server.URL)
	is.NoErr(err)
	tracer := trace.New(client)
	err = a(tracer, ctx)
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
	server := httptest.NewServer(trace.Handler())
	defer server.Close()
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
		client, err := trace.NewClient(server.URL)
		is.NoErr(err)
		tracer := trace.New(client)
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
	// Start Test
	client, err := trace.NewClient(server.URL)
	is.NoErr(err)
	tracer := trace.New(client)
	err = a(tracer, ctx)
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
	is := is.New(t)
	socketPath := filepath.Join(t.TempDir(), "trace.sock")
	server, err := trace.Serve(socketPath)
	is.NoErr(err)
	ctx := context.Background()
	client, err := trace.NewClient(socketPath)
	is.NoErr(err)
	tracer := trace.New(client)
	err = a(tracer, ctx)
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

func TestServerError(t *testing.T) {
	// Setup functions
	d := func(tracer *trace.Tracer, ctx context.Context) (err error) {
		_, span := tracer.Start(ctx, "d")
		defer span.End(&err)
		return fmt.Errorf("oh noz")
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
	server := httptest.NewServer(trace.Handler())
	defer server.Close()
	client, err := trace.NewClient(server.URL)
	is.NoErr(err)
	tracer := trace.New(client)
	err = a(tracer, ctx)
	is.Equal(err.Error(), "oh noz")
	tree, err := client.Print(ctx)
	is.NoErr(err)
	is.True(strings.Contains(tree, `) error="oh noz"`))
	is.True(strings.Contains(tree, "a ("))
	is.True(strings.Contains(tree, "├── b ("))
	is.True(strings.Contains(tree, "└── c ("))
	is.True(strings.Contains(tree, "    └── d ("))
}

// func TestAttributes(t *testing.T) {
// 	// Setup functions
// 	d := func(tracer *trace.Tracer, ctx context.Context) (err error) {
// 		_, span := tracer.Start(ctx, "d", "path", "/")
// 		defer span.End(&err)
// 		return nil
// 	}
// 	b := func(tracer *trace.Tracer, ctx context.Context) (err error) {
// 		_, span := tracer.Start(ctx, "b")
// 		defer span.End(&err)
// 		return nil
// 	}
// 	c := func(tracer *trace.Tracer, ctx context.Context) (err error) {
// 		ctx, span := tracer.Start(ctx, "c")
// 		defer span.End(&err)
// 		if err := d(tracer, ctx); err != nil {
// 			return err
// 		}
// 		return nil
// 	}
// 	a := func(tracer *trace.Tracer, ctx context.Context) (err error) {
// 		ctx, span := tracer.Start(ctx, "a", "port", 3000, "id", "10")
// 		defer span.End(&err)
// 		if err := b(tracer, ctx); err != nil {
// 			return err
// 		}
// 		if err := c(tracer, ctx); err != nil {
// 			return err
// 		}
// 		return nil
// 	}
// 	// Test
// 	is := is.New(t)
// 	ctx := context.Background()
// 	exporter := exporter()
// 	tracer := trace.New(exporter)
// 	err := a(tracer, ctx)
// 	is.NoErr(err)
// 	actual := exporter.Print()
// 	is.Equal(actual, `a {id:10} {port:3000} (b c (d {path:/}))`)
// }
