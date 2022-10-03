package budhttp_test

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/livebud/bud/package/budfs"
	"github.com/livebud/bud/package/log/testlog"

	"github.com/livebud/bud/framework/transform/transformrt"
	"github.com/livebud/bud/framework/view/dom"
	"github.com/livebud/bud/framework/view/ssr"
	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/pubsub"
	"github.com/livebud/bud/internal/testdir"
	"github.com/livebud/bud/internal/versions"
	"github.com/livebud/bud/package/budhttp"
	"github.com/livebud/bud/package/budhttp/budsvr"
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
	bfs.FileServer("bud/view", dom.New(module, transforms.DOM))
	bfs.FileServer("bud/node_modules", dom.NodeModules(module))
	bfs.FileGenerator("bud/view/_ssr.js", ssr.New(module, transforms.SSR))
	handler := budsvr.New(bfs, bus, log, vm)
	return httptest.NewServer(handler), nil
}

func TestOpen(t *testing.T) {
	ctx := context.Background()
	is := is.New(t)
	log := testlog.New()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["view/index.svelte"] = `
		<script>
			export let _string = "cupcake"
		</script>
		<h1>Hello, {_string}!</h1>
	`
	td.NodeModules["svelte"] = versions.Svelte
	td.NodeModules["livebud"] = "*"
	is.NoErr(td.Write(ctx))
	bus := pubsub.New()
	server, err := loadServer(bus, dir)
	is.NoErr(err)
	defer server.Close()
	client, err := budhttp.Load(log, server.URL)
	is.NoErr(err)

	// Check the entrypoint
	file, err := client.Open("bud/view/_index.svelte.js")
	is.NoErr(err)
	defer file.Close()
	code, err := io.ReadAll(file)
	is.NoErr(err)
	is.In(string(code), `"view/index.svelte"`)
	is.In(string(code), `"/bud/view/index.svelte"`)
	stat, err := file.Stat()
	is.NoErr(err)
	is.Equal(stat.Name(), "_index.svelte.js")
	is.Equal(stat.Size(), int64(3690))
	is.Equal(stat.Mode(), os.FileMode(0))
	is.Equal(stat.ModTime(), time.Time{})
	is.Equal(stat.IsDir(), false)
	is.Equal(stat.Sys(), nil)

	// Check the component
	file, err = client.Open("bud/view/index.svelte")
	is.NoErr(err)
	defer file.Close()
	stat, err = file.Stat()
	is.NoErr(err)
	is.Equal(stat.Name(), "index.svelte")
	is.Equal(stat.Size(), int64(3124))
	is.Equal(stat.Mode(), os.FileMode(0))
	is.Equal(stat.ModTime(), time.Time{})
	is.Equal(stat.IsDir(), false)
	is.Equal(stat.Sys(), nil)

	// Check the node_modules
	file, err = client.Open("bud/node_modules/svelte/internal")
	is.NoErr(err)
	defer file.Close()
	stat, err = file.Stat()
	is.NoErr(err)
	is.Equal(stat.Name(), "internal")
	is.Equal(stat.Size(), int64(56452))
	is.Equal(stat.Mode(), os.FileMode(0))
	is.Equal(stat.ModTime(), time.Time{})
	is.Equal(stat.IsDir(), false)
	is.Equal(stat.Sys(), nil)
}

func TestOpen404(t *testing.T) {
	ctx := context.Background()
	is := is.New(t)
	log := testlog.New()
	dir := t.TempDir()
	td := testdir.New(dir)
	is.NoErr(td.Write(ctx))
	bus := pubsub.New()
	server, err := loadServer(bus, dir)
	is.NoErr(err)
	defer server.Close()
	client, err := budhttp.Load(log, server.URL)
	is.NoErr(err)
	file, err := client.Open("public/favicon.ico")
	is.True(errors.Is(err, fs.ErrNotExist))
	is.Equal(file, nil)
}

func TestEvents(t *testing.T) {
	ctx := context.Background()
	is := is.New(t)
	log := testlog.New()
	dir := t.TempDir()
	td := testdir.New(dir)
	is.NoErr(td.Write(ctx))
	ps := pubsub.New()
	server, err := loadServer(ps, dir)
	is.NoErr(err)
	defer server.Close()
	client, err := budhttp.Load(log, server.URL)
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

func TestScript(t *testing.T) {
	ctx := context.Background()
	is := is.New(t)
	log := testlog.New()
	dir := t.TempDir()
	td := testdir.New(dir)
	is.NoErr(td.Write(ctx))
	ps := pubsub.New()
	server, err := loadServer(ps, dir)
	is.NoErr(err)
	defer server.Close()
	client, err := budhttp.Load(log, server.URL)
	is.NoErr(err)
	err = client.Script("script.js", "function a() { return 1 }")
	is.NoErr(err)
	err = client.Script("script.js", "function b() { return 1")
	is.True(err != nil)
	is.In(err.Error(), "SyntaxError: Unexpected end of input")
}

func TestScriptEval(t *testing.T) {
	ctx := context.Background()
	is := is.New(t)
	log := testlog.New()
	dir := t.TempDir()
	td := testdir.New(dir)
	is.NoErr(td.Write(ctx))
	ps := pubsub.New()
	server, err := loadServer(ps, dir)
	is.NoErr(err)
	defer server.Close()
	client, err := budhttp.Load(log, server.URL)
	is.NoErr(err)
	err = client.Script("script.js", "function a() { return 1 }")
	is.NoErr(err)
	val, err := client.Eval("script.js", "a()")
	is.NoErr(err)
	is.Equal(val, "1")
}
