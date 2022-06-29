package newcontroller_test

import (
	"context"
	"testing"

	"github.com/livebud/bud/internal/cli/testcli"
	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/testdir"
)

func TestNewControllerNoActions(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	err := td.Write(ctx)
	is.NoErr(err)
	cli := testcli.New(dir)
	result, err := cli.Run(ctx, "new", "controller", "hello")
	is.NoErr(err)
	is.Equal(result.Stdout(), "")
	is.Equal(result.Stderr(), "")
	is.NoErr(td.Exists("controller/hello/controller.go"))
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.GetJSON("/")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 404 Not Found
		Content-Type: text/plain; charset=utf-8
		X-Content-Type-Options: nosniff

		404 page not found
	`))
	is.NoErr(app.Close())
}

func TestNewControllerNoActionsRoute(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	err := td.Write(ctx)
	is.NoErr(err)
	cli := testcli.New(dir)
	result, err := cli.Run(ctx, "new", "controller", "hello:/")
	is.NoErr(err)
	is.Equal(result.Stdout(), "")
	is.Equal(result.Stderr(), "")
	is.NoErr(td.Exists("controller/controller.go"))
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.GetJSON("/")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 404 Not Found
		Content-Type: text/plain; charset=utf-8
		X-Content-Type-Options: nosniff

		404 page not found
	`))
	is.NoErr(app.Close())
}

func TestNewControllerIndexShow(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	err := td.Write(ctx)
	is.NoErr(err)
	cli := testcli.New(dir)
	result, err := cli.Run(ctx, "new", "controller", "posts:/", "index", "show")
	is.NoErr(err)
	is.Equal(result.Stdout(), "")
	is.Equal(result.Stderr(), "")
	is.NoErr(td.Exists("controller/controller.go"))
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.GetJSON("/")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		[]
	`))
	res, err = app.GetJSON("/10")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		{}
	`))
	is.NoErr(app.Close())
}
