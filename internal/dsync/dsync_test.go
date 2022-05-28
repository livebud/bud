package dsync_test

import (
	"errors"
	"io/fs"
	"path/filepath"
	"testing"
	"time"

	"github.com/livebud/bud/internal/dsync"
	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/conjure"
	"github.com/livebud/bud/package/vfs"
)

func TestFileSync(t *testing.T) {
	is := is.New(t)
	before := time.Date(2021, 8, 4, 14, 56, 0, 0, time.UTC)
	after := time.Date(2021, 8, 4, 14, 57, 0, 0, time.UTC)
	vfs.Now = func() time.Time { return after }

	// starting points
	sourceFS := vfs.Memory{
		"a.txt": &vfs.File{Data: []byte("a")},
		"b.txt": &vfs.File{Data: []byte("b")},
	}
	targetFS := vfs.Memory{
		"b.txt": &vfs.File{Data: []byte("bb"), ModTime: before},
		"c.txt": &vfs.File{Data: []byte("c"), ModTime: before},
	}

	// sync
	err := dsync.Dir(sourceFS, ".", targetFS, ".")
	is.NoErr(err)
	is.Equal(len(targetFS), 2)

	// a.txt
	code, err := fs.ReadFile(targetFS, "a.txt")
	is.NoErr(err)
	is.Equal(string(code), `a`)
	stat, err := fs.Stat(targetFS, "a.txt")
	is.NoErr(err)
	is.Equal(stat.IsDir(), false)
	is.Equal(stat.Mode(), fs.FileMode(0644))
	is.True(stat.ModTime().Equal(after))

	// b.txt
	code, err = fs.ReadFile(targetFS, "b.txt")
	is.NoErr(err)
	is.Equal(string(code), `b`)
	stat, err = fs.Stat(targetFS, "b.txt")
	is.NoErr(err)
	is.Equal(stat.IsDir(), false)
	is.Equal(stat.Mode(), fs.FileMode(0644))
	is.True(stat.ModTime().Equal(after))
}

func TestDirSync(t *testing.T) {
	is := is.New(t)
	before := time.Date(2021, 8, 4, 14, 56, 0, 0, time.UTC)
	after := time.Date(2021, 8, 4, 14, 57, 0, 0, time.UTC)
	vfs.Now = func() time.Time { return after }

	// starting points
	sourceFS := vfs.Memory{
		"duo/view/index.svelte":        &vfs.File{Data: []byte("<h1>index</h1>"), ModTime: after},
		"duo/view/about/about.svelte":  &vfs.File{Data: []byte("<h1>about</h1>"), ModTime: after},
		"duo/view/user/user.svelte":    &vfs.File{Data: []byte("<h1>user</h1>"), ModTime: before},
		"duo/controller/controller.go": &vfs.File{Data: []byte("package controller"), ModTime: after},
	}
	targetFS := vfs.Memory{
		"duo/view/index.svelte":       &vfs.File{Data: []byte("<h1>indexx</h1>"), ModTime: before},
		"duo/view/about/index.svelte": &vfs.File{Data: []byte("<h1>about</h1>"), ModTime: before},
		"duo/view/user/user.svelte":   &vfs.File{Data: []byte("<h1>user</h1>"), ModTime: before},
		"duo/session/session.go":      &vfs.File{Data: []byte("package session"), ModTime: before},
	}

	// sync
	err := dsync.Dir(sourceFS, "duo", targetFS, "duo")
	is.NoErr(err)
	is.Equal(len(targetFS), 5)

	// duo/view/index.svelte
	_, ok := targetFS["duo/view/index.svelte"]
	is.True(ok)
	code, err := fs.ReadFile(targetFS, "duo/view/index.svelte")
	is.NoErr(err)
	is.Equal(string(code), `<h1>index</h1>`)
	stat, err := fs.Stat(targetFS, "duo/view/index.svelte")
	is.NoErr(err)
	is.Equal(stat.IsDir(), false)
	is.Equal(stat.Mode(), fs.FileMode(0644))
	is.True(stat.ModTime().Equal(after))

	// duo/view/about/about.svelte
	_, ok = targetFS["duo/view/about/about.svelte"]
	is.True(ok)
	code, err = fs.ReadFile(targetFS, "duo/view/about/about.svelte")
	is.NoErr(err)
	is.Equal(string(code), `<h1>about</h1>`)
	stat, err = fs.Stat(targetFS, "duo/view/about/about.svelte")
	is.NoErr(err)
	is.Equal(stat.IsDir(), false)
	is.Equal(stat.Mode(), fs.FileMode(0644))
	is.True(stat.ModTime().Equal(after))

	// duo/view/user/user.svelte
	_, ok = targetFS["duo/view/user/user.svelte"]
	is.True(ok)
	code, err = fs.ReadFile(targetFS, "duo/view/user/user.svelte")
	is.NoErr(err)
	is.Equal(string(code), `<h1>user</h1>`)
	stat, err = fs.Stat(targetFS, "duo/view/user/user.svelte")
	is.NoErr(err)
	is.Equal(stat.IsDir(), false)
	is.Equal(stat.Mode(), fs.FileMode(0))
	is.True(stat.ModTime().Equal(before))

	// duo/controller
	_, ok = targetFS["duo/controller"]
	is.True(ok)
	stat, err = fs.Stat(targetFS, "duo/controller")
	is.NoErr(err)
	is.Equal(stat.IsDir(), true)
	is.Equal(stat.Mode(), fs.FileMode(0755|fs.ModeDir))
	is.True(stat.ModTime().Equal(after))

	// duo/controller/controller.go
	_, ok = targetFS["duo/controller/controller.go"]
	is.True(ok)
	code, err = fs.ReadFile(targetFS, "duo/controller/controller.go")
	is.NoErr(err)
	is.Equal(string(code), `package controller`)
	stat, err = fs.Stat(targetFS, "duo/controller/controller.go")
	is.NoErr(err)
	is.Equal(stat.IsDir(), false)
	is.Equal(stat.Mode(), fs.FileMode(0644))
	is.True(stat.ModTime().Equal(after))
}

func TestNoDuo(t *testing.T) {
	is := is.New(t)
	// before := time.Date(2021, 8, 4, 14, 56, 0, 0, time.UTC)
	after := time.Date(2021, 8, 4, 14, 57, 0, 0, time.UTC)
	vfs.Now = func() time.Time { return after }

	// starting points
	sourceFS := vfs.Memory{
		"duo/view/view.go": &vfs.File{Data: []byte("package view"), ModTime: after},
	}
	targetFS := vfs.Memory{}

	// sync
	err := dsync.Dir(sourceFS, "duo", targetFS, "duo")
	is.NoErr(err)
	is.Equal(len(targetFS), 2)

	// duo/view
	_, ok := targetFS["duo/view"]
	is.True(ok)
	stat, err := fs.Stat(targetFS, "duo/view")
	is.NoErr(err)
	is.Equal(stat.IsDir(), true)
	is.Equal(stat.Mode(), fs.FileMode(0755|fs.ModeDir))
	is.True(stat.ModTime().Equal(after))

	// duo/view/view.go
	_, ok = targetFS["duo/view/view.go"]
	is.True(ok)
	code, err := fs.ReadFile(targetFS, "duo/view/view.go")
	is.NoErr(err)
	is.Equal(string(code), `package view`)
	stat, err = fs.Stat(targetFS, "duo/view/view.go")
	is.NoErr(err)
	is.Equal(stat.IsDir(), false)
	is.Equal(stat.Mode(), fs.FileMode(0644))
	is.True(stat.ModTime().Equal(after))
}

func TestSkipNotExist(t *testing.T) {
	is := is.New(t)
	// before := time.Date(2021, 8, 4, 14, 56, 0, 0, time.UTC)
	after := time.Date(2021, 8, 4, 14, 57, 0, 0, time.UTC)
	vfs.Now = func() time.Time { return after }

	// starting points
	sourceFS := conjure.New()
	sourceFS.GenerateFile("bud/generate/main.go", func(file *conjure.File) error {
		return fs.ErrNotExist
	})
	targetFS := vfs.Memory{}

	// sync
	err := dsync.Dir(sourceFS, ".", targetFS, ".")
	is.NoErr(err)
	is.Equal(len(targetFS), 0)
}

func TestErrorGenerator(t *testing.T) {
	is := is.New(t)
	// before := time.Date(2021, 8, 4, 14, 56, 0, 0, time.UTC)
	after := time.Date(2021, 8, 4, 14, 57, 0, 0, time.UTC)
	vfs.Now = func() time.Time { return after }

	// starting points
	sourceFS := conjure.New()
	sourceFS.GenerateFile("bud/generate/main.go", func(file *conjure.File) error {
		return errors.New("uh oh")
	})
	targetFS := vfs.Memory{}

	// sync
	err := dsync.Dir(sourceFS, ".", targetFS, ".")
	is.True(err != nil)
	is.Equal(err.Error(), `conjure: generate "bud/generate/main.go". uh oh`)
	is.Equal(len(targetFS), 0)
}

func TestWithSkip(t *testing.T) {
	is := is.New(t)
	// starting points
	sourceFS := vfs.Memory{
		"index.svelte": &vfs.File{Data: []byte("<h1>index</h1>")},
	}
	targetFS := vfs.Memory{
		"node_modules/svelte/svelte.js": &vfs.File{Data: []byte("svelte")},
	}
	err := dsync.Dir(sourceFS, ".", targetFS, ".")
	is.NoErr(err)
	is.Equal(len(targetFS), 1) // this should have deleted node_modules
	// starting points
	sourceFS = vfs.Memory{
		"index.svelte":      &vfs.File{Data: []byte("<h1>index</h1>")},
		"bud/controller.go": &vfs.File{Data: []byte("package controller")},
	}
	targetFS = vfs.Memory{
		"node_modules/svelte/svelte.js": &vfs.File{Data: []byte("svelte")},
		"bud/generate.go":               &vfs.File{Data: []byte("package main")},
	}
	skip1 := func(name string, isDir bool) bool {
		return isDir && filepath.Base(name) == "node_modules"
	}
	// NOTE: if you don't have bud/controller.go
	skip2 := func(name string, isDir bool) bool {
		return !isDir && name == "bud/generate.go"
	}
	err = dsync.Dir(sourceFS, ".", targetFS, ".", dsync.WithSkip(skip1, skip2))
	is.NoErr(err)
	is.Equal(len(targetFS), 4) // this should have kept node_modules & generate
}

func TestAvoidDotCreate(t *testing.T) {
	is := is.New(t)
	// starting points
	sourceFS := vfs.Memory{
		".": &vfs.File{Mode: fs.ModeDir},
	}
	targetFS := vfs.Memory{}
	err := dsync.Dir(sourceFS, ".", targetFS, ".")
	is.NoErr(err)
	is.Equal(len(targetFS), 0)
}

func TestAvoidDotUpdate(t *testing.T) {
	is := is.New(t)
	// starting points
	sourceFS := vfs.Memory{
		".": &vfs.File{Mode: fs.ModeDir},
	}
	targetFS := vfs.Memory{
		".": &vfs.File{Mode: fs.ModeDir | 0755},
	}
	err := dsync.Dir(sourceFS, ".", targetFS, ".")
	is.NoErr(err)
	is.Equal(len(targetFS), 1)
}

// Avoid deleting the root of the target fs
func TestAvoidDotDelete(t *testing.T) {
	is := is.New(t)
	// starting points
	sourceFS := vfs.Memory{}
	targetFS := vfs.Memory{
		".": &vfs.File{Mode: fs.ModeDir},
	}
	err := dsync.Dir(sourceFS, ".", targetFS, ".")
	is.NoErr(err)
	// . should be ignored
	is.Equal(len(targetFS), 1)
}

// Relative path from source path
func TestRelativeSource(t *testing.T) {
	is := is.New(t)
	// starting points
	sourceFS := vfs.Memory{
		"bud/.cli/main.go": &vfs.File{Data: []byte("package main")},
		"bud/.cli/a/a.go":  &vfs.File{Data: []byte("package a")},
	}
	targetFS := vfs.Memory{
		"a/a.go": &vfs.File{Data: []byte("package aa")},
	}
	err := dsync.Dir(sourceFS, "bud/.cli", targetFS, ".")
	is.NoErr(err)
	_, ok := targetFS["main.go"]
	is.True(ok) // missing main.go
	data, err := fs.ReadFile(targetFS, "a/a.go")
	is.NoErr(err)
	is.Equal(string(data), "package a")
}

func TestRel(t *testing.T) {
	is := is.New(t)
	rel, err := dsync.Rel("bud/.cli", ".")("bud/.cli/a/a.go")
	is.NoErr(err)
	is.Equal(rel, "a/a.go")
	rel, err = dsync.Rel(".", ".")("a/a.go")
	is.NoErr(err)
	is.Equal(rel, "a/a.go")
	rel, err = dsync.Rel(".", "bud/.cli")("a/a.go")
	is.NoErr(err)
	is.Equal(rel, "bud/.cli/a/a.go")
	rel, err = dsync.Rel("bud", "app")("bud/a/a.go")
	is.NoErr(err)
	is.Equal(rel, "app/a/a.go")
}
