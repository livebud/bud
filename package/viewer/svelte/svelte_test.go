package svelte_test

import (
	"context"
	"io/fs"
	"testing"

	"github.com/livebud/bud/framework/transform2/transformrt"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log/testlog"
	"github.com/livebud/bud/package/viewer"
	"github.com/livebud/bud/package/virtual"

	"github.com/livebud/bud/internal/is"

	v8 "github.com/livebud/bud/package/js/v8"
	svelteCompiler "github.com/livebud/bud/package/svelte"
	"github.com/livebud/bud/package/viewer/svelte"
)

func loadViewer(t testing.TB, fsys fs.FS) *svelte.Viewer {
	is := is.New(t)
	is.Helper()
	module := gomod.New(t.TempDir())
	log := testlog.New()
	vm, err := v8.Load()
	is.NoErr(err)
	compiler, err := svelteCompiler.Load(vm)
	is.NoErr(err)
	transformer := transformrt.Load(log,
		&transformrt.Transform{
			From: ".svelte",
			To:   ".ssr.js",
			Func: func(file *transformrt.File) error {
				ssr, err := compiler.SSR(file.Path(), file.Data)
				if err != nil {
					return err
				}
				file.Data = []byte(ssr.JS)
				return nil
			},
		},
	)
	return svelte.New(fsys, log, module, transformer, vm)
}

func TestOnlyPage(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	view := loadViewer(t, virtual.Map{
		"view/posts/show.svelte": &virtual.File{
			Data: []byte(`<h1>Posts</h1>`),
		},
	})
	html, err := view.Render(ctx, &viewer.Page{
		Main: &viewer.View{
			Path:  "view/posts/show.svelte",
			Props: viewer.Props{},
		},
	})
	is.NoErr(err)
	is.In(string(html), `<h1>Posts</h1>`)
	// default layout
	is.In(string(html), `<!doctype html>`)
	is.In(string(html), `<meta charset="utf-8" />`)
}

func TestSimple(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	view := loadViewer(t, virtual.Map{
		"view/layout.svelte": &virtual.File{
			Data: []byte(`
				<script>
					export let title = "My Wonderful Blog"
					export let theme = "light"
				</script>
				<html lang="en">
					<head>
						<slot name="head" />
						<slot name="style" />
					</head>
					<body data-theme={theme}>
						<h1>{title}</h1>
						<main>
							<slot />
						</main>
					</body>
				</html>
			`),
		},
		"view/frame.svelte": &virtual.File{
			Data: []byte(`
				<script>
					export let categories = []
				</script>
				<main>
					<aside>
						{#each categories as category}
							<p>{category}</p>
						{/each}
					</aside>
					<article><slot /></article>
				</main>
				<style>
					aside p {
						padding: 10px;
					}
				</style>
			`),
		},
		"view/posts/frame.svelte": &virtual.File{
			Data: []byte(`
				<script>
					export let likes = 0
				</script>
				<div>
					<slot />
					<hr />
					{likes} people liked this.
				</div>
			`),
		},
		"view/posts/show.svelte": &virtual.File{
			Data: []byte(`
				<script>
					export let title = "Untitled"
					export let body = "No body"
				</script>
				<svelte:head>
					<title>{title}</title>
				</svelte:head>
				<h1>{title}</h1>
				<p>{body}</p>
				<style>
					h1 {
						color: red;
					}
				</style>
			`),
		},
	})
	html, err := view.Render(ctx, &viewer.Page{
		Layout: &viewer.View{
			Path: "view/layout.svelte",
			Props: map[string]interface{}{
				"title": "My Wonderful Blog",
				"theme": "dark",
			},
		},
		Frames: []*viewer.View{
			{
				Path: "view/frame.svelte",
				Props: map[string]interface{}{
					"categories": []string{"Science", "Technology", "Engineering", "Math"},
				},
			},
			{
				Path: "view/posts/frame.svelte",
				Props: map[string]interface{}{
					"likes": 42,
				},
			},
		},
		Main: &viewer.View{
			Path: "view/posts/show.svelte",
			Props: map[string]interface{}{
				"title": "Hello World",
				"body":  "This is my first post!",
			},
		},
	})
	is.NoErr(err)
	is.In(string(html), `data-theme="dark"`)
	is.In(string(html), `Science`)
	is.In(string(html), `Technology`)
	is.In(string(html), `42 people liked this.`)
	is.In(string(html), `<p>This is my first post!</p>`)
}

func TestError(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	view := loadViewer(t, virtual.Map{
		"view/posts/show.svelte": &virtual.File{
			Data: []byte(`<h1>{title}</h1>`),
		},
		"view/error.svelte": &virtual.File{
			Data: []byte(`
				<script>
					export let message = "Something went wrong"
				</script>
				<h1>An error occurred</h1>
				<pre>{message}</pre>
			`),
		},
	})
	html, err := view.Render(ctx, &viewer.Page{
		Main: &viewer.View{
			Path:  "view/posts/show.svelte",
			Props: viewer.Props{},
		},
	})
	is.True(err != nil)
	is.Equal(err.Error(), `ReferenceError: title is not defined`)
	is.Equal(html, nil)
	html = view.RenderError(ctx, &viewer.Page{
		Main: &viewer.View{
			Path: "view/error.svelte",
			Props: viewer.Props{
				"message": err.Error(),
			},
		},
	})
	is.In(string(html), `ReferenceError: title is not defined`)
	// default layout
	is.In(string(html), `<!doctype html>`)
	is.In(string(html), `<meta charset="utf-8" />`)
}
