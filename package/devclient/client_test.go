package devclient_test

import (
	"context"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/pubsub"
	"github.com/livebud/bud/internal/testdir"
	"github.com/livebud/bud/internal/versions"
	"github.com/livebud/bud/package/devclient"
	"github.com/livebud/bud/package/devserver"
	"github.com/livebud/bud/package/gomod"
	v8 "github.com/livebud/bud/package/js/v8"
	"github.com/livebud/bud/package/overlay"
	"github.com/livebud/bud/package/svelte"
	"github.com/livebud/bud/runtime/transform"
	"github.com/livebud/bud/runtime/view/dom"
	"github.com/livebud/bud/runtime/view/ssr"
)

func loadServer(dir string) (*httptest.Server, error) {
	vm, err := v8.Load()
	if err != nil {
		return nil, err
	}
	svelteCompiler, err := svelte.Load(vm)
	if err != nil {
		return nil, err
	}
	transforms, err := transform.Load(svelte.NewTransformable(svelteCompiler))
	if err != nil {
		return nil, err
	}
	module, err := gomod.Find(dir)
	if err != nil {
		return nil, err
	}
	genfs, err := overlay.Load(module)
	if err != nil {
		return nil, err
	}
	genfs.FileServer("bud/view", dom.New(module, transforms.DOM))
	genfs.FileServer("bud/node_modules", dom.NodeModules(module))
	genfs.FileGenerator("bud/view/_ssr.js", ssr.New(module, transforms.SSR))
	ps := pubsub.New()
	handler := devserver.New(genfs, ps, vm)
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
	server, err := loadServer(dir)
	is.NoErr(err)
	defer server.Close()
	client, err := devclient.Load(server.URL)
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
	server, err := loadServer(dir)
	is.NoErr(err)
	defer server.Close()
	client, err := devclient.Load(server.URL)
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
	server, err := loadServer(dir)
	is.NoErr(err)
	defer server.Close()
	client, err := devclient.Load(server.URL)
	is.NoErr(err)

	// Check the entrypoint
	req := httptest.NewRequest("GET", "/bud/view/_index.svelte.js", nil)
	rec := httptest.NewRecorder()
	client.Proxy(rec, req)
	res := rec.Result()
	defer res.Body.Close()
	is.Equal(res.StatusCode, 200)
	is.Equal(len(res.Header), 4)
	is.Equal(res.Header.Get("Content-Type"), "application/javascript")
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
