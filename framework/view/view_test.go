package view_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"
	"time"

	"github.com/lithammer/dedent"
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
	is.NoErr(td.Exists("bud/internal/web/view/view.go"))
	// Change svelte file
	indexFile := filepath.Join(dir, "view/index.svelte")
	is.NoErr(os.MkdirAll(filepath.Dir(indexFile), 0755))
	is.NoErr(os.WriteFile(indexFile, []byte(`<h1>hi</h1>`), 0644))
	// Wait for the app to be ready again
	readyCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	is.NoErr(app.Ready(readyCtx))
	cancel()
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
	// Change svelte file
	indexFile = filepath.Join(dir, "view/index.svelte")
	is.NoErr(os.MkdirAll(filepath.Dir(indexFile), 0755))
	is.NoErr(os.WriteFile(indexFile, []byte(`<h1>hola</h1>`), 0644))
	// Wait for the app to be ready again
	readyCtx, cancel = context.WithTimeout(ctx, 15*time.Second)
	is.NoErr(app.Ready(readyCtx))
	cancel()
	// Check that we received a hot reload event
	event, err = hot.Next(ctx)
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
	is.In(res.Body().String(), "<h1>hola</h1>")
	// Change svelte file one more time
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
	indexFile := filepath.Join(dir, "view/index.svelte")
	is.NoErr(os.MkdirAll(filepath.Dir(indexFile), 0755))
	is.NoErr(os.WriteFile(indexFile, []byte(`<h1>hi</h1>`), 0644))
	// Wait for the app to be ready again
	readyCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	is.NoErr(app.Ready(readyCtx))
	cancel()
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

func TestRenameView(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["controller/controller.go"] = `
		package controller
		type Controller struct {}
		func (c *Controller) Show() (id int) { return 10 }
	`
	td.Files["view/show.svelte"] = `
		<script>
			export let id = 0
		</script>
		<h1>{id}</h1>
	`
	td.NodeModules["svelte"] = versions.Svelte
	td.NodeModules["livebud"] = "*"
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	hot, err := app.Hot("/bud/hot/view/show.svelte")
	is.NoErr(err)
	defer hot.Close()
	res, err := app.Get("/10")
	is.NoErr(err)
	is.NoErr(res.DiffHeaders(`
		HTTP/1.1 200 OK
		Transfer-Encoding: chunked
		Content-Type: text/html
	`))
	is.In(res.Body().String(), "<h1>10</h1>")
	// Rename the file
	is.NoErr(os.Rename(
		filepath.Join(dir, "view/show.svelte"),
		filepath.Join(dir, "view/_show.svele"),
	))
	// Wait for the app to be ready again
	readyCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	is.NoErr(app.Ready(readyCtx))
	cancel()
	// Check that we received a hot reload event
	event, err := hot.Next(ctx)
	is.NoErr(err)
	is.In(string(event.Data), `{"reload":true}`)
	// Should change
	res, err = app.Get("/10")
	is.NoErr(err)
	is.NoErr(res.DiffHeaders(`
		HTTP/1.1 200 OK
		Content-Type: application/json
	`))
	is.Equal(res.Body().String(), "10")
	is.NoErr(app.Close())
}

func TestAddView(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["controller/controller.go"] = `
		package controller
		type Controller struct {}
		func (c *Controller) Show() (id int) { return 10 }
	`
	td.NodeModules["svelte"] = versions.Svelte
	td.NodeModules["livebud"] = "*"
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	hot, err := app.Hot("/bud/hot/view/show.svelte")
	is.NoErr(err)
	defer hot.Close()
	res, err := app.Get("/10")
	is.NoErr(err)
	is.NoErr(res.DiffHeaders(`
		HTTP/1.1 200 OK
		Content-Type: application/json
	`))
	is.Equal(res.Body().String(), "10")
	// Add the view
	showView := filepath.Join(dir, "view/show.svelte")
	is.NoErr(os.MkdirAll(filepath.Dir(showView), 0755))
	is.NoErr(os.WriteFile(showView, []byte(dedent.Dedent(`
		<script>
			export let id = 0
		</script>
		<h1>{id}</h1>
	`)), 0644))
	// Wait for the app to be ready again
	readyCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	is.NoErr(app.Ready(readyCtx))
	cancel()
	// Check that we received a hot reload event
	event, err := hot.Next(ctx)
	is.NoErr(err)
	is.In(string(event.Data), `{"reload":true}`)
	// Should change
	res, err = app.Get("/10")
	is.NoErr(err)
	is.NoErr(res.DiffHeaders(`
		HTTP/1.1 200 OK
		Transfer-Encoding: chunked
		Content-Type: text/html
	`))
	is.In(res.Body().String(), "<h1>10</h1>")
}

func TestSvelteImportFromNodeModule(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.NodeModules["svelte"] = versions.Svelte
	td.NodeModules["svelte-time"] = "*"
	td.Files["controller/controller.go"] = `
		package controller
		type Controller struct {}
		func (c *Controller) Index() string { return "" }
	`
	td.Files["view/index.svelte"] = `
		<script>
			import Time from 'svelte-time';
		</script>
		<p>The time is <Time timestamp="2022-07-19 10:19:00" />.</p>
	`
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
	is.In(res.Body().String(), "<time datetime=\"2022-07-19 10:19:00\">Jul 19, 2022</time>")
	is.NoErr(app.Close())
}

func TestSvelteImportFromOtherDir(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.NodeModules["svelte"] = versions.Svelte
	td.NodeModules["svelte-time"] = "*"
	td.NodeModules["livebud"] = "*"
	td.Files["controller/controller.go"] = `
		package controller
		type Controller struct {}
		func (c *Controller) Index() string { return "" }
	`
	td.Files["ui/Time.svelte"] = `
		<h1>The Time</h1>
	`
	td.Files["view/index.svelte"] = `
	 	<script>
			import Time from "../ui/Time.svelte"
		</script>
		<Time/>
	`
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
	is.In(res.Body().String(), "<h1>The Time</h1>")
	is.NoErr(app.Close())
}
