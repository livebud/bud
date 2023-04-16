package svelte_test

import (
	"context"
	"errors"
	"io"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/is"
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

func loadViewer(dir string) (*svelte.Viewer, error) {
	log := testlog.New()
	module, err := gomod.Find(dir)
	if err != nil {
		return nil, err
	}
	pages, err := viewer.Find(module)
	if err != nil {
		return nil, err
	}
	js := goja.New(&js.Console{
		Log:   os.Stdout,
		Error: os.Stderr,
	})
	flag := &framework.Flag{}
	esb := es.New(flag, log)
	svelteCompiler, err := svelte.Load(flag, js)
	if err != nil {
		return nil, err
	}
	tr := transpiler.New()
	tr.Add(".svelte", ".ssr.js", func(ctx context.Context, file *transpiler.File) error {
		ssr, err := svelteCompiler.SSR(ctx, file.Path(), file.Data)
		if err != nil {
			return err
		}
		file.Data = []byte(ssr.JS)
		return nil
	})
	tr.Add(".svelte", ".dom.js", func(ctx context.Context, file *transpiler.File) error {
		dom, err := svelteCompiler.DOM(ctx, file.Path(), file.Data)
		if err != nil {
			return err
		}
		file.Data = []byte(dom.JS)
		return nil
	})
	viewer := svelte.New(esb, flag, js, log, module, pages, tr)
	return viewer, nil
}

func TestPage(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	td, err := testdir.Load()
	is.NoErr(err)
	td.Files["index.svelte"] = `
		<script>
			export let planet = 'Mars'
		</script>
		<h1>Hello {planet}!</h1>
	`
	is.NoErr(td.Write(ctx))
	viewer, err := loadViewer(td.Directory())
	is.NoErr(err)
	html, err := viewer.Render(ctx, "index", map[string]interface{}{
		"index": map[string]interface{}{
			"planet": "Earth",
		},
	})
	is.NoErr(err)
	is.Equal(string(html), `<h1>Hello Earth!</h1>`)

	// Mount the client
	router := router.New()
	is.NoErr(viewer.Mount(router))

	// Entry
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/view/index.svelte.entry.js", nil)
	router.ServeHTTP(rec, req)
	res := rec.Result()
	body, err := io.ReadAll(res.Body)
	is.NoErr(err)
	is.In(string(body), `"https://esm.run/svelte@3.47.0/internal"`)
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
	is.In(string(body), `"https://esm.run/svelte@3.47.0/internal"`)
	is.In(string(body), `"Hello "`)
	is.In(string(body), `"planet"`)
	is.Equal(res.StatusCode, 200)
}

func TestLayout(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	td, err := testdir.Load()
	is.NoErr(err)
	td.Files["show.svelte"] = `
		<script>
			export let planet = 'Mars'
		</script>
		<h1>Hello {planet}!</h1>
	`
	td.Files["layout.svelte"] = `
		<script>
			export let title = 'default'
		</script>
		<html>
		<head>
			<title>{title}</title>
		</head>
		<body>
			<slot />
		</body>
		</html>
	`
	is.NoErr(td.Write(ctx))
	viewer, err := loadViewer(td.Directory())
	is.NoErr(err)
	html, err := viewer.Render(ctx, "show", map[string]interface{}{
		"layout": map[string]interface{}{
			"title": "Hello",
		},
		"show": map[string]interface{}{
			"planet": "Earth",
		},
	})
	is.NoErr(err)
	is.In(string(html), `<head><title>Hello</title></head>`)
	is.In(string(html), `<script id="bud_state" type="text/template">{"props":{"show":{"planet":"Earth"}}}</script>`)
	is.In(string(html), `<h1>Hello Earth!</h1>`)

	// Mount the client
	router := router.New()
	is.NoErr(viewer.Mount(router))

	// Entry
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/view/show.svelte.entry.js", nil)
	router.ServeHTTP(rec, req)
	res := rec.Result()
	body, err := io.ReadAll(res.Body)
	is.NoErr(err)
	is.In(string(body), `"https://esm.run/svelte@3.47.0/internal"`)
	is.In(string(body), `"/view/show.svelte.js"`)
	is.In(string(body), `mount({`)
	is.In(string(body), `key: "show"`)
	is.Equal(res.StatusCode, 200)

	// View
	rec = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/view/show.svelte.js", nil)
	router.ServeHTTP(rec, req)
	res = rec.Result()
	body, err = io.ReadAll(res.Body)
	is.NoErr(err)
	is.In(string(body), `"https://esm.run/svelte@3.47.0/internal"`)
	is.In(string(body), `"Hello "`)
	is.In(string(body), `"planet"`)
	is.Equal(res.StatusCode, 200)

	// Shouldn't expose the layout
	rec = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/view/layout.svelte.js", nil)
	router.ServeHTTP(rec, req)
	res = rec.Result()
	is.Equal(res.StatusCode, 404)
}

func TestRenderError(t *testing.T) {
	t.Skip("not done yet")
	is := is.New(t)
	ctx := context.Background()
	td, err := testdir.Load()
	is.NoErr(err)
	td.Files["index.svelte"] = `
		<script>
			export let planet = 'Mars'
		</script>
		<h1>Hello {planet}!</h1>
	`
	td.Files["layout.svelte"] = `
		<html><slot /></html>
	`
	td.Files["error.svelte"] = `
		<script>
			export let message = ''
		</script>
		<div class="error">{message}</div>
	`
	is.NoErr(td.Write(ctx))
	viewer, err := loadViewer(td.Directory())
	is.NoErr(err)
	html := viewer.RenderError(ctx, "index", map[string]interface{}{
		"index": map[string]interface{}{
			"Planet": "Earth",
		},
	}, errors.New("some error"))
	is.NoErr(err)
	is.Equal(string(html), `<html><div class="error">some error</div></html>`)
}

func TestBundle(t *testing.T) {
	t.Skip("TODO: test bundling")
}
