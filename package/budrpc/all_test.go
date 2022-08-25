package budrpc_test

import (
	"context"
	"encoding/json"
	"io/fs"
	"net"
	"path/filepath"
	"testing"

	"github.com/livebud/bud/framework/transform/transformrt"
	"github.com/livebud/bud/framework/view/dom"
	"github.com/livebud/bud/framework/view/ssr"
	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/pubsub"
	"github.com/livebud/bud/internal/testdir"
	"github.com/livebud/bud/internal/versions"
	"github.com/livebud/bud/package/budrpc"
	"github.com/livebud/bud/package/budrpc/budsvr"
	"github.com/livebud/bud/package/gomod"
	v8 "github.com/livebud/bud/package/js/v8"
	"github.com/livebud/bud/package/log/testlog"
	"github.com/livebud/bud/package/overlay"
	"github.com/livebud/bud/package/svelte"
	"golang.org/x/sync/errgroup"
)

func loadServer(bus pubsub.Client, dir string) (*testServer, error) {
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
	ln, err := net.Listen("unix", filepath.Join(dir, "unix.socket"))
	if err != nil {
		return nil, err
	}
	genfs, err := overlay.Load(log, module)
	if err != nil {
		return nil, err
	}
	genfs.FileServer("bud/view", dom.New(module, transforms.DOM))
	genfs.FileServer("bud/node_modules", dom.NodeModules(module))
	genfs.FileGenerator("bud/view/_ssr.js", ssr.New(module, transforms.SSR))
	server := budsvr.New(bus, genfs, log, vm)
	eg := new(errgroup.Group)
	eg.Go(func() error { return server.Serve(ln) })
	return &testServer{eg, ln}, nil
}

type testServer struct {
	eg *errgroup.Group
	ln net.Listener
}

func (t *testServer) URL() string {
	return t.ln.Addr().String()
}

func (t *testServer) Close() error {
	if err := t.ln.Close(); err != nil {
		return err
	}
	return t.eg.Wait()
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
	client, err := budrpc.Dial(ctx, server.URL())
	is.NoErr(err)
	res, err := client.Render("/", json.RawMessage(`{"_string": "marshmallow"}`))
	is.NoErr(err)
	is.Equal(res.Status, 200)
	is.Equal(len(res.Headers), 1)
	is.Equal(res.Headers["Content-Type"], "text/html")
	is.In(res.Body, `<h1>Hello, marshmallow!</h1>`)
	is.NoErr(server.Close())
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
	client, err := budrpc.Dial(ctx, server.URL())
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

func TestOpenFile(t *testing.T) {
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
	client, err := budrpc.Dial(ctx, server.URL())
	is.NoErr(err)
	// Check the entry
	code, err := fs.ReadFile(client, "bud/view/_index.svelte.js")
	is.NoErr(err)
	is.In(string(code), `Hello, `)
	is.In(string(code), `cupcake`)
	// Check the component
	code, err = fs.ReadFile(client, "bud/view/index.svelte")
	is.NoErr(err)
	is.In(string(code), `Hello, `)
	is.In(string(code), `cupcake`)
	// Check the node_module
	code, err = fs.ReadFile(client, "bud/node_modules/svelte/internal")
	is.NoErr(err)
	is.In(string(code), `function element(`)
	is.In(string(code), `function text(`)
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
	client, err := budrpc.Dial(ctx, server.URL())
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
