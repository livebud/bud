package public_test

import (
	"context"
	_ "embed"
	"testing"

	"github.com/livebud/bud/internal/cli"
	"github.com/livebud/bud/internal/cli/testcli"
	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/testdir"
)

func TestNoProject(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	cli := testcli.New(cli.New(dir))
	stdout, stderr, err := cli.Run(ctx)
	is.NoErr(err)
	is.In(stdout.String(), "bud")
	is.Equal(stderr.String(), "")
	is.NoErr(td.NotExists("bud/.app/public/public.go"))
}

func TestEmptyBuild(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	is.NoErr(td.Write(ctx))
	cli := testcli.New(cli.New(dir))
	stdout, stderr, err := cli.Run(ctx, "build")
	is.NoErr(err)
	is.In(stdout.String(), "")
	is.Equal(stderr.String(), "")
	// Empty builds don't generate public files
	is.NoErr(td.NotExists("bud/.app/public/public.go"))
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
	cli := testcli.New(cli.New(dir))
	app, stdout, stderr, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	// /favicon.ico
	res, err := app.Get("/favicon.ico")
	is.NoErr(err)
	is.Equal(200, res.Status())
	is.Equal(res.Body().Bytes(), favicon)
	// Ubuntu CI reports a different MIME type than OSX
	is.In(res.Header("Content-Type"), "image/")
	is.In(res.Header("Content-Type"), "icon")
	// /ga.js
	res, err = app.Get("/ga.js")
	is.NoErr(err)
	is.Equal(200, res.Status())
	is.Equal(res.Body().String(), ga)
	is.In(res.Header("Content-Type"), "javascript")
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
	// Test stdio
	is.Equal(stdout.String(), "")
	is.Equal(stderr.String(), "")
}

func TestPlugin(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Modules["github.com/livebud/bud-test-plugin"] = "v0.0.8"
	is.NoErr(td.Write(ctx))
	cli := testcli.New(cli.New(dir))
	app, stdout, stderr, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.Get("/tailwind/preflight.css")
	is.NoErr(err)
	is.Equal(200, res.Status())
	is.Equal(res.Body().String(), `/* tailwind */`)
	is.Equal(stdout.String(), "")
	is.Equal(stderr.String(), "")
}

func TestGetChangeGet(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.BFiles["public/favicon.ico"] = favicon
	is.NoErr(td.Write(ctx))
	cli := testcli.New(cli.New(dir))
	app, _, _, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.Get("/favicon.ico")
	is.NoErr(err)
	is.Equal(200, res.Status())
	is.Equal(res.Body().Bytes(), favicon)
	is.NoErr(td.Exists("bud/.app/public/public.go"))
	// Favicon2
	favicon2 := []byte{0x00, 0x00, 0x01}
	td.BFiles["public/favicon.ico"] = favicon2
	is.NoErr(td.Write(ctx))
	is.NoErr(td.Exists("bud/.app/public/public.go"))
	res, err = app.Get("/favicon.ico")
	is.NoErr(err)
	is.Equal(200, res.Status())
	is.Equal(res.Body().Bytes(), favicon2)
	// is.Equal(stdout.String(), "")
	// is.Equal(stderr.String(), "")
}

func TestEmbedFavicon(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.BFiles["public/favicon.ico"] = favicon
	is.NoErr(td.Write(ctx))
	cli := testcli.New(cli.New(dir))
	app, stdout, stderr, err := cli.Start(ctx, "run", "--embed")
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
	// Favicon shouldn't have changed
	res, err = app.Get("/favicon.ico")
	is.NoErr(err)
	is.Equal(200, res.Status())
	is.Equal(res.Body().Bytes(), favicon)
	is.Equal(stdout.String(), "")
	is.Equal(stderr.String(), "")
}

func TestAppPluginOverlap(t *testing.T) {
	t.SkipNow()
}

func TestPluginPluginOverlap(t *testing.T) {
	t.SkipNow()
}

//go:embed favicon.ico
var defaultFavicon []byte

//go:embed default.css
var defaultCSS []byte

func TestDefaults(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["view/index.svelte"] = `<h1>hello</h1>`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(cli.New(dir))
	app, stdout, stderr, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	// default favicon
	res, err := app.Get("/favicon.ico")
	is.NoErr(err)
	is.Equal(200, res.Status())
	is.Equal(res.Body().Bytes(), defaultFavicon)
	// default.css
	res, err = app.Get("/default.css")
	is.NoErr(err)
	is.Equal(200, res.Status())
	is.Equal(res.Body().Bytes(), defaultCSS)
	is.Equal(stdout.String(), "")
	is.Equal(stderr.String(), "")
}
