package view_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"
	"testing/fstest"

	"github.com/livebud/bud/pkg/view"
	"github.com/livebud/bud/pkg/view/gohtml"
	"github.com/livebud/bud/pkg/view/markdown"
	"github.com/matryer/is"
	"github.com/matthewmueller/diff"
)

func TestFindSample(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{
		"index.gohtml": &fstest.MapFile{Data: []byte("Hello {{ .Planet }}!")},
	}
	// Find the pages
	finder := view.New(fsys, map[string]view.Renderer{
		".gohtml": gohtml.New(),
	})
	page, err := finder.FindPage("index")
	is.NoErr(err)
	is.True(page.View != nil)
	is.Equal(page.View.Key(), "index")
	is.Equal(page.View.Path(), "index.gohtml")
	is.Equal(len(page.Frames), 0)
	is.Equal(page.Layout, nil)
	is.Equal(page.Error, nil)
}

func TestFindNested(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{
		"layout.gohtml":               &fstest.MapFile{Data: []byte(`<slot />`)},
		"frame.gohtml":                &fstest.MapFile{Data: []byte(`<slot />`)},
		"posts/frame.gohtml":          &fstest.MapFile{Data: []byte(`<slot />`)},
		"posts/index.gohtml":          &fstest.MapFile{Data: []byte(`<h1>Hello {planet}!</h1>`)},
		"posts/comments/frame.gohtml": &fstest.MapFile{Data: []byte(`<slot />`)},
		"posts/comments/index.gohtml": &fstest.MapFile{Data: []byte(`<h1>Hello {planet}!</h1>`)},
	}
	// Find the pages
	finder := view.New(fsys, map[string]view.Renderer{
		".gohtml": gohtml.New(),
	})
	// posts/index
	page, err := finder.FindPage("posts/index")
	is.NoErr(err)
	is.True(page.View != nil)
	is.Equal(page.View.Key(), "posts/index")
	is.Equal(page.View.Path(), "posts/index.gohtml")
	is.Equal(len(page.Frames), 2)
	is.Equal(page.Frames[0].Key(), "frame")
	is.Equal(page.Frames[0].Path(), "frame.gohtml")
	is.Equal(page.Frames[1].Key(), "posts/frame")
	is.Equal(page.Frames[1].Path(), "posts/frame.gohtml")
	is.Equal(page.Error, nil)
	is.True(page.Layout != nil)
	is.Equal(page.Layout.Key(), "layout")
	is.Equal(page.Layout.Path(), "layout.gohtml")
	// posts/comments/index
	page, err = finder.FindPage("posts/comments/index")
	is.NoErr(err)
	is.True(page.View != nil)
	is.Equal(page.View.Key(), "posts/comments/index")
	is.Equal(page.View.Path(), "posts/comments/index.gohtml")
	is.Equal(len(page.Frames), 3)
	is.Equal(page.Frames[0].Key(), "frame")
	is.Equal(page.Frames[0].Path(), "frame.gohtml")
	is.Equal(page.Frames[1].Key(), "posts/frame")
	is.Equal(page.Frames[1].Path(), "posts/frame.gohtml")
	is.Equal(page.Frames[2].Key(), "posts/comments/frame")
	is.Equal(page.Frames[2].Path(), "posts/comments/frame.gohtml")
	is.Equal(page.Error, nil)
	is.True(page.Layout != nil)
	is.Equal(page.Layout.Key(), "layout")
	is.Equal(page.Layout.Path(), "layout.gohtml")
}

func TestFindView(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{
		"layout.gohtml":               &fstest.MapFile{Data: []byte(`<slot />`)},
		"frame.gohtml":                &fstest.MapFile{Data: []byte(`<slot />`)},
		"posts/frame.gohtml":          &fstest.MapFile{Data: []byte(`<slot />`)},
		"posts/index.gohtml":          &fstest.MapFile{Data: []byte(`<h1>Hello {planet}!</h1>`)},
		"posts/comments/frame.gohtml": &fstest.MapFile{Data: []byte(`<slot />`)},
		"posts/comments/index.gohtml": &fstest.MapFile{Data: []byte(`<h1>Hello {planet}!</h1>`)},
	}
	// Find the pages
	finder := view.New(fsys, map[string]view.Renderer{
		".gohtml": gohtml.New(),
	})
	// layout
	v, err := finder.FindView("layout")
	is.NoErr(err)
	is.True(v != nil)
	is.Equal(v.Key(), "layout")
	is.Equal(v.Path(), "layout.gohtml")
	// frame
	v, err = finder.FindView("frame")
	is.NoErr(err)
	is.True(v != nil)
	is.Equal(v.Key(), "frame")
	is.Equal(v.Path(), "frame.gohtml")
	// posts/index
	v, err = finder.FindView("posts/index")
	is.NoErr(err)
	is.True(v != nil)
	is.Equal(v.Key(), "posts/index")
	is.Equal(v.Path(), "posts/index.gohtml")
	// posts/frame
	v, err = finder.FindView("posts/frame")
	is.NoErr(err)
	is.True(v != nil)
	is.Equal(v.Key(), "posts/frame")
	// posts/comments/index
	v, err = finder.FindView("posts/comments/index")
	is.NoErr(err)
	is.True(v != nil)
	is.Equal(v.Key(), "posts/comments/index")
	is.Equal(v.Path(), "posts/comments/index.gohtml")
	// posts/comments/frame
	v, err = finder.FindView("posts/comments/frame")
	is.NoErr(err)
	is.True(v != nil)
	is.Equal(v.Key(), "posts/comments/frame")
	is.Equal(v.Path(), "posts/comments/frame.gohtml")
	// posts/layout
	v, err = finder.FindView("posts/layout")
	is.True(err != nil)
	is.Equal(v, nil)
	is.True(errors.Is(err, view.ErrNotFound))
}

func newSlotBuffer() *slotBuffer {
	defaultBuffer := new(bytes.Buffer)
	buffers := map[string]*bytes.Buffer{}
	return &slotBuffer{defaultBuffer, buffers}
}

type slotBuffer struct {
	*bytes.Buffer
	buffers map[string]*bytes.Buffer
}

var _ view.Slot = (*slotBuffer)(nil)

func (s *slotBuffer) Reader(name string) io.Reader {
	if s.buffers[name] == nil {
		s.buffers[name] = new(bytes.Buffer)
	}
	return s.buffers[name]
}

func (s *slotBuffer) Writer(name string) io.Writer {
	if s.buffers[name] == nil {
		s.buffers[name] = new(bytes.Buffer)
	}
	return s.buffers[name]
}

func TestRenderSample(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{
		"index.gohtml": &fstest.MapFile{Data: []byte("Hello {{ .Planet }}!")},
		"about.md":     &fstest.MapFile{Data: []byte("# About Me")},
	}
	renderers := map[string]view.Renderer{
		".md":     markdown.New(),
		".gohtml": gohtml.New(),
	}
	vf := view.New(fsys, renderers)
	ctx := context.Background()
	writer := newSlotBuffer()
	page, err := vf.FindPage("index")
	is.NoErr(err)
	err = page.View.Render(ctx, writer, map[string]interface{}{"Planet": "World"})
	is.NoErr(err)
	diff.TestString(t, writer.String(), "Hello World!")
	writer = newSlotBuffer()
	page, err = vf.FindPage("about")
	is.NoErr(err)
	err = page.View.Render(ctx, writer, nil)
	is.NoErr(err)
	diff.TestString(t, writer.String(), "<h1>About Me</h1>\n")
}

func TestRenderMultiple(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{
		"index.html": &fstest.MapFile{Data: []byte("index")},
		"edit.html":  &fstest.MapFile{Data: []byte(`edit`)},
	}
	renderers := map[string]view.Renderer{
		".html": gohtml.New(),
	}
	vf := view.New(fsys, renderers)
	index, err := vf.FindPage("index")
	is.NoErr(err)
	edit, err := vf.FindPage("edit")
	is.NoErr(err)
	is.Equal(index.View.Key(), "index")
	is.Equal(index.View.Path(), "index.html")
	is.Equal(edit.View.Key(), "edit")
	is.Equal(edit.View.Path(), "edit.html")
	indexBytes := newSlotBuffer()
	ctx := context.Background()
	err = index.View.Render(ctx, indexBytes, nil)
	is.NoErr(err)
	diff.TestString(t, indexBytes.String(), "index")
	editBytes := newSlotBuffer()
	err = edit.View.Render(ctx, editBytes, nil)
	is.NoErr(err)
	diff.TestString(t, editBytes.String(), "edit")
}
