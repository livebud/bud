package svelte_test

import (
	"context"
	"io/fs"
	"strings"
	"testing"

	"github.com/livebud/bud/internal/testdir"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/esb"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/js"
	v8 "github.com/livebud/bud/package/js/v8"
	"github.com/livebud/bud/package/log/testlog"
	"github.com/livebud/bud/runtime/transpiler"
	"github.com/livebud/bud/runtime/view"
	"github.com/livebud/bud/runtime/view/svelte"
)

// func loadViewer(dir string, flag *framework.Flag, fsys fs.FS, module *gomod.Module) (*svelte.Viewer, error) {
// 	esbuilder := esb.New(module)
// 	js, err := v8.Load()
// 	if err != nil {
// 		return nil, err
// 	}
// 	// js, err := v8.Load(&js.Console{
// 	// 	Log:   os.Stdout,
// 	// 	Error: os.Stderr,
// 	// })
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
// 	tr := transpiler.NewClient(fsys)
// 	viewer := svelte.New(nil, esbuilder, flag, js, module, tr)
// 	return viewer, nil
// }

func testTranspiler(ctx context.Context, flag *framework.Flag, fsys fs.FS, js js.VM) (transpiler.Transpiler, error) {
	compiler, err := svelte.Load(flag, js)
	if err != nil {
		return nil, err
	}
	tr := transpiler.NewTester(fsys)
	tr.Add(".svelte", ".ssr.js", func(file *transpiler.File) error {
		ssr, err := compiler.SSR(ctx, file.Path(), file.Data)
		if err != nil {
			return err
		}
		file.Data = []byte(ssr.JS)
		return nil
	})
	return tr, nil
}

func TestEntryNoProps(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["view/index.svelte"] = `
		<script>
			export let planet = 'Earth'
		</script>
		<h1>Hello {planet}!</h1>
	`
	is.NoErr(td.Write(ctx))
	flag := &framework.Flag{}
	module, err := gomod.Find(dir)
	is.NoErr(err)
	pages, err := view.Find(module)
	is.NoErr(err)
	esbuilder := esb.New(module, log)
	js, err := v8.Load()
	is.NoErr(err)
	tr, err := testTranspiler(ctx, flag, module, js)
	is.NoErr(err)
	viewer := svelte.New(nil, esbuilder, flag, js, module, tr)
	html, err := viewer.Render(ctx, pages["view/index"], nil)
	is.NoErr(err)
	is.True(strings.Contains(string(html), `<h1>Hello Earth!</h1>`))
	// svelte, err := loadViewer(dir, flag, module, module)
	// is.NoErr(err)
	// fmt.Println(pages)
	// fmt.Println(svelte)
	// viewers := view.Viewers{
	// 	"index": svelte,
	// }
	// ctx := context.Background()
	// html, err := viewers.Render(ctx, "index", nil)
	// is.NoErr(err)
	// is.True(strings.Contains(string(html), `<h1>Hello Earth!</h1>`))
}

// func TestEntryProps(t *testing.T) {
// 	is := is.New(t)
// 	dir := t.TempDir()
// 	fsys := fstest.MapFS{
// 		"index.svelte": &fstest.MapFile{Data: []byte(`
// 			<script>
// 				export let planet = 'Earth'
// 			</script>
// 			<h1>Hello {planet}!</h1>
// 		`)},
// 	}
// 	flag := &framework.Flag{}
// 	pages, err := view.Find(fsys)
// 	is.NoErr(err)
// 	svelte, err := loadViewer(dir, flag, fsys, pages)
// 	is.NoErr(err)
// 	viewers := view.Viewers{
// 		"index": svelte,
// 	}
// 	ctx := context.Background()
// 	html, err := viewers.Render(ctx, "index", map[string]interface{}{
// 		"index": map[string]interface{}{
// 			"planet": "Mars",
// 		},
// 	})
// 	is.NoErr(err)
// 	is.True(strings.Contains(string(html), `<h1>Hello Mars!</h1>`))
// }

// func TestFrameProps(t *testing.T) {
// 	is := is.New(t)
// 	dir := t.TempDir()
// 	fsys := fstest.MapFS{
// 		"frame.svelte": &fstest.MapFile{Data: []byte(`
// 			<script>
// 				export let id = '123'
// 			</script>
// 			<article id="slug-{id}"><slot /></article>
// 		`)},
// 		"index.svelte": &fstest.MapFile{Data: []byte(`
// 			<script>
// 				export let planet = 'Earth'
// 			</script>
// 			<h1>Hello {planet}!</h1>
// 		`)},
// 	}
// 	flag := &framework.Flag{}
// 	pages, err := view.Find(fsys)
// 	is.NoErr(err)
// 	svelte, err := loadViewer(dir, flag, fsys, pages)
// 	is.NoErr(err)
// 	viewers := view.Viewers{
// 		"index": svelte,
// 	}
// 	ctx := context.Background()
// 	html, err := viewers.Render(ctx, "index", map[string]interface{}{
// 		"index": map[string]interface{}{
// 			"planet": "Mars",
// 		},
// 		"frame": map[string]interface{}{
// 			"id": 456,
// 		},
// 	})
// 	is.NoErr(err)
// 	is.True(strings.Contains(string(html), `<article id="slug-456"><h1>Hello Mars!</h1></article>`))
// }

// func TestLayoutFrameNoProps(t *testing.T) {
// 	is := is.New(t)
// 	dir := t.TempDir()
// 	fsys := fstest.MapFS{
// 		"layout.svelte": &fstest.MapFile{Data: []byte(`
// 			<script>
// 				export let theme = 'light'
// 			</script>
// 			<html>
// 			<head>
// 				<title>My App</title>
// 			</head>
// 			<body>
// 				<main data-theme={theme}><slot /></main>
// 			</body>
// 			</html>
// 		`)},
// 		"frame.svelte": &fstest.MapFile{Data: []byte(`
// 			<script>
// 				export let id = '123'
// 			</script>
// 			<article id="slug-{id}"><slot /></article>
// 		`)},
// 		"index.svelte": &fstest.MapFile{Data: []byte(`
// 			<script>
// 				export let planet = 'Earth'
// 			</script>
// 			<h1>Hello {planet}!</h1>
// 		`)},
// 	}
// 	flag := &framework.Flag{}
// 	pages, err := view.Find(fsys)
// 	is.NoErr(err)
// 	svelte, err := loadViewer(dir, flag, fsys, pages)
// 	is.NoErr(err)
// 	viewers := view.Viewers{
// 		"index": svelte,
// 	}
// 	ctx := context.Background()
// 	html, err := viewers.Render(ctx, "index", map[string]interface{}{
// 		"index": map[string]interface{}{
// 			"planet": "Mars",
// 		},
// 		"frame": map[string]interface{}{
// 			"id": 456,
// 		},
// 	})
// 	is.NoErr(err)
// 	is.True(strings.Contains(string(html), `<main data-theme="light"><div id="bud_target"><article id="slug-456"><h1>Hello Mars!</h1></article></div></main>`))
// }

// func TestLayoutFrameProps(t *testing.T) {
// 	is := is.New(t)
// 	dir := t.TempDir()
// 	fsys := fstest.MapFS{
// 		"layout.svelte": &fstest.MapFile{Data: []byte(`
// 			<script>
// 				export let theme = 'light'
// 			</script>
// 			<html>
// 			<head>
// 				<title>My App</title>
// 			</head>
// 			<body>
// 				<main data-theme={theme}><slot /></main>
// 			</body>
// 			</html>
// 		`)},
// 		"frame.svelte": &fstest.MapFile{Data: []byte(`
// 			<script>
// 				export let id = '123'
// 			</script>
// 			<article id="slug-{id}"><slot /></article>
// 		`)},
// 		"index.svelte": &fstest.MapFile{Data: []byte(`
// 			<script>
// 				export let planet = 'Earth'
// 			</script>
// 			<h1>Hello {planet}!</h1>
// 		`)},
// 	}
// 	flag := &framework.Flag{}
// 	pages, err := view.Find(fsys)
// 	is.NoErr(err)
// 	svelte, err := loadViewer(dir, flag, fsys, pages)
// 	is.NoErr(err)
// 	viewers := view.Viewers{
// 		"index": svelte,
// 	}
// 	ctx := context.Background()
// 	html, err := viewers.Render(ctx, "index", map[string]interface{}{
// 		"index": map[string]interface{}{
// 			"planet": "Mars",
// 		},
// 		"frame": map[string]interface{}{
// 			"id": 456,
// 		},
// 		"layout": map[string]interface{}{
// 			"theme": "dark",
// 		},
// 	})
// 	is.NoErr(err)
// 	is.True(strings.Contains(string(html), `<main data-theme="dark"><div id="bud_target"><article id="slug-456"><h1>Hello Mars!</h1></article></div></main>`))
// }

// func TestMultipleFrames(t *testing.T) {
// 	is := is.New(t)
// 	dir := t.TempDir()
// 	fsys := fstest.MapFS{
// 		"layout.svelte": &fstest.MapFile{Data: []byte(`
// 			<script>
// 				export let theme = 'light'
// 			</script>
// 			<html>
// 			<head>
// 				<title>My App</title>
// 			</head>
// 			<body>
// 				<main data-theme={theme}><slot /></main>
// 			</body>
// 			</html>
// 		`)},
// 		"frame.svelte": &fstest.MapFile{Data: []byte(`
// 			<script>
// 				export let title = "My blog"
// 			</script>
// 			<header><h1>{title}</h1></header><slot />
// 		`)},
// 		"posts/frame.svelte": &fstest.MapFile{Data: []byte(`
// 			<script>
// 				export let id = '123'
// 			</script>
// 			<article id="slug-{id}"><slot /></article>
// 		`)},
// 		"posts/index.svelte": &fstest.MapFile{Data: []byte(`
// 			<script>
// 				export let planet = 'Earth'
// 			</script>
// 			<h1>Hello {planet}!</h1>
// 		`)},
// 	}
// 	flag := &framework.Flag{}
// 	pages, err := view.Find(fsys)
// 	is.NoErr(err)
// 	svelte, err := loadViewer(dir, flag, fsys, pages)
// 	is.NoErr(err)
// 	viewers := view.Viewers{
// 		"posts/index": svelte,
// 	}
// 	ctx := context.Background()
// 	html, err := viewers.Render(ctx, "posts/index", map[string]interface{}{
// 		"posts/index": map[string]interface{}{
// 			"planet": "Mars",
// 		},
// 		"posts/frame": map[string]interface{}{
// 			"id": 456,
// 		},
// 		"frame": map[string]interface{}{
// 			"title": "Matt's Blog",
// 		},
// 		"layout": map[string]interface{}{
// 			"theme": "navy",
// 		},
// 	})
// 	is.NoErr(err)
// 	is.True(strings.Contains(string(html), `<main data-theme="navy"><div id="bud_target"><article id="slug-456"><header><h1>Matt's Blog</h1></header><h1>Hello Mars!</h1></article></div></main></body></html>`))
// }

// func TestBundle(t *testing.T) {
// 	is := is.New(t)
// 	dir := t.TempDir()
// 	fsys := fstest.MapFS{
// 		"layout.svelte": &fstest.MapFile{Data: []byte(`
// 			<script>
// 				export let theme = 'light'
// 			</script>
// 			<html>
// 			<head>
// 				<title>My App</title>
// 			</head>
// 			<body>
// 				<main data-theme={theme}><slot /></main>
// 			</body>
// 			</html>
// 		`)},
// 		"frame.svelte": &fstest.MapFile{Data: []byte(`
// 			<script>
// 				export let title = "My blog"
// 			</script>
// 			<header><h1>{title}</h1></header><slot />
// 		`)},
// 		"posts/frame.svelte": &fstest.MapFile{Data: []byte(`
// 			<script>
// 				export let id = '123'
// 			</script>
// 			<article id="slug-{id}"><slot /></article>
// 		`)},
// 		"posts/index.svelte": &fstest.MapFile{Data: []byte(`
// 			<script>
// 				export let planet = 'Earth'
// 			</script>
// 			<h1>Hello {planet}!</h1>
// 		`)},
// 		"index.svelte": &fstest.MapFile{Data: []byte(`
// 			<script>
// 				export let planet = 'Mercury'
// 			</script>
// 			<h1>Hello {planet}!</h1>
// 		`)},
// 	}
// 	flag := &framework.Flag{}
// 	pages, err := view.Find(fsys)
// 	is.NoErr(err)
// 	svelte, err := loadViewer(dir, flag, fsys, pages)
// 	is.NoErr(err)
// 	viewers := view.Viewers{
// 		"index":       svelte,
// 		"posts/index": svelte,
// 	}
// 	ctx := context.Background()
// 	outfs := virtual.Map{}
// 	err = viewers.Bundle(ctx, outfs)
// 	is.NoErr(err)
// 	is.True(outfs[".ssr.js"] != nil)
// 	flag = &framework.Flag{Embed: true, Minify: true}
// 	svelte, err = loadViewer(dir, flag, outfs, pages)
// 	is.NoErr(err)
// 	viewers = view.Viewers{
// 		"index":       svelte,
// 		"posts/index": svelte,
// 	}
// 	// Try with index
// 	html, err := viewers.Render2(ctx, "index", map[string]interface{}{
// 		"index": map[string]interface{}{
// 			"planet": "Venus",
// 		},
// 		"frame": map[string]interface{}{
// 			"title": "Mark's Blog",
// 		},
// 		"layout": map[string]interface{}{
// 			"theme": "mango",
// 		},
// 	})
// 	is.NoErr(err)
// 	is.True(strings.Contains(string(html), `<main data-theme="mango"><div id="bud_target"><header><h1>Mark's Blog</h1></header><h1>Hello Venus!</h1></div></main></body></html>`))

// 	// Try with posts/index
// 	html, err = viewers.Render2(ctx, "posts/index", map[string]interface{}{
// 		"posts/index": map[string]interface{}{
// 			"planet": "Mars",
// 		},
// 		"posts/frame": map[string]interface{}{
// 			"id": 456,
// 		},
// 		"frame": map[string]interface{}{
// 			"title": "Matt's Blog",
// 		},
// 		"layout": map[string]interface{}{
// 			"theme": "navy",
// 		},
// 	})
// 	is.NoErr(err)
// 	is.True(strings.Contains(string(html), `<main data-theme="navy"><div id="bud_target"><article id="slug-456"><header><h1>Matt's Blog</h1></header><h1>Hello Mars!</h1></article></div></main></body></html>`))
// }
