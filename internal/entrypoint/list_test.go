package entrypoint_test

import (
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/internal/entrypoint"
	"gitlab.com/mnm/bud/vfs"
)

func TestList(t *testing.T) {
	is := is.New(t)
	// TODO: add view/ to everything. It won't make a difference but it will be
	// more realistic
	fsys := vfs.Memory{
		"about.jsx":                      &vfs.File{},
		"first-post.md":                  &vfs.File{},
		"Frame.jsx":                      &vfs.File{},
		"Frame.svelte":                   &vfs.File{},
		"Error.svelte":                   &vfs.File{},
		"index.svelte":                   &vfs.File{},
		"Layout.jsx":                     &vfs.File{},
		"Layout.svelte":                  &vfs.File{},
		".dot.svelte":                    &vfs.File{},
		"_underscore.svelte":             &vfs.File{},
		"Component.svelte":               &vfs.File{},
		"user/Frame.svelte":              &vfs.File{},
		"user/edit.svelte":               &vfs.File{},
		"user/index.svelte":              &vfs.File{},
		"user/Error.svelte":              &vfs.File{},
		"visitor/comments/Frame.svelte":  &vfs.File{},
		"visitor/comments/edit.svelte":   &vfs.File{},
		"visitor/comments/index.svelte":  &vfs.File{},
		"visitor/comments/show.svelte":   &vfs.File{},
		"visitor/comments/Error.svelte":  &vfs.File{},
		"visitor/comments/Layout.svelte": &vfs.File{},
	}
	views, err := entrypoint.List(fsys)
	is.NoErr(err)
	is.Equal(len(views), 8)
	// about.jsx
	is.Equal(views[0].Page, entrypoint.Path("about.jsx"))
	is.Equal(len(views[0].Frames), 1)
	is.Equal(views[0].Frames[0], entrypoint.Path("Frame.jsx"))
	is.Equal(views[0].Layout, entrypoint.Path("Layout.jsx"))
	is.Equal(views[0].Error, entrypoint.Path(""))
	is.Equal(views[0].Type, "jsx")
	is.Equal(views[0].Route, "/about")
	is.Equal(views[0].Client, "bud/_about.jsx")
	// index.svelte
	is.Equal(views[2].Page, entrypoint.Path("index.svelte"))
	is.Equal(len(views[2].Frames), 1)
	is.Equal(views[2].Frames[0], entrypoint.Path("Frame.svelte"))
	is.Equal(views[2].Layout, entrypoint.Path("Layout.svelte"))
	is.Equal(views[2].Error, entrypoint.Path("Error.svelte"))
	is.Equal(views[2].Type, "svelte")
	is.Equal(views[2].Route, "/")
	is.Equal(views[2].Client, "bud/_index.svelte")
	// user/edit.svelte
	is.Equal(views[3].Page, entrypoint.Path("user/edit.svelte"))
	is.Equal(len(views[3].Frames), 2)
	is.Equal(views[3].Frames[0], entrypoint.Path("Frame.svelte"))
	is.Equal(views[3].Frames[1], entrypoint.Path("user/Frame.svelte"))
	is.Equal(views[3].Layout, entrypoint.Path("Layout.svelte"))
	is.Equal(views[3].Error, entrypoint.Path("user/Error.svelte"))
	is.Equal(views[3].Type, "svelte")
	is.Equal(views[3].Route, "/user/:id/edit")
	is.Equal(views[3].Client, "bud/user/_edit.svelte")
	// user/index.svelte
	is.Equal(views[4].Page, entrypoint.Path("user/index.svelte"))
	is.Equal(len(views[4].Frames), 2)
	is.Equal(views[4].Frames[0], entrypoint.Path("Frame.svelte"))
	is.Equal(views[4].Frames[1], entrypoint.Path("user/Frame.svelte"))
	is.Equal(views[4].Layout, entrypoint.Path("Layout.svelte"))
	is.Equal(views[4].Error, entrypoint.Path("user/Error.svelte"))
	is.Equal(views[4].Type, "svelte")
	is.Equal(views[4].Route, "/user")
	is.Equal(views[4].Client, "bud/user/_index.svelte")
	// visitor/comments/index.svelte
	is.Equal(views[5].Page, entrypoint.Path("visitor/comments/edit.svelte"))
	is.Equal(len(views[5].Frames), 2)
	is.Equal(views[5].Frames[0], entrypoint.Path("Frame.svelte"))
	is.Equal(views[5].Frames[1], entrypoint.Path("visitor/comments/Frame.svelte"))
	is.Equal(views[5].Layout, entrypoint.Path("visitor/comments/Layout.svelte"))
	is.Equal(views[5].Error, entrypoint.Path("visitor/comments/Error.svelte"))
	is.Equal(views[5].Type, "svelte")
	is.Equal(views[5].Route, "/visitor/:visitor_id/comments/:id/edit")
	is.Equal(views[5].Client, "bud/visitor/comments/_edit.svelte")
}
