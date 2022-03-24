package public_test

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/internal/budtest"
	"gitlab.com/mnm/bud/package/modcache"
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
		"gitlab.com/mnm/bud-tailwind@v0.0.1": modcache.Files{
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

func TestAppPluginOverlap(t *testing.T) {
	t.SkipNow()
}

func TestPluginPluginOverlap(t *testing.T) {
	t.SkipNow()
}
