package viewer_test

import (
	"testing"
	"testing/fstest"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/viewer"
)

func TestIndex(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{
		"index.gohtml": &fstest.MapFile{Data: []byte("Hello {{ .Planet }}!")},
	}
	// Find the pages
	pages, err := viewer.Find(fsys)
	is.NoErr(err)
	is.Equal(len(pages), 1)
	is.True(pages["index"] != nil)
	is.Equal(pages["index"].Path, "index.gohtml")
	is.Equal(len(pages["index"].Frames), 0)
	is.Equal(pages["index"].Layout, nil)
	is.Equal(pages["index"].Error, nil)
}

func TestNested(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{
		"layout.svelte":      &fstest.MapFile{Data: []byte(`<slot />`)},
		"frame.svelte":       &fstest.MapFile{Data: []byte(`<slot />`)},
		"posts/frame.svelte": &fstest.MapFile{Data: []byte(`<slot />`)},
		"posts/index.svelte": &fstest.MapFile{Data: []byte(`<h1>Hello {planet}!</h1>`)},
	}
	// Find the pages
	pages, err := viewer.Find(fsys)
	is.NoErr(err)
	is.Equal(len(pages), 1)
	is.True(pages["posts/index"] != nil)
	is.Equal(pages["posts/index"].Path, "posts/index.svelte")
	is.True(pages["posts/index"].Client != nil)
	is.Equal(pages["posts/index"].Client.Route, "/view/posts/index.svelte.entry.js")
	is.Equal(pages["posts/index"].Client.Path, "posts/index.svelte.entry.js")
	is.True(pages["posts/index"].View.Client != nil)
	is.Equal(pages["posts/index"].View.Client.Route, "/view/posts/index.svelte.js")
	is.Equal(pages["posts/index"].View.Client.Path, "posts/index.svelte.js")
	is.Equal(pages["posts/index"].Route, "/posts")

	// Frames
	is.Equal(len(pages["posts/index"].Frames), 2)
	is.Equal(pages["posts/index"].Frames[0].Key, "posts/frame")
	is.Equal(pages["posts/index"].Frames[0].Path, "posts/frame.svelte")
	is.True(pages["posts/index"].Frames[0].Client != nil)
	is.Equal(pages["posts/index"].Frames[0].Client.Route, "/view/posts/frame.svelte.js")
	is.Equal(pages["posts/index"].Frames[0].Client.Path, "posts/frame.svelte.js")
	is.Equal(pages["posts/index"].Frames[1].Key, "frame")
	is.Equal(pages["posts/index"].Frames[1].Path, "frame.svelte")
	is.True(pages["posts/index"].Frames[1].Client != nil)
	is.Equal(pages["posts/index"].Frames[1].Client.Route, "/view/frame.svelte.js")
	is.Equal(pages["posts/index"].Frames[1].Client.Path, "frame.svelte.js")

	// Error page
	is.Equal(pages["posts/index"].Error, nil)

	// Layout
	is.True(pages["posts/index"].Layout != nil)
	is.Equal(pages["posts/index"].Layout.Key, "layout")
	is.Equal(pages["posts/index"].Layout.Path, "layout.svelte")
	// Layout is server-side only
	is.Equal(pages["posts/index"].Layout.Client, nil)
}

func TestError(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{
		"layout.svelte":      &fstest.MapFile{Data: []byte(`<slot />`)},
		"frame.svelte":       &fstest.MapFile{Data: []byte(`<slot />`)},
		"error.svelte":       &fstest.MapFile{Data: []byte(`<h1>Oops!</h1>`)},
		"posts/frame.svelte": &fstest.MapFile{Data: []byte(`<slot />`)},
		"posts/index.svelte": &fstest.MapFile{Data: []byte(`<h1>Hello {planet}!</h1>`)},
	}
	// Find the pages
	pages, err := viewer.Find(fsys)
	is.NoErr(err)
	is.Equal(len(pages), 2)

	is.True(pages["error"] != nil)
	is.Equal(pages["error"].Key, "error")
	is.Equal(pages["error"].Path, "error.svelte")
	is.True(pages["error"].Client != nil)
	is.Equal(pages["error"].Client.Route, "/view/error.svelte.entry.js")
	is.Equal(pages["error"].Client.Path, "error.svelte.entry.js")
	is.True(pages["error"].View.Client != nil)
	is.Equal(pages["error"].View.Client.Route, "/view/error.svelte.js")
	is.Equal(pages["error"].View.Client.Path, "error.svelte.js")
	is.Equal(pages["error"].Route, "/error")

	// Frames
	is.Equal(len(pages["error"].Frames), 1)
	is.Equal(pages["error"].Frames[0].Key, "frame")
	is.Equal(pages["error"].Frames[0].Path, "frame.svelte")
	is.True(pages["error"].Frames[0].Client != nil)
	is.Equal(pages["error"].Frames[0].Client.Route, "/view/frame.svelte.js")
	is.Equal(pages["error"].Frames[0].Client.Path, "frame.svelte.js")

	// Error page
	is.Equal(pages["error"].Error, nil)

	// Layout
	is.True(pages["error"].Layout != nil)
	is.Equal(pages["error"].Layout.Key, "layout")
	is.Equal(pages["error"].Layout.Path, "layout.svelte")
	// Layout is server-side only
	is.Equal(pages["error"].Layout.Client, nil)
}

func TestComponent(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{
		"Component.svelte": &fstest.MapFile{Data: []byte(`<h1>Hello</h1>`)},
	}
	// Find the pages
	pages, err := viewer.Find(fsys)
	is.NoErr(err)
	is.Equal(len(pages), 0)
}

func TestUnderscore(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{
		"_index.svelte": &fstest.MapFile{Data: []byte(`<h1>Hello</h1>`)},
	}
	// Find the pages
	pages, err := viewer.Find(fsys)
	is.NoErr(err)
	is.Equal(len(pages), 0)
}
