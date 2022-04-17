package trace_test

import (
	"context"
	"strings"
	"testing"

	"github.com/livebud/bud/package/trace"

	"github.com/matryer/is"
)

func TestStart(t *testing.T) {
	// Setup functions
	d := func(ctx context.Context) (err error) {
		_, span := trace.Start(ctx, "d")
		defer span.End(&err)
		return nil
	}
	b := func(ctx context.Context) (err error) {
		_, span := trace.Start(ctx, "b")
		defer span.End(&err)
		return nil
	}
	c := func(ctx context.Context) (err error) {
		ctx, span := trace.Start(ctx, "c")
		defer span.End(&err)
		if err := d(ctx); err != nil {
			return err
		}
		return nil
	}
	a := func(ctx context.Context) (err error) {
		ctx, span := trace.Start(ctx, "a")
		defer span.End(&err)
		if err := b(ctx); err != nil {
			return err
		}
		if err := c(ctx); err != nil {
			return err
		}
		return nil
	}
	is := is.New(t)
	ctx := context.Background()
	tracer, ctx, err := trace.Serve(ctx)
	is.NoErr(err)
	err = a(ctx)
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

func TestCodec(t *testing.T) {
	// Setup functions
	e := func(ctx context.Context) (err error) {
		_, span := trace.Start(ctx, "e")
		defer span.End(&err)
		return nil
	}
	d := func(ctx context.Context) (err error) {
		ctx, span := trace.Start(ctx, "d")
		defer span.End(&err)
		if err := e(ctx); err != nil {
			return err
		}
		return nil
	}
	subprocess := func(data []byte) (err error) {
		ctx := context.Background()
		ctx, err = trace.Decode(ctx, data)
		if err != nil {
			return err
		}
		ctx, span := trace.Start(ctx, "subprocess")
		defer span.End(&err)
		if err := d(ctx); err != nil {
			return err
		}
		return nil
	}
	b := func(ctx context.Context) (err error) {
		_, span := trace.Start(ctx, "b")
		defer span.End(&err)
		return nil
	}
	c := func(ctx context.Context) (err error) {
		ctx, span := trace.Start(ctx, "c")
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
	a := func(ctx context.Context) (err error) {
		ctx, span := trace.Start(ctx, "a")
		defer span.End(&err)
		if err := b(ctx); err != nil {
			return err
		}
		if err := c(ctx); err != nil {
			return err
		}
		return nil
	}
	// Start Test
	is := is.New(t)
	ctx := context.Background()
	tracer, ctx, err := trace.Serve(ctx)
	is.NoErr(err)
	err = a(ctx)
	is.NoErr(err)
	tree, err := tracer.Print(ctx)
	is.NoErr(err)
	is.True(strings.Contains(tree, "a ("))
	is.True(strings.Contains(tree, "├── b ("))
	is.True(strings.Contains(tree, "└── c ("))
	is.True(strings.Contains(tree, "    └── subprocess ("))
	is.True(strings.Contains(tree, "        └── d ("))
	is.True(strings.Contains(tree, "            └── e ("))
	err = tracer.Shutdown(ctx)
	is.NoErr(err)
}
