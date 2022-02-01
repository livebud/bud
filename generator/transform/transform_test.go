package transform_test

import (
	"io"
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/internal/test"
	"gitlab.com/mnm/bud/pkg/modcache"
)

func TestEmpty(t *testing.T) {
	is := is.New(t)
	generator := test.Generator(t)
	app, err := generator.Generate()
	is.NoErr(err)
	is.Equal(false, app.Exists("bud/transform/transform.go"))
}

func TestSvelteView(t *testing.T) {
	is := is.New(t)
	generator := test.Generator(t)
	generator.Files["view/index.svelte"] = []byte(`<h1>hello world!</h1>`)
	app, err := generator.Generate()
	is.NoErr(err)
	is.Equal(true, app.Exists("bud/transform/transform.go"))
	is.Equal(true, app.Exists("bud/main.go"))
}

func TestMarkdownPlugin(t *testing.T) {
	t.SkipNow()
	is := is.New(t)
	generator := test.Generator(t)
	generator.Modules = map[string]modcache.Files{
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
	generator.Files["view/index.md"] = []byte(`# hello`)
	app, err := generator.Generate()
	is.NoErr(err)
	is.Equal(true, app.Exists("bud/transform/transform.go"))
	is.Equal(true, app.Exists("bud/main.go"))
	server, err := app.Start()
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
