package vfs_test

import (
	"errors"
	"io/fs"
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/vfs"
)

func TestMemory(t *testing.T) {
	is := is.New(t)
	fsys := vfs.Memory{
		"duo/view/index.svelte": &vfs.File{
			Data: []byte(`<h1>index</h1>`),
		},
	}

	// Read duo/view/index.svelte
	code, err := fs.ReadFile(fsys, "duo/view/index.svelte")
	is.NoErr(err)
	is.Equal(string(code), `<h1>index</h1>`)

	// stat duo/
	stat, err := fs.Stat(fsys, "duo")
	is.NoErr(err)
	is.Equal(stat.Name(), "duo")
	is.Equal(stat.IsDir(), true)
	is.Equal(stat.Mode(), fs.FileMode(fs.ModeDir))

	// mkdir duo/controller
	err = fsys.MkdirAll("duo/controller", 0755)
	is.NoErr(err)
	stat, err = fs.Stat(fsys, "duo/controller")
	is.NoErr(err)
	is.Equal(stat.Name(), "controller")
	is.Equal(stat.IsDir(), true)
	is.Equal(stat.Mode(), fs.FileMode(0755|fs.ModeDir))

	// write duo/controller/controller.go
	err = fsys.WriteFile("duo/controller/controller.go", []byte(`package controller`), 0644)
	is.NoErr(err)

	// read duo/controller/controller.go
	code, err = fs.ReadFile(fsys, "duo/controller/controller.go")
	is.NoErr(err)
	is.Equal(string(code), `package controller`)

	// remove duo/view
	err = fsys.RemoveAll("duo/view")
	is.NoErr(err)

	// Read duo/view/index.svelte
	code, err = fs.ReadFile(fsys, "duo/view/index.svelte")
	is.Equal(errors.Is(err, fs.ErrNotExist), true)
	is.Equal(code, nil)
}

func TestWriteAll(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	dirfs := vfs.OS(dir)
	fsys := vfs.Memory{
		"duo/view/index.svelte": &vfs.File{
			Data: []byte(`<h1>index</h1>`),
		},
	}

	// write duo/controller/controller.go
	err := fsys.WriteFile("duo/controller/controller.go", []byte(`package controller`), 0644)
	is.NoErr(err)

	// remove duo/view
	err = fsys.RemoveAll("duo/view")
	is.NoErr(err)

	err = vfs.WriteAll(".", dir, fsys)
	is.NoErr(err)

	// _tmp/duo has real entries
	des, err := fs.ReadDir(dirfs, "duo")
	is.NoErr(err)
	is.Equal(len(des), 1)
	is.Equal(des[0].Name(), "controller")

	// duo/controller/controller.go exists in _tmp
	code, err := fs.ReadFile(dirfs, "duo/controller/controller.go")
	is.NoErr(err)
	is.Equal(string(code), `package controller`)

	// duo/view/index.svelte doesn't exist in _tmp
	code, err = fs.ReadFile(dirfs, "duo/view/index.svelte")
	is.Equal(errors.Is(err, fs.ErrNotExist), true)
	is.Equal(code, nil)
}

func TestWriteRead(t *testing.T) {
	is := is.New(t)
	fsys := vfs.Memory{}
	err := fsys.MkdirAll("a/b", 0755)
	is.NoErr(err)
	err = fsys.WriteFile("a/b/c.txt", []byte("c"), 0644)
	is.NoErr(err)
	stat, err := fs.Stat(fsys, "a/b")
	is.NoErr(err)
	is.Equal(stat.Name(), "b")
	is.Equal(stat.IsDir(), true)
}
