package budclient_test

import (
	"context"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/livebud/bud/package/budfs"
	"github.com/livebud/bud/package/log/testlog"

	"github.com/livebud/bud/framework/transform/transformrt"
	"github.com/livebud/bud/framework/view/dom"
	"github.com/livebud/bud/framework/view/nodemods"
	"github.com/livebud/bud/framework/view/ssr"
	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/pubsub"
	"github.com/livebud/bud/internal/testdir"
	"github.com/livebud/bud/internal/versions"
	"github.com/livebud/bud/package/budclient"
	"github.com/livebud/bud/package/budserver"
	"github.com/livebud/bud/package/gomod"
	v8 "github.com/livebud/bud/package/js/v8"
	"github.com/livebud/bud/package/svelte"
)

func loadServer(bus pubsub.Client, dir string) (*httptest.Server, error) {
	log := testlog.New()
	vm, err := v8.Load()
	if err != nil {
		return nil, err
	}
	svelteCompiler, err := svelte.Load(vm)
	if err != nil {
		return nil, err
	}
	transforms, err := transformrt.Load(svelte.NewTransformable(svelteCompiler))
	if err != nil {
		return nil, err
	}
	module, err := gomod.Find(dir)
	if err != nil {
		return nil, err
	}
	bfs := budfs.New(module, log)
	bfs.DirGenerator("bud/view", dom.New(module, transforms.DOM))
	bfs.DirGenerator("bud/node_modules", nodemods.New(module))
	bfs.FileGenerator("bud/view/_ssr.js", ssr.New(module, transforms.SSR))
	handler := budserver.New(bfs, bus, log, vm)
	return httptest.NewServer(handler), nil
}

func TestRender(t *testing.T) {
	ctx := context.Background()
	is := is.New(t)
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["view/index.svelte"] = `
		<script>
			export let _string = "cupcake"
		</script>
		<h1>Hello, {_string}!</h1>
	`
	td.NodeModules["svelte"] = versions.Svelte
	is.NoErr(td.Write(ctx))
	ps := pubsub.New()
	server, err := loadServer(ps, dir)
	is.NoErr(err)
	defer server.Close()
	client, err := budclient.Load(server.URL)
	is.NoErr(err)
	res, err := client.Render("/", map[string]interface{}{
		"_string": "marshmallow",
	})
	is.NoErr(err)
	is.Equal(res.Status, 200)
	is.Equal(len(res.Headers), 1)
	is.Equal(res.Headers["Content-Type"], "text/html")
	is.In(res.Body, `<h1>Hello, marshmallow!</h1>`)
}

func TestRenderNested(t *testing.T) {
	ctx := context.Background()
	is := is.New(t)
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["view/posts/comments/edit.svelte"] = `
		<script>
			export let _string = "cupcake"
		</script>
		<h1>Hello, {_string}!</h1>
	`
	td.NodeModules["svelte"] = versions.Svelte
	is.NoErr(td.Write(ctx))
	ps := pubsub.New()
	server, err := loadServer(ps, dir)
	is.NoErr(err)
	defer server.Close()
	client, err := budclient.Load(server.URL)
	is.NoErr(err)
	res, err := client.Render("/posts/:post_id/comments/:id/edit", map[string]interface{}{
		"_string": "marshmallow",
	})
	is.NoErr(err)
	is.Equal(res.Status, 200)
	is.Equal(len(res.Headers), 1)
	is.Equal(res.Headers["Content-Type"], "text/html")
	is.In(res.Body, `<h1>Hello, marshmallow!</h1>`)
}

func TestProxyFile(t *testing.T) {
	ctx := context.Background()
	is := is.New(t)
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["view/index.svelte"] = `
		<script>
			export let _string = "cupcake"
		</script>
		<h1>Hello, {_string}!</h1>
	`
	td.NodeModules["svelte"] = versions.Svelte
	is.NoErr(td.Write(ctx))
	bus := pubsub.New()
	server, err := loadServer(bus, dir)
	is.NoErr(err)
	defer server.Close()
	client, err := budclient.Load(server.URL)
	is.NoErr(err)

	// Check the entrypoint
	req := httptest.NewRequest("GET", "/bud/view/_index.svelte.js", nil)
	rec := httptest.NewRecorder()
	client.Proxy(rec, req)
	res := rec.Result()
	defer res.Body.Close()
	is.Equal(res.StatusCode, 200)
	is.Equal(len(res.Header), 4)
	is.In(res.Header.Get("Content-Type"), "/javascript")
	is.Equal(res.Header.Get("Accept-Ranges"), "bytes")
	is.True(res.Header["Content-Length"] != nil)
	is.True(res.Header["Date"] != nil)
	body, err := io.ReadAll(res.Body)
	is.NoErr(err)
	is.In(string(body), `Hello, `)
	is.In(string(body), `cupcake`)

	// Check the component
	req = httptest.NewRequest("GET", "/bud/view/index.svelte", nil)
	rec = httptest.NewRecorder()
	client.Proxy(rec, req)
	res = rec.Result()
	defer res.Body.Close()
	is.Equal(res.StatusCode, 200)
	is.Equal(len(res.Header), 4)
	is.Equal(res.Header.Get("Content-Type"), "application/javascript")
	is.Equal(res.Header.Get("Accept-Ranges"), "bytes")
	is.True(res.Header["Content-Length"] != nil)
	is.True(res.Header["Date"] != nil)
	body, err = io.ReadAll(res.Body)
	is.NoErr(err)
	is.In(string(body), `Hello, `)
	is.In(string(body), `cupcake`)

	// Check the node_modules
	req = httptest.NewRequest("GET", "/bud/node_modules/svelte/internal", nil)
	rec = httptest.NewRecorder()
	client.Proxy(rec, req)
	res = rec.Result()
	defer res.Body.Close()
	is.Equal(res.StatusCode, 200)
	is.Equal(len(res.Header), 4)
	is.Equal(res.Header.Get("Content-Type"), "application/javascript")
	is.Equal(res.Header.Get("Accept-Ranges"), "bytes")
	is.True(res.Header["Content-Length"] != nil)
	is.True(res.Header["Date"] != nil)
	body, err = io.ReadAll(res.Body)
	is.NoErr(err)
	is.In(string(body), `function element(`)
	is.In(string(body), `function text(`)
}

func TestEvents(t *testing.T) {
	ctx := context.Background()
	is := is.New(t)
	dir := t.TempDir()
	td := testdir.New(dir)
	is.NoErr(td.Write(ctx))
	ps := pubsub.New()
	server, err := loadServer(ps, dir)
	is.NoErr(err)
	defer server.Close()
	client, err := budclient.Load(server.URL)
	is.NoErr(err)
	sub := ps.Subscribe("ready")
	defer sub.Close()
	err = client.Publish("ready", []byte("test"))
	is.NoErr(err)
	select {
	case payload := <-sub.Wait():
		is.Equal(string(payload), "test")
	default:
		t.Fatalf("missing event")
	}
}
