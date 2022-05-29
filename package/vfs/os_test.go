package vfs_test

import (
	"errors"
	"io/fs"
	"os"
	"testing"
	"testing/fstest"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/vfs"
)

func TestOS(t *testing.T) {
	is := is.New(t)

	// Prepare
	err := os.RemoveAll("_tmp")
	is.NoErr(err)
	defer func() {
		if !t.Failed() {
			err := os.RemoveAll("_tmp")
			is.NoErr(err)
		}
	}()

	// Write initial files
	mem := vfs.Memory{
		"duo/view/index.svelte": &fstest.MapFile{
			Data: []byte(`<h1>index</h1>`),
			Mode: 0644,
		},
	}
	err = vfs.WriteAll(".", "_tmp", mem)
	is.NoErr(err)

	// Initialize
	fsys := vfs.OS("_tmp")

	// Read duo/view/index.svelte
	code, err := fs.ReadFile(fsys, "duo/view/index.svelte")
	is.NoErr(err)
	is.Equal(string(code), `<h1>index</h1>`)

	// stat duo/
	stat, err := fs.Stat(fsys, "duo")
	is.NoErr(err)
	is.Equal(stat.Name(), "duo")
	is.Equal(stat.IsDir(), true)
	is.Equal(stat.Mode(), fs.FileMode(0755|fs.ModeDir))

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
