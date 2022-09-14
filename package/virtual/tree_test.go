package virtual_test

import (
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/virtual"
)

func TestFSTree(t *testing.T) {
	is := is.New(t)
	fsys := virtual.Tree{
		"bud/view/index.svelte": &virtual.File{
			Data: []byte(`<h1>index</h1>`),
		},
		"bud/view/about/index.svelte": &virtual.File{
			Data: []byte(`<h1>about</h1>`),
		},
	}
	err := fstest.TestFS(fsys, "bud/view/index.svelte", "bud/view/about/index.svelte")
	is.NoErr(err)
}

func TestReadDir(t *testing.T) {
	is := is.New(t)
	fsys := virtual.Tree{
		"bud/view/index.svelte": &virtual.File{
			Data: []byte(`<h1>index</h1>`),
		},
		"bud/view/about/index.svelte": &virtual.File{
			Data: []byte(`<h1>about</h1>`),
		},
	}
	des, err := fs.ReadDir(fsys, "bud/view")
	is.NoErr(err)
	is.Equal(len(des), 2)
	is.Equal(des[0].Name(), "about")
	is.Equal(des[1].Name(), "index.svelte")
}

func TestTreeWriteRead(t *testing.T) {
	is := is.New(t)
	fsys := virtual.Tree{}

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

func TestTreeRoot(t *testing.T) {
	is := is.New(t)
	fsys := virtual.Tree{}
	des, err := fs.ReadDir(fsys, ".")
	is.NoErr(err)
	is.Equal(len(des), 0)
}
