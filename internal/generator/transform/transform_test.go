package transform_test

import (
	"context"
	"testing"

	"github.com/livebud/bud/internal/budtest"
	"github.com/livebud/bud/package/modcache"
	"github.com/matryer/is"
)

func TestEmpty(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	bud := budtest.New(dir)
	project, err := bud.Compile(ctx)
	is.NoErr(err)
	app, err := project.Build(ctx)
	is.NoErr(err)
	is.NoErr(app.Exists("bud/.cli/transform/transform.go"))
}

func TestSvelteView(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	bud := budtest.New(dir)
	bud.Files["view/index.svelte"] = `<h1>hello</h1>`
	project, err := bud.Compile(ctx)
	is.NoErr(err)
	app, err := project.Build(ctx)
	is.NoErr(err)
	is.NoErr(app.Exists("bud/.cli/transform/transform.go"))
}

func TestMarkdownPlugin(t *testing.T) {
	t.SkipNow()
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	bud := budtest.New(dir)
	bud.Modules = map[string]modcache.Files{
		"github.com/livebud/bud-markdown@v0.0.1": modcache.Files{
			"transform/markdown/transform.go": `
				package markdown
				import "github.com/livebud/bud/runtime/transform"
				type Transform struct {}
				func (t *Transform) MdToSvelte(file *transform.File) error {
					file.Code = "<h1>hello</h1>"
				}
			`,
		},
	}
	bud.Files["view/index.md"] = `# hello`
	project, err := bud.Compile(ctx)
	is.NoErr(err)
	app, err := project.Build(ctx)
	is.NoErr(err)
	is.NoErr(app.Exists("bud/.cli/transform/transform.go"))
	is.NoErr(app.Exists("bud/.app/main.go"))
	server, err := app.Start(ctx)
	is.NoErr(err)
	defer server.Close()
	res, err := server.Get("/")
	is.NoErr(err)
	defer res.Body.Close()
	// HTML response
	is.NoErr(res.ExpectHeaders(`
		HTTP/1.1 200 OK
		Content-Type: text/html
		Date: Fri, 31 Dec 2021 00:00:00 GMT
	`))
	is.NoErr(res.ContainsBody(`<h1>hello</h1>`))
}
