package vfs_test

import (
	"errors"
	"io/fs"
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/virtual/vfs"
)

func TestMap(t *testing.T) {
	is := is.New(t)
	fsys := vfs.Map{
		"bud/view/index.svelte": &vfs.Entry{
			Data: []byte(`<h1>index</h1>`),
		},
		"bud/view/about/index.svelte": &vfs.Entry{
			Data: []byte(`<h1>about</h1>`),
		},
	}

	// Read bud/view/index.svelte
	code, err := fs.ReadFile(fsys, "bud/view/index.svelte")
	is.NoErr(err)
	is.Equal(string(code), `<h1>index</h1>`)

	// Read bud/view/about/index.svelte
	code, err = fs.ReadFile(fsys, "bud/view/about/index.svelte")
	is.NoErr(err)
	is.Equal(string(code), `<h1>about</h1>`)

	// Remove bud/view/index.svelte
	err = fsys.RemoveAll("bud/view/index.svelte")
	is.NoErr(err)

	// Read bud/view/index.svelte
	code, err = fs.ReadFile(fsys, "bud/view/index.svelte")
	is.Equal(errors.Is(err, fs.ErrNotExist), true)
	is.Equal(code, nil)
}

func TestMapWriteRead(t *testing.T) {
	is := is.New(t)
	fsys := vfs.Map{}

	// Create a directory
	err := fsys.MkdirAll("a/b", 0755)
	is.NoErr(err)
	stat, err := fs.Stat(fsys, "a/b")
	is.NoErr(err)
	is.Equal(stat.Name(), "b")
	is.Equal(stat.IsDir(), true, "a/b should be a directory")

	// Write a file
	err = fsys.WriteFile("a/b/c.txt", []byte("c"), 0644)
	is.NoErr(err)
	code, err := fs.ReadFile(fsys, "a/b/c.txt")
	is.NoErr(err)
	is.Equal(string(code), `c`)
}
