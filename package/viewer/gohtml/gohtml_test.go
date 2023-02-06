package gohtml_test

import (
	"context"
	"testing"
	"testing/fstest"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/viewer"
	"github.com/livebud/bud/package/viewer/gohtml"
	"github.com/livebud/bud/runtime/transpiler"
)

func TestEntry(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{
		"index.gohtml": &fstest.MapFile{Data: []byte("Hello {{ .Planet }}!")},
	}
	// Find the pages
	pages, err := viewer.Find(fsys)
	is.NoErr(err)
	tr := transpiler.New()
	gohtml := gohtml.New(fsys, tr, pages)
	viewers := viewer.Viewers{
		"index": gohtml,
	}
	ctx := context.Background()
	html, err := viewers.Render(ctx, "index", map[string]interface{}{
		"index": map[string]interface{}{
			"Planet": "Earth",
		},
	})
	is.NoErr(err)
	is.Equal(string(html), "Hello Earth!")
}
