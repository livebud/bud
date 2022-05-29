package view_test

import (
	"context"
	"testing"
	"time"

	"github.com/livebud/bud/internal/cli"
	"github.com/livebud/bud/internal/cli/testcli"
	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/testdir"
	"github.com/livebud/bud/internal/version"
	"github.com/matthewmueller/diff"
)

func TestHello(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["controller/controller.go"] = `
		package controller
		type Controller struct {}
		func (c *Controller) Index() string { return "" }
	`
	td.Files["view/index.svelte"] = `<h1>hello</h1>`
	td.NodeModules["svelte"] = version.Svelte
	td.NodeModules["livebud"] = "*"
	is.NoErr(td.Write(ctx))
	cli := testcli.New(cli.New(dir))
	app, stdout, stderr, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.Get("/")
	is.NoErr(err)
	diff.TestHTTP(t, res.Headers().String(), `
		HTTP/1.1 200 OK
		Content-Type: text/html
	`)
	is.In(res.Body().String(), "<h1>hello</h1>")
	is.NoErr(td.Exists("bud/.app/view/view.go"))
	// Change svelte file
	td.Files["view/index.svelte"] = `<h1>hi</h1>`
	is.NoErr(td.Write(ctx))
	// Should change
	res, err = app.Get("/")
	is.NoErr(err)
	diff.TestHTTP(t, res.Headers().String(), `
			HTTP/1.1 200 OK
			Content-Type: text/html
		`)
	is.In(res.Body().String(), "<h1>hi</h1>")
	is.Equal(stdout.String(), "")
	is.In(stderr.String(), "info: Ready on")
}

func TestHelloEmbed(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["controller/controller.go"] = `
		package controller
		type Controller struct {}
		func (c *Controller) Index() string { return "" }
	`
	td.Files["view/index.svelte"] = `<h1>hello</h1>`
	td.NodeModules["svelte"] = version.Svelte
	td.NodeModules["livebud"] = "*"
	is.NoErr(td.Write(ctx))
	cli := testcli.New(cli.New(dir))
	app, stdout, stderr, err := cli.Start(ctx, "run", "--embed")
	is.NoErr(err)
	defer app.Close()
	hot, err := app.Hot("/bud/view/index.svelte")
	is.NoErr(err)
	defer hot.Close()
	res, err := app.Get("/")
	is.NoErr(err)
	diff.TestHTTP(t, res.Headers().String(), `
		HTTP/1.1 200 OK
		Content-Type: text/html
	`)
	is.In(res.Body().String(), "<h1>hello</h1>")
	// Change svelte file
	td = testdir.New(dir)
	td.Files["view/index.svelte"] = `<h1>hi</h1>`
	is.NoErr(td.Write(ctx))
	// Wait for the change event
	eventCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	event, err := hot.Next(eventCtx)
	is.NoErr(err)
	is.In(string(event.Data), `{"scripts":["?ts=`)
	// Shouldn't be any change
	res, err = app.Get("/")
	is.NoErr(err)
	diff.TestHTTP(t, res.Headers().String(), `
		HTTP/1.1 200 OK
		Content-Type: text/html
	`)
	is.In(res.Body().String(), "<h1>hello</h1>")
	is.Equal(stdout.String(), "")
	is.In(stderr.String(), "")
}
