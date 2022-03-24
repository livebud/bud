package transform_test

import (
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/internal/budtest"
	"gitlab.com/mnm/bud/package/modcache"
)

// TODO: re-enable once bud/app is always built
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
	is.NoErr(app.NotExists("bud/.app/transform/transform.go"))
}

func TestSvelteView(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	bud := budtest.New(dir)
	bud.Files["view/index.svelte"] = `<h1>hello world!</h1>`
	project, err := bud.Compile(ctx)
	is.NoErr(err)
	app, err := project.Build(ctx)
	is.NoErr(err)
	is.NoErr(app.Exists("bud/.app/transform/transform.go"))
}

func TestMarkdownPlugin(t *testing.T) {
	t.SkipNow()
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	bud := budtest.New(dir)
	bud.Modules = map[string]modcache.Files{
		"gitlab.com/mnm/bud-markdown@v0.0.1": modcache.Files{
			"transform/markdown/transform.go": `
				package markdown
				import "gitlab.com/mnm/bud/runtime/transform"
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
	is.NoErr(app.Exists("bud/transform/transform.go"))
	is.NoErr(app.Exists("bud/main.go"))
	fmt.Println("starting")
	server, err := app.Start(ctx)
	is.NoErr(err)
	defer server.Close()
	res, err := server.Get("/")
	is.NoErr(err)
	defer res.Body.Close()
	is.Equal(200, res.StatusCode)
	body, err := io.ReadAll(res.Body)
	is.NoErr(err)
	is.Equal(`<h1>hello</h1>`, string(body))
}
