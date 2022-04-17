package public_test

import (
	"bytes"
	"context"
	_ "embed"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/matryer/is"
	"github.com/livebud/bud/internal/budtest"
	"github.com/livebud/bud/package/modcache"
)

// TODO: bud/.app/main.go should always be generated, but public will not exist
// if there are no files
func TestEmpty(t *testing.T) {
	t.SkipNow()
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	bud := budtest.New(dir)
	project, err := bud.Compile(ctx)
	is.NoErr(err)
	app, err := project.Build(ctx)
	is.NoErr(err)
	is.NoErr(app.NotExists("bud/.app/public/public.go"))
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

func TestFavicon(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	bud := budtest.New(dir)
	bud.BFiles["public/favicon.ico"] = favicon
	project, err := bud.Compile(ctx)
	is.NoErr(err)
	app, err := project.Build(ctx)
	is.NoErr(err)
	is.NoErr(app.Exists("bud/.app/public/public.go"))
	server, err := app.Start(ctx)
	is.NoErr(err)
	defer server.Close()
	res, err := server.Get("/favicon.ico")
	is.NoErr(err)
	defer res.Body.Close()
	is.Equal(200, res.StatusCode)
	body, err := io.ReadAll(res.Body)
	is.NoErr(err)
	is.True(bytes.Equal(favicon, body))
}

func TestNested(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	bud := budtest.New(dir)
	css := `* { box-sizing: border-box; }`
	bud.Files["public/normalize/normalize.css"] = css
	project, err := bud.Compile(ctx)
	is.NoErr(err)
	app, err := project.Build(ctx)
	is.NoErr(err)
	is.NoErr(app.Exists("bud/.app/public/public.go"))
	server, err := app.Start(ctx)
	is.NoErr(err)
	defer server.Close()
	res, err := server.Get("/normalize/normalize.css")
	is.NoErr(err)
	defer res.Body.Close()
	is.Equal(200, res.StatusCode)
	body, err := io.ReadAll(res.Body)
	is.NoErr(err)
	is.Equal(string(body), css)
}

func TestPlugin(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	bud := budtest.New(dir)
	preflight := `/* tailwind */`
	bud.Modules = map[string]modcache.Files{
		"github.com/livebud/bud-tailwind@v0.0.1": modcache.Files{
			"public/tailwind/preflight.css": preflight,
		},
	}
	project, err := bud.Compile(ctx)
	is.NoErr(err)
	app, err := project.Build(ctx)
	is.NoErr(err)
	is.NoErr(app.Exists("bud/.app/public/public.go"))
	server, err := app.Start(ctx)
	is.NoErr(err)
	defer server.Close()
	res, err := server.Get("/tailwind/preflight.css")
	is.NoErr(err)
	defer res.Body.Close()
	is.Equal(200, res.StatusCode)
	body, err := io.ReadAll(res.Body)
	is.NoErr(err)
	is.Equal(preflight, string(body))
}

func TestGetChangeGet(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	bud := budtest.New(dir)
	bud.BFiles["public/favicon.ico"] = favicon
	project, err := bud.Compile(ctx)
	is.NoErr(err)
	app, err := project.Build(ctx)
	is.NoErr(err)
	is.NoErr(app.Exists("bud/.app/public/public.go"))
	server, err := app.Start(ctx)
	is.NoErr(err)
	defer server.Close()
	res, err := server.Get("/favicon.ico")
	is.NoErr(err)
	defer res.Body.Close()
	is.Equal(200, res.StatusCode)
	body, err := io.ReadAll(res.Body)
	is.NoErr(err)
	is.True(bytes.Equal(favicon, body))
	// Favicon2
	favicon2 := []byte{0x00, 0x00, 0x01}
	bud.BFiles["public/favicon.ico"] = favicon2
	err = project.Rewrite()
	is.NoErr(err)
	// Rebuild
	app, err = project.Build(ctx)
	is.NoErr(err)
	is.NoErr(app.Exists("bud/.app/public/public.go"))
	err = server.Restart(ctx)
	is.NoErr(err)
	res, err = server.Get("/")
	is.NoErr(err)
	res, err = server.Get("/favicon.ico")
	is.NoErr(err)
	defer res.Body.Close()
	is.Equal(200, res.StatusCode)
	body, err = io.ReadAll(res.Body)
	is.NoErr(err)
	is.True(bytes.Equal(favicon2, body))
}

func TestEmbedFavicon(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	bud := budtest.New(dir)
	bud.Flag.Embed = true
	bud.BFiles["public/favicon.ico"] = favicon
	project, err := bud.Compile(ctx)
	is.NoErr(err)
	app, err := project.Build(ctx)
	is.NoErr(err)
	is.NoErr(app.Exists("bud/.app/public/public.go"))
	// Remove file to ensure it's been embedded
	is.NoErr(os.Remove(filepath.Join(dir, "public/favicon.ico")))
	// Try requesting the favicon, it should be in memory now
	server, err := app.Start(ctx)
	is.NoErr(err)
	defer server.Close()
	res, err := server.Get("/favicon.ico")
	is.NoErr(err)
	defer res.Body.Close()
	is.Equal(200, res.StatusCode)
	body, err := io.ReadAll(res.Body)
	is.NoErr(err)
	is.True(bytes.Equal(favicon, body))
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
	bud := budtest.New(dir)
	bud.Files["view/index.svelte"] = `<h1>hello</h1>`
	project, err := bud.Compile(ctx)
	is.NoErr(err)
	app, err := project.Build(ctx)
	is.NoErr(err)
	is.NoErr(app.Exists("bud/.app/public/public.go"))
	server, err := app.Start(ctx)
	is.NoErr(err)
	defer server.Close()
	// favicon.ico
	res, err := server.Get("/favicon.ico")
	is.NoErr(err)
	defer res.Body.Close()
	is.Equal(200, res.StatusCode)
	body, err := io.ReadAll(res.Body)
	is.NoErr(err)
	is.True(bytes.Equal(defaultFavicon, body))
	// default.css
	res, err = server.Get("/default.css")
	is.NoErr(err)
	defer res.Body.Close()
	is.Equal(200, res.StatusCode)
	body, err = io.ReadAll(res.Body)
	is.NoErr(err)
	is.True(bytes.Equal(defaultCSS, body))
}

func TestDefaultsPublic(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	bud := budtest.New(dir)
	ga := `function ga(track){}`
	bud.Files["public/ga.js"] = ga
	project, err := bud.Compile(ctx)
	is.NoErr(err)
	app, err := project.Build(ctx)
	is.NoErr(err)
	is.NoErr(app.Exists("bud/.app/public/public.go"))
	server, err := app.Start(ctx)
	is.NoErr(err)
	defer server.Close()
	// favicon.ico
	res, err := server.Get("/ga.js")
	is.NoErr(err)
	defer res.Body.Close()
	is.Equal(200, res.StatusCode)
	body, err := io.ReadAll(res.Body)
	is.NoErr(err)
	is.Equal(ga, string(body))
	// favicon.ico
	res, err = server.Get("/favicon.ico")
	is.NoErr(err)
	defer res.Body.Close()
	is.Equal(200, res.StatusCode)
	body, err = io.ReadAll(res.Body)
	is.NoErr(err)
	is.True(bytes.Equal(defaultFavicon, body))
	// default.css
	res, err = server.Get("/default.css")
	is.NoErr(err)
	defer res.Body.Close()
	is.Equal(200, res.StatusCode)
	body, err = io.ReadAll(res.Body)
	is.NoErr(err)
	is.True(bytes.Equal(defaultCSS, body))
}
