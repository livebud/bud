package entrypoint_test

import (
	"testing"
	"testing/fstest"

	"gitlab.com/mnm/bud/internal/entrypoint"
	"github.com/matryer/is"
)

func TestList(t *testing.T) {
	is := is.New(t)
	// TODO: add view/ to everything. It won't make a difference but it will be
	// more realistic
	fsys := fstest.MapFS{
		"about.jsx":          &fstest.MapFile{},
		"first-post.md":      &fstest.MapFile{},
		"frame.jsx":          &fstest.MapFile{},
		"frame.svelte":       &fstest.MapFile{},
		"error.svelte":       &fstest.MapFile{},
		"index.svelte":       &fstest.MapFile{},
		"layout.jsx":         &fstest.MapFile{},
		"layout.svelte":      &fstest.MapFile{},
		".dot.svelte":        &fstest.MapFile{},
		"_underscore.svelte": &fstest.MapFile{},
		"Component.svelte":   &fstest.MapFile{},
		"user/frame.svelte":  &fstest.MapFile{},
		"user/edit.svelte":   &fstest.MapFile{},
		"user/index.svelte":  &fstest.MapFile{},
		"user/error.svelte":  &fstest.MapFile{},
	}
	views, err := entrypoint.List(fsys)
	is.NoErr(err)
	is.Equal(len(views), 4)
	// about.jsx
	is.Equal(views[0].Page, entrypoint.Path("about.jsx"))
	is.Equal(len(views[0].Frames), 1)
	is.Equal(views[0].Frames[0], entrypoint.Path("frame.jsx"))
	is.Equal(views[0].Layout, entrypoint.Path("layout.jsx"))
	is.Equal(views[0].Type, "jsx")
	is.Equal(views[0].Route, "/about")
	is.Equal(views[0].Client, "bud/_about.jsx")
	// index.svelte
	is.Equal(views[1].Page, entrypoint.Path("index.svelte"))
	is.Equal(len(views[1].Frames), 1)
	is.Equal(views[1].Frames[0], entrypoint.Path("frame.svelte"))
	is.Equal(views[1].Layout, entrypoint.Path("layout.svelte"))
	is.Equal(views[1].Type, "svelte")
	is.Equal(views[1].Route, "/")
	is.Equal(views[1].Client, "bud/_index.svelte")
	// user/edit.svelte
	is.Equal(views[2].Page, entrypoint.Path("user/edit.svelte"))
	is.Equal(len(views[2].Frames), 2)
	is.Equal(views[2].Frames[0], entrypoint.Path("frame.svelte"))
	is.Equal(views[2].Frames[1], entrypoint.Path("user/frame.svelte"))
	is.Equal(views[2].Layout, entrypoint.Path("layout.svelte"))
	is.Equal(views[2].Type, "svelte")
	is.Equal(views[2].Route, "/user/:id/edit")
	is.Equal(views[2].Client, "bud/user/_edit.svelte")
	// user/index.svelte
	is.Equal(views[3].Page, entrypoint.Path("user/index.svelte"))
	is.Equal(len(views[3].Frames), 2)
	is.Equal(views[3].Frames[0], entrypoint.Path("frame.svelte"))
	is.Equal(views[3].Frames[1], entrypoint.Path("user/frame.svelte"))
	is.Equal(views[3].Layout, entrypoint.Path("layout.svelte"))
	is.Equal(views[3].Type, "svelte")
	is.Equal(views[3].Route, "/user")
	is.Equal(views[3].Client, "bud/user/_index.svelte")
}
