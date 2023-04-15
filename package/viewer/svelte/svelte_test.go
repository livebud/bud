package svelte_test

import (
	"context"
	"io"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/versions"
	"github.com/livebud/bud/package/es"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log/testlog"
	"github.com/livebud/bud/package/router"
	"github.com/livebud/bud/package/testdir"
	"github.com/livebud/bud/package/transpiler"
	"github.com/livebud/bud/package/viewer"
	"github.com/livebud/bud/package/viewer/svelte"
	"github.com/livebud/js"
	"github.com/livebud/js/goja"
)

func TestPage(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	log := testlog.New()
	td, err := testdir.Load()
	is.NoErr(err)
	td.NodeModules["livebud"] = "*"
	td.NodeModules["svelte"] = versions.Svelte
	td.Files["index.svelte"] = `
		<script>
			export let planet = 'Mars'
		</script>
		<h1>Hello {planet}!</h1>
	`
	is.NoErr(td.Write(ctx))
	module, err := gomod.Find(td.Directory())
	is.NoErr(err)
	pages, err := viewer.Find(module)
	is.NoErr(err)
	js := goja.New(&js.Console{
		Log:   os.Stdout,
		Error: os.Stderr,
	})
	flag := &framework.Flag{}
	esb := es.New(flag, log, module)
	svelteCompiler, err := svelte.Load(flag, js)
	is.NoErr(err)
	tr := transpiler.New()
	tr.Add(".svelte", ".ssr.js", func(file *transpiler.File) error {
		ssr, err := svelteCompiler.SSR(ctx, file.Path(), file.Data)
		if err != nil {
			return err
		}
		file.Data = []byte(ssr.JS)
		return nil
	})
	tr.Add(".svelte", ".dom.js", func(file *transpiler.File) error {
		dom, err := svelteCompiler.DOM(ctx, file.Path(), file.Data)
		if err != nil {
			return err
		}
		file.Data = []byte(dom.JS)
		return nil
	})
	viewer := svelte.New(esb, flag, js, log, module, pages, tr)
	html, err := viewer.Render(ctx, "index", map[string]interface{}{
		"index": map[string]interface{}{
			"planet": "Earth",
		},
	})
	is.NoErr(err)
	is.Equal(string(html), `<h1>Hello Earth!</h1>`)
	router := router.New()
	is.NoErr(viewer.Mount(router))

	// Entry
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/view/index.svelte.entry.js", nil)
	router.ServeHTTP(rec, req)
	res := rec.Result()
	body, err := io.ReadAll(res.Body)
	is.NoErr(err)
	is.In(string(body), `"/node_modules/svelte/internal"`)
	is.In(string(body), `"/view/index.svelte.js"`)
	is.In(string(body), `mount({`)
	is.In(string(body), `key: "index"`)
	is.Equal(res.StatusCode, 200)

	// View
	rec = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/view/index.svelte.js", nil)
	router.ServeHTTP(rec, req)
	res = rec.Result()
	body, err = io.ReadAll(res.Body)
	is.NoErr(err)
	is.In(string(body), `"/node_modules/svelte/internal"`)
	is.In(string(body), `"Hello "`)
	is.In(string(body), `"planet"`)
	is.Equal(res.StatusCode, 200)
}

func TestLayout(t *testing.T) {
	t.Skip("TODO: test layouts")
}

func TestRenderError(t *testing.T) {
	t.Skip("TODO: test render error")
}

func TestBundle(t *testing.T) {
	t.Skip("TODO: test bundling")
}
