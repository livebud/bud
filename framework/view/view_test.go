package view_test

import (
	"context"
	"testing"

	"github.com/livebud/bud/internal/cli/testcli"
	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/testdir"
	"github.com/livebud/bud/internal/versions"
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
	td.NodeModules["svelte"] = versions.Svelte
	td.NodeModules["livebud"] = "*"
	is.NoErr(td.Write(ctx))

	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	hot, err := app.Hot("/bud/hot/view/index.svelte")
	is.NoErr(err)
	defer hot.Close()
	res, err := app.Get("/")
	is.NoErr(err)
	diff.TestHTTP(t, res.Headers().String(), `
		HTTP/1.1 200 OK
		Transfer-Encoding: chunked
		Content-Type: text/html
	`)
	is.In(res.Body().String(), "<h1>hello</h1>")
	is.NoErr(td.Exists("bud/internal/app/view/view.go"))
	// Change svelte file
	td = testdir.New(dir)
	td.Files["view/index.svelte"] = `<h1>hi</h1>`
	is.NoErr(td.Write(ctx))
	// Wait for the app to be ready again
	app.Ready(ctx)
	// Check that we received a hot reload event
	event, err := hot.Next(ctx)
	is.NoErr(err)
	is.In(string(event.Data), `{"scripts":["/bud/view/index.svelte?ts=`)
	// Should change
	res, err = app.Get("/")
	is.NoErr(err)
	diff.TestHTTP(t, res.Headers().String(), `
			HTTP/1.1 200 OK
			Transfer-Encoding: chunked
			Content-Type: text/html
		`)
	is.In(res.Body().String(), "<h1>hi</h1>")
	is.NoErr(app.Close())
}

// Note: if this test is failing due to context deadline exceeding, you
// probably just need update the timeout. Right now we don't have a signal
// that Start() has built and started the app.
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
	td.NodeModules["svelte"] = versions.Svelte
	td.NodeModules["livebud"] = "*"
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run", "--embed")
	is.NoErr(err)
	defer app.Close()
	hot, err := app.Hot("/bud/hot/view/index.svelte")
	is.NoErr(err)
	defer hot.Close()
	res, err := app.Get("/")
	is.NoErr(err)
	diff.TestHTTP(t, res.Headers().String(), `
		HTTP/1.1 200 OK
		Transfer-Encoding: chunked
		Content-Type: text/html
	`)
	is.In(res.Body().String(), "<h1>hello</h1>")
	// Change svelte file
	td = testdir.New(dir)
	td.Files["view/index.svelte"] = `<h1>hi</h1>`
	is.NoErr(td.Write(ctx))
	// Wait for the the app to be ready again
	is.NoErr(app.Ready(ctx))
	// Ensure that we got a hot reload event
	event, err := hot.Next(ctx)
	is.NoErr(err)
	is.In(string(event.Data), `{"scripts":["/bud/view/index.svelte?ts=`)
	// Shouldn't be any change
	res, err = app.Get("/")
	is.NoErr(err)
	diff.TestHTTP(t, res.Headers().String(), `
		HTTP/1.1 200 OK
		Transfer-Encoding: chunked
		Content-Type: text/html
	`)
	is.In(res.Body().String(), "<h1>hello</h1>")
	is.NoErr(app.Close())
}
