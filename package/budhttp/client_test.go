package budhttp_test

import (
	"context"
	"testing"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/dag"

	"github.com/livebud/bud/package/genfs"
	"github.com/livebud/bud/package/socket"

	"github.com/livebud/bud/package/log/testlog"

	"github.com/livebud/bud/framework/transform/transformrt"
	"github.com/livebud/bud/framework/view/dom"
	"github.com/livebud/bud/framework/view/nodemodules"
	"github.com/livebud/bud/framework/view/ssr"
	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/pubsub"
	"github.com/livebud/bud/internal/testdir"
	"github.com/livebud/bud/package/budhttp"
	"github.com/livebud/bud/package/budhttp/budsvr"
	"github.com/livebud/bud/package/gomod"
	v8 "github.com/livebud/bud/package/js/v8"
	"github.com/livebud/bud/package/svelte"
)

func loadServer(bus pubsub.Client, dir string) (*budsvr.Server, error) {
	flag := new(framework.Flag)
	log := testlog.New()
	vm, err := v8.Load()
	if err != nil {
		return nil, err
	}
	svelteCompiler, err := svelte.Load(vm)
	if err != nil {
		return nil, err
	}
	transforms, err := transformrt.Default(log, svelteCompiler)
	if err != nil {
		return nil, err
	}
	module, err := gomod.Find(dir)
	if err != nil {
		return nil, err
	}
	budln, err := socket.Listen(":0")
	if err != nil {
		return nil, err
	}
	cache, err := dag.Load(log, ":memory:")
	if err != nil {
		return nil, err
	}
	gfs := genfs.New(cache, module, log)
	gfs.FileServer("bud/view", dom.New(module, transforms))
	gfs.FileServer("bud/node_modules", nodemodules.New(module))
	gfs.FileGenerator("bud/view/_ssr.js", ssr.New(module, transforms))
	return budsvr.New(budln, bus, flag, gfs, log, vm), nil
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
	server.Start(ctx)
	defer server.Close()
	client, err := budhttp.Load(log, server.Address())
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
	server.Start(ctx)
	defer server.Close()
	client, err := budhttp.Load(log, server.Address())
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
	server.Start(ctx)
	defer server.Close()
	client, err := budhttp.Load(log, server.Address())
	is.NoErr(err)
	err = client.Script("script.js", "function a() { return 1 }")
	is.NoErr(err)
	val, err := client.Eval("script.js", "a()")
	is.NoErr(err)
	is.Equal(val, "1")
}
