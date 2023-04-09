package gohtml_test

import (
	"context"
	"errors"
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/log/testlog"
	"github.com/livebud/bud/package/viewer"
	"github.com/livebud/bud/package/viewer/gohtml"
	"github.com/livebud/bud/package/virtual"
	"github.com/livebud/bud/runtime/transpiler"
)

func TestPage(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	fsys := virtual.Map{
		"index.gohtml": "Hello {{ .Planet }}!",
	}
	pages, err := viewer.Find(fsys)
	is.NoErr(err)
	viewer := gohtml.New(fsys, log, pages, transpiler.New())
	ctx := context.Background()
	html, err := viewer.Render(ctx, "index", map[string]interface{}{
		"index": map[string]interface{}{
			"Planet": "Earth",
		},
	})
	is.NoErr(err)
	is.Equal(string(html), "Hello Earth!")
}

func TestLayout(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	fsys := virtual.Map{
		"index.gohtml":  "Hello {{ .Planet }}!",
		"layout.gohtml": "<html>{{ . }}</html>",
	}
	pages, err := viewer.Find(fsys)
	is.NoErr(err)
	viewer := gohtml.New(fsys, log, pages, transpiler.New())
	ctx := context.Background()
	html, err := viewer.Render(ctx, "index", map[string]interface{}{
		"index": map[string]interface{}{
			"Planet": "Earth",
		},
	})
	is.NoErr(err)
	is.Equal(string(html), "<html>Hello Earth!</html>")
}

func TestRenderError(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	fsys := virtual.Map{
		"index.gohtml":  "Hello {{ .Planet }}!",
		"layout.gohtml": "<html>{{ . }}</html>",
		"error.gohtml":  `<div class="error">{{ .Message }}</div>`,
	}
	pages, err := viewer.Find(fsys)
	is.NoErr(err)
	viewer := gohtml.New(fsys, log, pages, transpiler.New())
	ctx := context.Background()
	html := viewer.RenderError(ctx, "index", map[string]interface{}{
		"index": map[string]interface{}{
			"Planet": "Earth",
		},
	}, errors.New("some error"))
	is.Equal(string(html), `<html><div class="error">some error</div></html>`)
}

func TestBundle(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	fsys := virtual.Map{
		"index.gohtml":  "Hello {{ .Planet }}!",
		"layout.gohtml": "<html>{{ . }}</html>",
	}
	pages, err := viewer.Find(fsys)
	is.NoErr(err)
	viewer := gohtml.New(fsys, log, pages, transpiler.New())
	out := virtual.Tree{}
	ctx := context.Background()
	err = viewer.Bundle(ctx, out)
	is.NoErr(err)
	viewer = gohtml.New(out, log, pages, transpiler.New())
	html, err := viewer.Render(ctx, "index", map[string]interface{}{
		"index": map[string]interface{}{
			"Planet": "Earth",
		},
	})
	is.NoErr(err)
	is.Equal(string(html), "<html>Hello Earth!</html>")
}
