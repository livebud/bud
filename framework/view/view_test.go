package view_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/livebud/bud/internal/cli/testcli"
	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/testdir"
	"github.com/livebud/bud/internal/versions"
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
	is.NoErr(res.DiffHeaders(`
		HTTP/1.1 200 OK
		Transfer-Encoding: chunked
		Content-Type: text/html
	`))
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
	is.NoErr(res.DiffHeaders(`
		HTTP/1.1 200 OK
		Transfer-Encoding: chunked
		Content-Type: text/html
	`))
	is.In(res.Body().String(), "<h1>hi</h1>")
	is.NoErr(app.Close())
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
	is.NoErr(res.DiffHeaders(`
		HTTP/1.1 200 OK
		Transfer-Encoding: chunked
		Content-Type: text/html
	`))
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
	is.NoErr(res.DiffHeaders(`
		HTTP/1.1 200 OK
		Transfer-Encoding: chunked
		Content-Type: text/html
	`))
	is.In(res.Body().String(), "<h1>hello</h1>")
	// Try the entrypoint
	res, err = app.Get("/bud/view/_index.svelte.js")
	is.NoErr(err)
	is.NoErr(res.DiffHeaders(`
		HTTP/1.1 200 OK
		Accept-Ranges: bytes
		Content-Type: application/javascript
	`))
	is.In(res.Body().String(), "bud_target")
	is.NoErr(app.Close())
}

var chunkRe = regexp.MustCompile(`chunk-[A-Za-z0-9]+\.js`)

func findChunk(name, src string) (string, error) {
	chunks := chunkRe.FindAllString(src, -1)
	if len(chunks) == 0 {
		return "", fmt.Errorf("unable to find a chunk in %q", name)
	}
	return chunks[0], nil
}

func TestChunks(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["controller/controller.go"] = `
		package controller
		type Controller struct {}
		func (c *Controller) Index() string { return "" }
		func (c *Controller) Show() string { return "" }
	`
	td.Files["view/index.svelte"] = `<h1>index</h1>`
	td.Files["view/show.svelte"] = `<h1>show</h1>`
	td.NodeModules["svelte"] = versions.Svelte
	td.NodeModules["livebud"] = "*"
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run", "--embed")
	is.NoErr(err)
	defer app.Close()
	// Ensure we have an index
	res, err := app.Get("/")
	is.NoErr(err)
	is.NoErr(res.DiffHeaders(`
		HTTP/1.1 200 OK
		Transfer-Encoding: chunked
		Content-Type: text/html
	`))
	is.In(res.Body().String(), "<h1>index</h1>")
	// Try the index entrypoint
	res, err = app.Get("/bud/view/_index.svelte.js")
	is.NoErr(err)
	is.NoErr(res.DiffHeaders(`
		HTTP/1.1 200 OK
		Accept-Ranges: bytes
		Content-Type: application/javascript
	`))
	is.In(res.Body().String(), "bud_target")
	// Ensure we have a show
	res, err = app.Get("/10")
	is.NoErr(err)
	is.NoErr(res.DiffHeaders(`
		HTTP/1.1 200 OK
		Transfer-Encoding: chunked
		Content-Type: text/html
	`))
	is.In(res.Body().String(), "<h1>show</h1>")
	// Try the show entrypoint
	res, err = app.Get("/bud/view/_show.svelte.js")
	is.NoErr(err)
	is.NoErr(res.DiffHeaders(`
		HTTP/1.1 200 OK
		Accept-Ranges: bytes
		Content-Type: application/javascript
	`))
	is.In(res.Body().String(), "bud_target")
	// Ensure the code's been split and find the name of the chunk
	chunkName, err := findChunk("bud/view/_show.svelte.js", res.Body().String())
	is.NoErr(err)
	res, err = app.Get("/bud/view/" + chunkName)
	is.NoErr(err)
	is.NoErr(res.DiffHeaders(`
		HTTP/1.1 200 OK
		Accept-Ranges: bytes
		Content-Type: application/javascript
	`))
	is.In(res.Body().String(), "bud_props")
	is.NoErr(app.Close())
}

func TestConsoleLog(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["controller/controller.go"] = `
		package controller
		type Controller struct {}
		func (c *Controller) Index() string { return "" }
	`
	td.Files["view/index.svelte"] = `
		<script>
			console.log("log", "!", "!")
		</script>
		<h1>hello</h1>
	`
	td.NodeModules["svelte"] = versions.Svelte
	td.NodeModules["livebud"] = "*"
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.Get("/")
	is.NoErr(err)
	is.NoErr(res.DiffHeaders(`
		HTTP/1.1 200 OK
		Transfer-Encoding: chunked
		Content-Type: text/html
	`))
	is.In(res.Body().String(), "<h1>hello</h1>")
	is.NoErr(app.Close())
	is.In(app.Stdout(), "") // TODO: console.log() should go to stdout
	is.In(app.Stderr(), "")
}

func TestConsoleError(t *testing.T) {
	// TODO: console.error needs to be added to:
	// https://github.com/kuoruan/v8go-polyfills
	t.SkipNow()
}
