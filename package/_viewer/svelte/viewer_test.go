package svelte_test

import (
	"context"
	"io/fs"
	"os"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/es"
	"github.com/livebud/bud/package/viewer"
	"github.com/livebud/bud/package/viewer/svelte"
	"github.com/livebud/bud/package/virtual"
	"github.com/livebud/bud/runtime/transpiler"
	"github.com/livebud/js"
	v8 "github.com/livebud/js/v8"
)

func loadViewer(dir string, flag *framework.Flag, fsys fs.FS, pages map[string]*viewer.Page) (*svelte.Viewer, error) {
	esbuilder := es.New(dir)
	js, err := v8.Load(&js.Console{
		Log:   os.Stdout,
		Error: os.Stderr,
	})
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	compiler, err := svelte.Load(ctx, js)
	if err != nil {
		return nil, err
	}
	tr := transpiler.New()
	tr.Add(".svelte", ".ssr.js", func(file *transpiler.File) error {
		ssr, err := compiler.SSR(ctx, file.Path(), file.Data)
		if err != nil {
			return err
		}
		file.Data = []byte(ssr.JS)
		return nil
	})
	viewer := svelte.New(esbuilder, flag, fsys, js, tr, pages)
	return viewer, nil
}

func TestEntryNoProps(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	fsys := fstest.MapFS{
		"index.svelte": &fstest.MapFile{Data: []byte(`
			<script>
				export let planet = 'Earth'
			</script>
			<h1>Hello {planet}!</h1>
		`)},
	}
	flag := &framework.Flag{}
	pages, err := viewer.Find(fsys)
	is.NoErr(err)
	svelte, err := loadViewer(dir, flag, fsys, pages)
	is.NoErr(err)
	viewers := viewer.Viewers{
		"index": svelte,
	}
	ctx := context.Background()
	html, err := viewers.Render(ctx, "index", nil)
	is.NoErr(err)
	is.True(strings.Contains(string(html), `<h1>Hello Earth!</h1>`))
}

func TestEntryProps(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	fsys := fstest.MapFS{
		"index.svelte": &fstest.MapFile{Data: []byte(`
			<script>
				export let planet = 'Earth'
			</script>
			<h1>Hello {planet}!</h1>
		`)},
	}
	flag := &framework.Flag{}
	pages, err := viewer.Find(fsys)
	is.NoErr(err)
	svelte, err := loadViewer(dir, flag, fsys, pages)
	is.NoErr(err)
	viewers := viewer.Viewers{
		"index": svelte,
	}
	ctx := context.Background()
	html, err := viewers.Render(ctx, "index", map[string]interface{}{
		"index": map[string]interface{}{
			"planet": "Mars",
		},
	})
	is.NoErr(err)
	is.True(strings.Contains(string(html), `<h1>Hello Mars!</h1>`))
}

func TestFrameProps(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	fsys := fstest.MapFS{
		"frame.svelte": &fstest.MapFile{Data: []byte(`
			<script>
				export let id = '123'
			</script>
			<article id="slug-{id}"><slot /></article>
		`)},
		"index.svelte": &fstest.MapFile{Data: []byte(`
			<script>
				export let planet = 'Earth'
			</script>
			<h1>Hello {planet}!</h1>
		`)},
	}
	flag := &framework.Flag{}
	pages, err := viewer.Find(fsys)
	is.NoErr(err)
	svelte, err := loadViewer(dir, flag, fsys, pages)
	is.NoErr(err)
	viewers := viewer.Viewers{
		"index": svelte,
	}
	ctx := context.Background()
	html, err := viewers.Render(ctx, "index", map[string]interface{}{
		"index": map[string]interface{}{
			"planet": "Mars",
		},
		"frame": map[string]interface{}{
			"id": 456,
		},
	})
	is.NoErr(err)
	is.True(strings.Contains(string(html), `<article id="slug-456"><h1>Hello Mars!</h1></article>`))
}

func TestLayoutFrameNoProps(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	fsys := fstest.MapFS{
		"layout.svelte": &fstest.MapFile{Data: []byte(`
			<script>
				export let theme = 'light'
			</script>
			<html>
			<head>
				<title>My App</title>
			</head>
			<body>
				<main data-theme={theme}><slot /></main>
			</body>
			</html>
		`)},
		"frame.svelte": &fstest.MapFile{Data: []byte(`
			<script>
				export let id = '123'
			</script>
			<article id="slug-{id}"><slot /></article>
		`)},
		"index.svelte": &fstest.MapFile{Data: []byte(`
			<script>
				export let planet = 'Earth'
			</script>
			<h1>Hello {planet}!</h1>
		`)},
	}
	flag := &framework.Flag{}
	pages, err := viewer.Find(fsys)
	is.NoErr(err)
	svelte, err := loadViewer(dir, flag, fsys, pages)
	is.NoErr(err)
	viewers := viewer.Viewers{
		"index": svelte,
	}
	ctx := context.Background()
	html, err := viewers.Render(ctx, "index", map[string]interface{}{
		"index": map[string]interface{}{
			"planet": "Mars",
		},
		"frame": map[string]interface{}{
			"id": 456,
		},
	})
	is.NoErr(err)
	is.True(strings.Contains(string(html), `<main data-theme="light"><div id="bud_target"><article id="slug-456"><h1>Hello Mars!</h1></article></div></main>`))
}

func TestLayoutFrameProps(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	fsys := fstest.MapFS{
		"layout.svelte": &fstest.MapFile{Data: []byte(`
			<script>
				export let theme = 'light'
			</script>
			<html>
			<head>
				<title>My App</title>
			</head>
			<body>
				<main data-theme={theme}><slot /></main>
			</body>
			</html>
		`)},
		"frame.svelte": &fstest.MapFile{Data: []byte(`
			<script>
				export let id = '123'
			</script>
			<article id="slug-{id}"><slot /></article>
		`)},
		"index.svelte": &fstest.MapFile{Data: []byte(`
			<script>
				export let planet = 'Earth'
			</script>
			<h1>Hello {planet}!</h1>
		`)},
	}
	flag := &framework.Flag{}
	pages, err := viewer.Find(fsys)
	is.NoErr(err)
	svelte, err := loadViewer(dir, flag, fsys, pages)
	is.NoErr(err)
	viewers := viewer.Viewers{
		"index": svelte,
	}
	ctx := context.Background()
	html, err := viewers.Render(ctx, "index", map[string]interface{}{
		"index": map[string]interface{}{
			"planet": "Mars",
		},
		"frame": map[string]interface{}{
			"id": 456,
		},
		"layout": map[string]interface{}{
			"theme": "dark",
		},
	})
	is.NoErr(err)
	is.True(strings.Contains(string(html), `<main data-theme="dark"><div id="bud_target"><article id="slug-456"><h1>Hello Mars!</h1></article></div></main>`))
}

func TestMultipleFrames(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	fsys := fstest.MapFS{
		"layout.svelte": &fstest.MapFile{Data: []byte(`
			<script>
				export let theme = 'light'
			</script>
			<html>
			<head>
				<title>My App</title>
			</head>
			<body>
				<main data-theme={theme}><slot /></main>
			</body>
			</html>
		`)},
		"frame.svelte": &fstest.MapFile{Data: []byte(`
			<script>
				export let title = "My blog"
			</script>
			<header><h1>{title}</h1></header><slot />
		`)},
		"posts/frame.svelte": &fstest.MapFile{Data: []byte(`
			<script>
				export let id = '123'
			</script>
			<article id="slug-{id}"><slot /></article>
		`)},
		"posts/index.svelte": &fstest.MapFile{Data: []byte(`
			<script>
				export let planet = 'Earth'
			</script>
			<h1>Hello {planet}!</h1>
		`)},
	}
	flag := &framework.Flag{}
	pages, err := viewer.Find(fsys)
	is.NoErr(err)
	svelte, err := loadViewer(dir, flag, fsys, pages)
	is.NoErr(err)
	viewers := viewer.Viewers{
		"posts/index": svelte,
	}
	ctx := context.Background()
	html, err := viewers.Render(ctx, "posts/index", map[string]interface{}{
		"posts/index": map[string]interface{}{
			"planet": "Mars",
		},
		"posts/frame": map[string]interface{}{
			"id": 456,
		},
		"frame": map[string]interface{}{
			"title": "Matt's Blog",
		},
		"layout": map[string]interface{}{
			"theme": "navy",
		},
	})
	is.NoErr(err)
	is.True(strings.Contains(string(html), `<main data-theme="navy"><div id="bud_target"><article id="slug-456"><header><h1>Matt's Blog</h1></header><h1>Hello Mars!</h1></article></div></main></body></html>`))
}

func TestBundle(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	fsys := fstest.MapFS{
		"layout.svelte": &fstest.MapFile{Data: []byte(`
			<script>
				export let theme = 'light'
			</script>
			<html>
			<head>
				<title>My App</title>
			</head>
			<body>
				<main data-theme={theme}><slot /></main>
			</body>
			</html>
		`)},
		"frame.svelte": &fstest.MapFile{Data: []byte(`
			<script>
				export let title = "My blog"
			</script>
			<header><h1>{title}</h1></header><slot />
		`)},
		"posts/frame.svelte": &fstest.MapFile{Data: []byte(`
			<script>
				export let id = '123'
			</script>
			<article id="slug-{id}"><slot /></article>
		`)},
		"posts/index.svelte": &fstest.MapFile{Data: []byte(`
			<script>
				export let planet = 'Earth'
			</script>
			<h1>Hello {planet}!</h1>
		`)},
		"index.svelte": &fstest.MapFile{Data: []byte(`
			<script>
				export let planet = 'Mercury'
			</script>
			<h1>Hello {planet}!</h1>
		`)},
	}
	flag := &framework.Flag{}
	pages, err := viewer.Find(fsys)
	is.NoErr(err)
	svelte, err := loadViewer(dir, flag, fsys, pages)
	is.NoErr(err)
	viewers := viewer.Viewers{
		"index":       svelte,
		"posts/index": svelte,
	}
	ctx := context.Background()
	outfs := virtual.Map{}
	err = viewers.Bundle(ctx, outfs)
	is.NoErr(err)
	is.True(outfs[".ssr.js"] != nil)
	flag = &framework.Flag{Embed: true, Minify: true}
	svelte, err = loadViewer(dir, flag, outfs, pages)
	is.NoErr(err)
	viewers = viewer.Viewers{
		"index":       svelte,
		"posts/index": svelte,
	}
	// Try with index
	html, err := viewers.Render2(ctx, "index", map[string]interface{}{
		"index": map[string]interface{}{
			"planet": "Venus",
		},
		"frame": map[string]interface{}{
			"title": "Mark's Blog",
		},
		"layout": map[string]interface{}{
			"theme": "mango",
		},
	})
	is.NoErr(err)
	is.True(strings.Contains(string(html), `<main data-theme="mango"><div id="bud_target"><header><h1>Mark's Blog</h1></header><h1>Hello Venus!</h1></div></main></body></html>`))

	// Try with posts/index
	html, err = viewers.Render2(ctx, "posts/index", map[string]interface{}{
		"posts/index": map[string]interface{}{
			"planet": "Mars",
		},
		"posts/frame": map[string]interface{}{
			"id": 456,
		},
		"frame": map[string]interface{}{
			"title": "Matt's Blog",
		},
		"layout": map[string]interface{}{
			"theme": "navy",
		},
	})
	is.NoErr(err)
	is.True(strings.Contains(string(html), `<main data-theme="navy"><div id="bud_target"><article id="slug-456"><header><h1>Matt's Blog</h1></header><h1>Hello Mars!</h1></article></div></main></body></html>`))
}
