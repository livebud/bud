package gohtml_test

import (
	"testing"
)

func TestPage(t *testing.T) {
	// is := is.New(t)
	// log := testlog.New()
	// fsys := virtual.Map{
	// 	"index.gohtml": "Hello {{ .Planet }}!",
	// }
	// renderer := gohtml.New(fsys, log, transpiler.New())
	// ctx := context.Background()
	// html := new(bytes.Buffer)
	// err := renderer.Render(ctx, html, "index", map[string]interface{}{
	// 	"Planet": "Earth",
	// })
	// is.NoErr(err)
	// is.Equal(html.String(), "Hello Earth!")
}

// func TestLayout(t *testing.T) {
// 	is := is.New(t)
// 	log := testlog.New()
// 	fsys := virtual.Map{
// 		"index.gohtml":  "Hello {{ .Planet }}!",
// 		"layout.gohtml": "<html>{{ . }}</html>",
// 	}
// 	pages, err := viewer.Find(fsys)
// 	is.NoErr(err)
// 	viewer := gohtml.New(fsys, log, pages, transpiler.New())
// 	ctx := context.Background()
// 	html, err := viewer.Render(ctx, "index", map[string]interface{}{
// 		"index": map[string]interface{}{
// 			"Planet": "Earth",
// 		},
// 	})
// 	is.NoErr(err)
// 	is.Equal(string(html), "<html>Hello Earth!</html>")
// }

// func TestRenderError(t *testing.T) {
// 	is := is.New(t)
// 	log := testlog.New()
// 	fsys := virtual.Map{
// 		"index.gohtml":  "Hello {{ .Planet }}!",
// 		"frame.gohtml":  `<main>{{ . }}</main>`,
// 		"layout.gohtml": "<html>{{ . }}</html>",
// 		"error.gohtml":  `<div class="error">{{ .Message }}</div>`,
// 	}
// 	pages, err := viewer.Find(fsys)
// 	is.NoErr(err)
// 	gohtml := gohtml.New(fsys, log, pages, transpiler.New())
// 	ctx := context.Background()
// 	html := gohtml.RenderError(ctx, "index", map[string]interface{}{
// 		"index": map[string]interface{}{"Planet": "Earth"},
// 	}, errors.New("some error"))
// 	is.Equal(string(html), `<html><main><div class="error">some error</div></main></html>`)
// }

// // TODO: this should have a default page
// func TestRenderErrorNoPage(t *testing.T) {
// 	is := is.New(t)
// 	log := testlog.New()
// 	fsys := virtual.Map{
// 		"index.gohtml":  "Hello {{ .Planet }}!",
// 		"frame.gohtml":  `<main>{{ . }}</main>`,
// 		"layout.gohtml": "<html>{{ . }}</html>",
// 	}
// 	pages, err := viewer.Find(fsys)
// 	is.NoErr(err)
// 	gohtml := gohtml.New(fsys, log, pages, transpiler.New())
// 	ctx := context.Background()
// 	html := gohtml.RenderError(ctx, "index", map[string]interface{}{
// 		"index": map[string]interface{}{"Planet": "Earth"},
// 	}, errors.New("some error"))
// 	is.Equal(string(html), `gohtml: page "index" has no error page to render error. some error`)
// }

// func TestRenderErrorWithFrames(t *testing.T) {
// 	is := is.New(t)
// 	log := testlog.New()
// 	fsys := virtual.Map{
// 		"posts/index.gohtml": "Hello {{ .Planet }}!",
// 		"posts/frame.gohtml": `<div class="posts">{{ . }}</div>`,
// 		"frame.gohtml":       `<main>{{ . }}</main>`,
// 		"layout.gohtml":      "<html>{{ . }}</html>",
// 		"error.gohtml":       `<div class="error">{{ .Message }}</div>`,
// 	}
// 	pages, err := viewer.Find(fsys)
// 	is.NoErr(err)
// 	gohtml := gohtml.New(fsys, log, pages, transpiler.New())
// 	ctx := context.Background()
// 	html := gohtml.RenderError(ctx, "posts/index", map[string]interface{}{
// 		"posts/index": map[string]interface{}{"Planet": "Earth"},
// 	}, errors.New("some error"))
// 	is.Equal(string(html), `<html><main><div class="error">some error</div></main></html>`)
// }

// func TestBundle(t *testing.T) {
// 	is := is.New(t)
// 	log := testlog.New()
// 	fsys := virtual.Map{
// 		"index.gohtml":  "Hello {{ .Planet }}!",
// 		"layout.gohtml": "<html>{{ . }}</html>",
// 	}
// 	pages, err := viewer.Find(fsys)
// 	is.NoErr(err)
// 	viewer := gohtml.New(fsys, log, pages, transpiler.New())
// 	out := virtual.Tree{}
// 	ctx := context.Background()
// 	err = viewer.Bundle(ctx, out)
// 	is.NoErr(err)
// 	viewer = gohtml.New(out, log, pages, transpiler.New())
// 	html, err := viewer.Render(ctx, "index", map[string]interface{}{
// 		"index": map[string]interface{}{
// 			"Planet": "Earth",
// 		},
// 	})
// 	is.NoErr(err)
// 	is.Equal(string(html), "<html>Hello Earth!</html>")
// }
