package svelte_test

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sync"
	"testing"

	"github.com/livebud/bud/runtime/transpiler"

	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log/testlog"
	"github.com/livebud/bud/package/router"
	"github.com/livebud/bud/package/viewer"
	"github.com/livebud/bud/package/virtual"

	"github.com/livebud/bud/internal/is"

	v8 "github.com/livebud/bud/package/js/v8"
	"github.com/livebud/bud/package/viewer/svelte"
)

func newView(svelte *svelte.Viewer) *View {
	return &View{
		Svelte: svelte,
		router: router.New(),
	}
}

type View struct {
	Svelte *svelte.Viewer

	registerOnce sync.Once
	router       *router.Router
}

func (v *View) Register(r *router.Router) error {
	staticPage := &viewer.Page{
		View: &viewer.View{
			Key:     "posts/show",
			Path:    "view/posts/show.svelte",
			Props:   viewer.Props{},
			Context: viewer.Context{},
		},
		Client: "/bud/view/posts/show.svelte.js",
		Frames: []*viewer.View{
			{
				Key:     "posts/frame",
				Path:    "view/posts/frame.svelte",
				Props:   viewer.Props{},
				Context: viewer.Context{},
			},
			{
				Key:     "frame",
				Path:    "view/frame.svelte",
				Props:   viewer.Props{},
				Context: viewer.Context{},
			},
		},
	}
	// Register the page if it doesn't exist already
	r.Get("/posts/:id", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		html, err := v.Render(r.Context(), staticPage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(html)
	}))
	v.Svelte.RegisterClient(r, staticPage)
	return nil
}

var keyPaths = map[string]string{
	"frame":       "view/frame.svelte",
	"posts/show":  "view/posts/show.svelte",
	"posts/frame": "view/posts/frame.svelte",
}

func (v *View) Viewer(path string) (viewer viewer.Viewer, err error) {
	switch filepath.Ext(path) {
	case ".svelte":
		return v.Svelte, nil
	default:
		return nil, fmt.Errorf("view: no viewer for %q", path)
	}
}

func (v *View) Render(ctx context.Context, page *viewer.Page) ([]byte, error) {
	// Attach the paths to the views
	if path, ok := keyPaths[page.Key]; ok {
		page.Path = path
	}
	for _, frame := range page.Frames {
		if path, ok := keyPaths[frame.Key]; ok {
			frame.Path = path
		}
	}
	if page.Layout != nil {
		if path, ok := keyPaths[page.Layout.Key]; ok {
			page.Layout.Path = path
		}
	}
	viewer, err := v.Viewer(page.Path)
	if err != nil {
		return nil, err
	}
	return viewer.Render(ctx, page)
}

// ServeHTTP allows View to be used standalone.
func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.registerOnce.Do(func() { v.Register(v.router) })
	v.router.ServeHTTP(w, r)
}

func loadViewer(t testing.TB, fsys fs.FS) *svelte.Viewer {
	is := is.New(t)
	is.Helper()
	module := gomod.New(t.TempDir())
	log := testlog.New()
	vm, err := v8.Load()
	is.NoErr(err)
	// compiler, err := svelteCompiler.Load(vm)
	// is.NoErr(err)
	tr := transpiler.New()
	// transformer := transformrt.Load(log,
	// 	&transformrt.Transform{
	// 		From: ".svelte",
	// 		To:   ".ssr.js",
	// 		Func: func(file *transformrt.File) error {
	// 			ssr, err := compiler.SSR(file.Path(), file.Data)
	// 			if err != nil {
	// 				return err
	// 			}
	// 			file.Data = []byte(ssr.JS)
	// 			return nil
	// 		},
	// 	},
	// 	&transformrt.Transform{
	// 		From: ".svelte",
	// 		To:   ".js",
	// 		Func: func(file *transformrt.File) error {
	// 			dom, err := compiler.DOM(file.Path(), file.Data)
	// 			if err != nil {
	// 				return err
	// 			}
	// 			file.Data = []byte(dom.JS)
	// 			return nil
	// 		},
	// 	},
	// )
	return svelte.New(fsys, log, module, tr, vm)
}

func TestServeViewSSR(t *testing.T) {
	is := is.New(t)
	viewer := loadViewer(t, virtual.Map{
		"view/posts/show.svelte": &virtual.File{
			Data: []byte(`<h1>Posts</h1>`),
		},
		"view/posts/frame.svelte": &virtual.File{
			Data: []byte(`<div><slot /></div>`),
		},
		"view/frame.svelte": &virtual.File{
			Data: []byte(`<article><slot /></article>`),
		},
	})
	view := newView(viewer)
	r := httptest.NewRequest("GET", "/posts/10", nil)
	w := httptest.NewRecorder()
	view.ServeHTTP(w, r)
	res := w.Result()
	body, err := io.ReadAll(res.Body)
	is.NoErr(err)
	is.In(string(body), `<div id="bud_target"><article><div><h1>Posts</h1></div></article></div>`)
	// default layout
	is.In(string(body), `<!doctype html>`)
	is.In(string(body), `<meta charset="utf-8" />`)
	is.In(string(body), `<script src="/bud/view/posts/show.svelte.js" defer></script>`)
	is.In(string(body), `<script id="bud_data" type="text/data" defer>{`)
	is.Equal(res.StatusCode, 200)
	is.Equal(res.Header.Get("Content-Type"), "text/html; charset=utf-8")
}

func TestServeViewDOM(t *testing.T) {
	is := is.New(t)
	viewer := loadViewer(t, virtual.Map{
		"view/posts/show.svelte": &virtual.File{
			Data: []byte(`<b>Posts</b>`),
		},
		"view/posts/frame.svelte": &virtual.File{
			Data: []byte(`<u><slot /></u>`),
		},
		"view/frame.svelte": &virtual.File{
			Data: []byte(`<em><slot /></em>`),
		},
	})
	view := newView(viewer)
	// Get the client
	r := httptest.NewRequest("GET", "/bud/view/posts/show.svelte.js", nil)
	w := httptest.NewRecorder()
	view.ServeHTTP(w, r)
	res := w.Result()
	body, err := io.ReadAll(res.Body)
	is.NoErr(err)
	is.In(string(body), `"em"`)
	is.In(string(body), `"u"`)
	is.In(string(body), `"b"`)
	is.In(string(body), `hydrator.hydrate(document.getElementById("bud_target"), document.getElementById("bud_data"))`)
	is.Equal(res.StatusCode, 200)
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
		View: &viewer.View{
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
		View: &viewer.View{
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
		View: &viewer.View{
			Path:  "view/posts/show.svelte",
			Props: viewer.Props{},
		},
	})
	is.True(err != nil)
	is.Equal(err.Error(), `ReferenceError: title is not defined`)
	is.Equal(html, nil)
	html = view.RenderError(ctx, &viewer.Page{
		View: &viewer.View{
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
