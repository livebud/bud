package virtual_test

import (
	"errors"
	"io/fs"
	"strings"
	"testing"
	"time"

	"github.com/livebud/bud/internal/dag"
	"github.com/livebud/bud/internal/dsync"
	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/genfs"
	"github.com/livebud/bud/package/log/testlog"
	"github.com/livebud/bud/package/virtual"
)

func TestSyncFile(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	before := time.Date(2021, 8, 4, 14, 56, 0, 0, time.UTC)
	after := time.Date(2021, 8, 4, 14, 57, 0, 0, time.UTC)
	virtual.Now = func() time.Time { return after }

	// starting points
	from := virtual.Tree{
		"a.txt": &virtual.File{Data: []byte("a")},
		"b.txt": &virtual.File{Data: []byte("b")},
	}
	to := virtual.Tree{
		"b.txt": &virtual.File{Data: []byte("bb"), ModTime: before},
		"c.txt": &virtual.File{Data: []byte("c"), ModTime: before},
	}

	// sync
	err := virtual.Sync(log, from, to, ".")
	is.NoErr(err)
	is.Equal(len(to), 2)

	// a.txt
	code, err := fs.ReadFile(to, "a.txt")
	is.NoErr(err)
	is.Equal(string(code), `a`)
	stat, err := fs.Stat(to, "a.txt")
	is.NoErr(err)
	is.Equal(stat.IsDir(), false)
	is.Equal(stat.Mode(), fs.FileMode(0644))
	is.True(stat.ModTime().Equal(after))

	// b.txt
	code, err = fs.ReadFile(to, "b.txt")
	is.NoErr(err)
	is.Equal(string(code), `b`)
	stat, err = fs.Stat(to, "b.txt")
	is.NoErr(err)
	is.Equal(stat.IsDir(), false)
	is.Equal(stat.Mode(), fs.FileMode(0644))
	is.True(stat.ModTime().Equal(after))
}

func TestSyncDir(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	before := time.Date(2021, 8, 4, 14, 56, 0, 0, time.UTC)
	after := time.Date(2021, 8, 4, 14, 57, 0, 0, time.UTC)
	virtual.Now = func() time.Time { return after }

	// starting points
	from := virtual.Tree{
		"duo/view/index.svelte":        &virtual.File{Data: []byte("<h1>index</h1>"), ModTime: after},
		"duo/view/about/about.svelte":  &virtual.File{Data: []byte("<h1>about</h1>"), ModTime: after},
		"duo/view/user/user.svelte":    &virtual.File{Data: []byte("<h1>user</h1>"), ModTime: before},
		"duo/controller/controller.go": &virtual.File{Data: []byte("package controller"), ModTime: after},
	}
	to := virtual.Tree{
		"duo/view/index.svelte":       &virtual.File{Data: []byte("<h1>indexx</h1>"), ModTime: before},
		"duo/view/about/index.svelte": &virtual.File{Data: []byte("<h1>about</h1>"), ModTime: before},
		"duo/view/user/user.svelte":   &virtual.File{Data: []byte("<h1>user</h1>"), ModTime: before},
		"duo/session/session.go":      &virtual.File{Data: []byte("package session"), ModTime: before},
	}

	// sync
	err := virtual.Sync(log, from, to, "duo")
	is.NoErr(err)
	is.Equal(len(to), 5)

	// duo/view/index.svelte
	_, ok := to["duo/view/index.svelte"]
	is.True(ok)
	code, err := fs.ReadFile(to, "duo/view/index.svelte")
	is.NoErr(err)
	is.Equal(string(code), `<h1>index</h1>`)
	stat, err := fs.Stat(to, "duo/view/index.svelte")
	is.NoErr(err)
	is.Equal(stat.IsDir(), false)
	is.Equal(stat.Mode(), fs.FileMode(0644))
	is.True(stat.ModTime().Equal(after))

	// duo/view/about/about.svelte
	_, ok = to["duo/view/about/about.svelte"]
	is.True(ok)
	code, err = fs.ReadFile(to, "duo/view/about/about.svelte")
	is.NoErr(err)
	is.Equal(string(code), `<h1>about</h1>`)
	stat, err = fs.Stat(to, "duo/view/about/about.svelte")
	is.NoErr(err)
	is.Equal(stat.IsDir(), false)
	is.Equal(stat.Mode(), fs.FileMode(0644))
	is.True(stat.ModTime().Equal(after))

	// duo/view/user/user.svelte
	_, ok = to["duo/view/user/user.svelte"]
	is.True(ok)
	code, err = fs.ReadFile(to, "duo/view/user/user.svelte")
	is.NoErr(err)
	is.Equal(string(code), `<h1>user</h1>`)
	stat, err = fs.Stat(to, "duo/view/user/user.svelte")
	is.NoErr(err)
	is.Equal(stat.IsDir(), false)
	is.Equal(stat.Mode(), fs.FileMode(0))
	is.True(stat.ModTime().Equal(before))

	// duo/controller
	_, ok = to["duo/controller"]
	is.True(ok)
	stat, err = fs.Stat(to, "duo/controller")
	is.NoErr(err)
	is.Equal(stat.IsDir(), true)
	is.Equal(stat.Mode(), fs.FileMode(0755|fs.ModeDir))
	is.True(stat.ModTime().Equal(after))

	// duo/controller/controller.go
	_, ok = to["duo/controller/controller.go"]
	is.True(ok)
	code, err = fs.ReadFile(to, "duo/controller/controller.go")
	is.NoErr(err)
	is.Equal(string(code), `package controller`)
	stat, err = fs.Stat(to, "duo/controller/controller.go")
	is.NoErr(err)
	is.Equal(stat.IsDir(), false)
	is.Equal(stat.Mode(), fs.FileMode(0644))
	is.True(stat.ModTime().Equal(after))
}

func TestSyncNoDuo(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	// before := time.Date(2021, 8, 4, 14, 56, 0, 0, time.UTC)
	after := time.Date(2021, 8, 4, 14, 57, 0, 0, time.UTC)
	virtual.Now = func() time.Time { return after }

	// starting points
	from := virtual.Tree{
		"duo/view/view.go": &virtual.File{Data: []byte("package view"), ModTime: after},
	}
	to := virtual.Tree{}

	// sync
	err := virtual.Sync(log, from, to, "duo")
	is.NoErr(err)
	is.Equal(len(to), 2)

	// duo/view
	_, ok := to["duo/view"]
	is.True(ok)
	stat, err := fs.Stat(to, "duo/view")
	is.NoErr(err)
	is.Equal(stat.IsDir(), true)
	is.Equal(stat.Mode(), fs.FileMode(0755|fs.ModeDir))
	is.True(stat.ModTime().Equal(after))

	// duo/view/view.go
	_, ok = to["duo/view/view.go"]
	is.True(ok)
	code, err := fs.ReadFile(to, "duo/view/view.go")
	is.NoErr(err)
	is.Equal(string(code), `package view`)
	stat, err = fs.Stat(to, "duo/view/view.go")
	is.NoErr(err)
	is.Equal(stat.IsDir(), false)
	is.Equal(stat.Mode(), fs.FileMode(0644))
	is.True(stat.ModTime().Equal(after))
}

func TestSyncSkipNotExist(t *testing.T) {
	is := is.New(t)
	log := testlog.New()

	// starting points
	from := genfs.New(dag.Discard, virtual.List{}, log)
	from.GenerateFile("bud/generate/main.go", func(fsys genfs.FS, file *genfs.File) error {
		return fs.ErrNotExist
	})
	to := virtual.List{}

	// sync
	err := virtual.Sync(log, from, &to, ".")
	is.NoErr(err)
	is.Equal(len(to), 0)
}

func TestSyncSkipDirNotExist(t *testing.T) {
	is := is.New(t)
	log := testlog.New()

	// starting points
	from := genfs.New(dag.Discard, virtual.List{}, log)
	from.GenerateDir("bud/generate", func(fsys genfs.FS, dir *genfs.Dir) error {
		return fs.ErrNotExist
	})
	to := virtual.List{}

	// sync
	err := virtual.Sync(log, from, &to, ".")
	is.NoErr(err)
	is.Equal(len(to), 0)
}

func TestSyncErrorGenerator(t *testing.T) {
	is := is.New(t)
	log := testlog.New()

	// before := time.Date(2021, 8, 4, 14, 56, 0, 0, time.UTC)
	after := time.Date(2021, 8, 4, 14, 57, 0, 0, time.UTC)
	virtual.Now = func() time.Time { return after }

	// starting points
	from := genfs.New(dag.Discard, virtual.List{}, log)
	from.GenerateFile("bud/generate/main.go", func(fsys genfs.FS, file *genfs.File) error {
		return errors.New("uh oh")
	})
	to := virtual.List{}

	// sync
	err := virtual.Sync(log, from, &to, ".")
	is.True(err != nil)
	is.In(err.Error(), `uh oh`)
	is.True(!errors.Is(err, fs.ErrNotExist))
	is.Equal(len(to), 0)
}

func TestSyncExclude(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	// starting points
	from := virtual.Tree{
		"index.svelte": &virtual.File{Data: []byte("<h1>index</h1>")},
	}
	to := virtual.Tree{
		"node_modules/svelte/svelte.js": &virtual.File{Data: []byte("svelte")},
	}
	err := virtual.Sync(log, from, to, ".")
	is.NoErr(err)
	is.Equal(len(to), 1) // this should have deleted node_modules
	// starting points
	from = virtual.Tree{
		"index.svelte":      &virtual.File{Data: []byte("<h1>index</h1>")},
		"bud/controller.go": &virtual.File{Data: []byte("package controller")},
	}
	to = virtual.Tree{
		"node_modules/svelte/svelte.js": &virtual.File{Data: []byte("svelte")},
		"bud/generate.go":               &virtual.File{Data: []byte("package main")},
	}
	excluded := virtual.Exclude(to, func(path string) bool {
		return strings.HasPrefix(path, "node_modules") ||
			path == "bud/generate.go"
	})
	err = virtual.Sync(log, from, excluded, ".")
	is.NoErr(err)
	is.Equal(len(to), 4) // this should have kept node_modules & generate
}

func TestSyncAvoidDotCreate(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	// starting points
	from := virtual.Tree{
		".": &virtual.File{Mode: fs.ModeDir},
	}
	to := virtual.Tree{}
	err := virtual.Sync(log, from, to, ".")
	is.NoErr(err)
	is.Equal(len(to), 0)
}

func TestSyncAvoidDotUpdate(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	// starting points
	from := virtual.Tree{
		".": &virtual.File{Mode: fs.ModeDir},
	}
	to := virtual.Tree{
		".": &virtual.File{Mode: fs.ModeDir | 0755},
	}
	err := virtual.Sync(log, from, to, ".")
	is.NoErr(err)
	is.Equal(len(to), 1)
}

// Avoid deleting the root of the target fs
func TestSyncAvoidDotDelete(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	// starting points
	from := virtual.Tree{}
	to := virtual.Tree{
		".": &virtual.File{Mode: fs.ModeDir},
	}
	err := virtual.Sync(log, from, to, ".")
	is.NoErr(err)
	// . should be ignored
	is.Equal(len(to), 1)
}

// Relative path from source path
func TestSyncRelativeSource(t *testing.T) {
	is := is.New(t)
	// starting points
	from := virtual.Tree{
		"bud/.cli/main.go": &virtual.File{Data: []byte("package main")},
		"bud/.cli/a/a.go":  &virtual.File{Data: []byte("package a")},
	}
	to := virtual.Tree{
		"a/a.go": &virtual.File{Data: []byte("package aa")},
	}
	err := dsync.Dir(from, "bud/.cli", to, ".")
	is.NoErr(err)
	_, ok := to["main.go"]
	is.True(ok) // missing main.go
	data, err := fs.ReadFile(to, "a/a.go")
	is.NoErr(err)
	is.Equal(string(data), "package a")
}

func TestSyncDeleteNotExist(t *testing.T) {
	is := is.New(t)
	log := testlog.New()

	// starting points
	from := genfs.New(dag.Discard, virtual.List{}, log)
	notExist := false
	from.GenerateFile("bud/generate/main.go", func(fsys genfs.FS, file *genfs.File) error {
		if notExist {
			return fs.ErrNotExist
		}
		file.Data = []byte("package main")
		return nil
	})
	to := virtual.OS(t.TempDir())

	// sync
	err := virtual.Sync(log, from, to, ".")
	is.NoErr(err)
	data, err := fs.ReadFile(to, "bud/generate/main.go")
	is.NoErr(err)
	is.Equal(string(data), "package main")

	// set not exist and sync again
	notExist = true
	err = virtual.Sync(log, from, to, ".")
	is.NoErr(err)
	data, err = fs.ReadFile(to, "bud/generate/main.go")
	is.True(errors.Is(err, fs.ErrNotExist))
	is.Equal(data, nil)
}
