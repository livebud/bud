package public_test

import (
	"context"
	"testing"
	"time"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/testcli"
	"github.com/livebud/bud/internal/testdir"
)

func TestNoProject(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	is.NoErr(td.NotExists("bud/internal/web/public"))
	result, err := cli.Run(ctx)
	is.NoErr(err)
	is.In(result.Stdout(), "bud")
	is.Equal(result.Stderr(), "")
	is.NoErr(td.NotExists("bud/internal/web/public"))
}

func TestEmptyBuild(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	is.NoErr(td.NotExists("bud/internal/web/public"))
	result, err := cli.Run(ctx, "build")
	is.NoErr(err)
	is.In(result.Stdout(), "")
	is.Equal(result.Stderr(), "")
	// Empty builds don't generate public files
	is.NoErr(td.NotExists("bud/internal/web/public"))
}

// Pulled from: https://github.com/mathiasbynens/small
// Built with: xxd -i small.ico
var favicon = []byte{
	0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x01, 0x01, 0x00, 0x00, 0x01, 0x00,
	0x18, 0x00, 0x30, 0x00, 0x00, 0x00, 0x16, 0x00, 0x00, 0x00, 0x28, 0x00,
	0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x01, 0x00,
	0x18, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00,
}

// Small valid gif: https://github.com/mathiasbynens/small/blob/master/gif.gif
var gif = []byte{
	0x47, 0x49, 0x46, 0x38, 0x39, 0x61, 0x01,
	0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x3b,
}

func TestPublic(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.BFiles["public/favicon.ico"] = favicon
	ga := `function ga(track){}`
	td.Files["public/ga.js"] = ga
	css := `* { box-sizing: border-box; }`
	td.Files["public/normalize/normalize.css"] = css
	td.BFiles["public/lol.gif"] = gif
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	// /favicon.ico
	res, err := app.Get("/favicon.ico")
	is.NoErr(err)
	is.Equal(200, res.Status(), "unable to user-defined /favicon.ico")
	is.Equal(res.Body().Bytes(), favicon)
	// Ubuntu CI reports a different MIME type than OSX
	is.In(res.Header("Content-Type"), "image/")
	is.In(res.Header("Content-Type"), "icon")
	// /ga.js
	res, err = app.Get("/ga.js")
	is.NoErr(err)
	is.Equal(200, res.Status())
	is.Equal(res.Body().String(), ga)
	is.In(res.Header("Content-Type"), "/javascript")
	// /normalize/normalize.css
	res, err = app.Get("/normalize/normalize.css")
	is.NoErr(err)
	is.Equal(200, res.Status())
	is.Equal(res.Body().String(), css)
	is.In(res.Header("Content-Type"), "css")
	// /normalize/normalize.css
	res, err = app.Get("/lol.gif")
	is.NoErr(err)
	is.Equal(200, res.Status())
	is.Equal(res.Body().Bytes(), gif)
	is.In(res.Header("Content-Type"), "image/")
	is.In(res.Header("Content-Type"), "gif")
}

func TestPlugin(t *testing.T) {
	t.SkipNow()
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Modules["github.com/livebud/bud-test-plugin"] = "v0.0.9"
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.Get("/tailwind/preflight.css")
	is.NoErr(err)
	is.Equal(200, res.Status(), "unable to get tailwind plugin")
	is.Equal(res.Body().String(), `/* tailwind */`)
}

func TestGetChangeGet(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.BFiles["public/favicon.ico"] = favicon
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.Get("/favicon.ico")
	is.NoErr(err)
	is.Equal(200, res.Status())
	is.Equal(res.Body().Bytes(), favicon)
	is.NoErr(td.Exists("bud/internal/web/public/public.go"))
	// Favicon2
	favicon2 := []byte{0x00, 0x00, 0x01}
	td.BFiles["public/favicon.ico"] = favicon2
	is.NoErr(td.Write(ctx))
	readyCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	is.NoErr(app.Ready(readyCtx))
	cancel()
	is.NoErr(td.Exists("bud/internal/web/public/public.go"))
	res, err = app.Get("/favicon.ico")
	is.NoErr(err)
	is.Equal(200, res.Status())
	is.Equal(res.Body().Bytes(), favicon2)
	// is.Equal(result.Stdout(), "")
	// is.Equal(result.Stderr(), "")
}

func TestEmbedFavicon(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.BFiles["public/favicon.ico"] = favicon
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run", "--embed", "--hot=false")
	is.NoErr(err)
	defer app.Close()
	res, err := app.Get("/favicon.ico")
	is.NoErr(err)
	is.Equal(200, res.Status())
	is.Equal(res.Body().Bytes(), favicon)
	// Replace favicon
	favicon2 := []byte{0x00, 0x00, 0x01}
	td.BFiles["public/favicon.ico"] = favicon2
	is.NoErr(td.Write(ctx))
	readyCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	is.NoErr(app.Ready(readyCtx))
	cancel()
	// Favicon shouldn't have changed because non-Go files don't trigger
	// full rebuilds and server restarts
	res, err = app.Get("/favicon.ico")
	is.NoErr(err)
	is.Equal(200, res.Status())
	is.Equal(res.Body().Bytes(), favicon)
	is.NoErr(app.Close())
}
