package app_test

import (
	"context"
	"io"
	"testing"

	"github.com/livebud/bud/internal/cli/testcli"
	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/testdir"
)

func TestWelcome(t *testing.T) {
	is := is.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	dir := t.TempDir()
	td := testdir.New(dir)
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	is.NoErr(td.NotExists("bud/app"))
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	res, err := app.Get("/")
	is.NoErr(err)
	is.Equal(res.Status(), 200)
	is.In(res.Body().String(), "Hey Bud")
	is.In(res.Body().String(), "Hey Bud") // should work multiple times
	is.NoErr(app.Close())
	stdout, err := io.ReadAll(app.Stdout())
	is.NoErr(err)
	is.In(string(stdout), "")
	stderr, err := io.ReadAll(app.Stderr())
	is.NoErr(err)
	is.In(string(stderr), "")
	is.NoErr(td.Exists("bud/app"))
}
