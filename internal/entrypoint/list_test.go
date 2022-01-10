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
	fsys := vfs.Map{
		"view/about.jsx":                      "",
		"view/first-post.md":                  "",
		"view/Frame.jsx":                      "",
		"view/Frame.svelte":                   "",
		"view/Error.svelte":                   "",
		"view/index.svelte":                   "",
		"view/Layout.jsx":                     "",
		"view/Layout.svelte":                  "",
		"view/.dot.svelte":                    "",
		"view/_underscore.svelte":             "",
		"view/Component.svelte":               "",
		"view/user/Frame.svelte":              "",
		"view/user/edit.svelte":               "",
		"view/user/index.svelte":              "",
		"view/user/Error.svelte":              "",
		"view/visitor/comments/Frame.svelte":  "",
		"view/visitor/comments/edit.svelte":   "",
		"view/visitor/comments/index.svelte":  "",
		"view/visitor/comments/show.svelte":   "",
		"view/visitor/comments/Error.svelte":  "",
		"view/visitor/comments/Layout.svelte": "",
	}
	views, err := entrypoint.List(fsys)
	is.NoErr(err)
	is.Equal(len(views), 6)
	// index.svelte
	is.Equal(views[0].Page, entrypoint.Path("view/index.svelte"))
	is.Equal(len(views[0].Frames), 1)
	is.Equal(views[0].Frames[0], entrypoint.Path("view/Frame.svelte"))
	is.Equal(views[0].Layout, entrypoint.Path("view/Layout.svelte"))
	is.Equal(views[0].Error, entrypoint.Path("view/Error.svelte"))
	is.Equal(views[0].Type, "svelte")
	is.Equal(views[0].Route, "/")
	is.Equal(views[0].Client, "bud/view/_index.svelte")
	// user/edit.svelte
	is.Equal(views[1].Page, entrypoint.Path("view/user/edit.svelte"))
	is.Equal(len(views[1].Frames), 2)
	is.Equal(views[1].Frames[0], entrypoint.Path("view/Frame.svelte"))
	is.Equal(views[1].Frames[1], entrypoint.Path("view/user/Frame.svelte"))
	is.Equal(views[1].Layout, entrypoint.Path("view/Layout.svelte"))
	is.Equal(views[1].Error, entrypoint.Path("view/user/Error.svelte"))
	is.Equal(views[1].Type, "svelte")
	is.Equal(views[1].Route, "/user/:id/edit")
	is.Equal(views[1].Client, "bud/view/user/_edit.svelte")
	// user/index.svelte
	is.Equal(views[2].Page, entrypoint.Path("view/user/index.svelte"))
	is.Equal(len(views[2].Frames), 2)
	is.Equal(views[2].Frames[0], entrypoint.Path("view/Frame.svelte"))
	is.Equal(views[2].Frames[1], entrypoint.Path("view/user/Frame.svelte"))
	is.Equal(views[2].Layout, entrypoint.Path("view/Layout.svelte"))
	is.Equal(views[2].Error, entrypoint.Path("view/user/Error.svelte"))
	is.Equal(views[2].Type, "svelte")
	is.Equal(views[2].Route, "/user")
	is.Equal(views[2].Client, "bud/view/user/_index.svelte")
	// visitor/comments/index.svelte
	is.Equal(views[3].Page, entrypoint.Path("view/visitor/comments/edit.svelte"))
	is.Equal(len(views[3].Frames), 2)
	is.Equal(views[3].Frames[0], entrypoint.Path("view/Frame.svelte"))
	is.Equal(views[3].Frames[1], entrypoint.Path("view/visitor/comments/Frame.svelte"))
	is.Equal(views[3].Layout, entrypoint.Path("view/visitor/comments/Layout.svelte"))
	is.Equal(views[3].Error, entrypoint.Path("view/visitor/comments/Error.svelte"))
	is.Equal(views[3].Type, "svelte")
	is.Equal(views[3].Route, "/visitor/:visitor_id/comments/:id/edit")
	is.Equal(views[3].Client, "bud/view/visitor/comments/_edit.svelte")
}

func TestListUnderscore(t *testing.T) {
	is := is.New(t)
	// TODO: add view/ to everything. It won't make a difference but it will be
	// more realistic
	fsys := vfs.Map{
		"admin_users/comments/show.svelte": "",
	}
	views, err := entrypoint.List(fsys)
	is.NoErr(err)
	is.Equal(len(views), 1)
	is.Equal(views[0].Page, entrypoint.Path("admin_users/comments/show.svelte"))
	is.Equal(len(views[0].Frames), 0)
	is.Equal(views[0].Layout, entrypoint.Path(""))
	is.Equal(views[0].Error, entrypoint.Path(""))
	is.Equal(views[0].Type, "svelte")
	is.Equal(views[0].Route, "/admin_users/:admin_user_id/comments/:id")
	is.Equal(views[0].Client, "bud/admin_users/comments/_show.svelte")
}
