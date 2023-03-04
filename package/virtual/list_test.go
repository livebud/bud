package virtual_test

import (
	"errors"
	"io/fs"
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/virtual"
)

func TestList(t *testing.T) {
	is := is.New(t)
	fsys := virtual.List{
		&virtual.File{
			Path: "bud/view/index.svelte",
			Data: []byte(`<h1>index</h1>`),
		},
		&virtual.File{
			Path: "bud/view/about/index.svelte",
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

func TestListWriteRead(t *testing.T) {
	is := is.New(t)
	fsys := virtual.List{}

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

func TestListRoot(t *testing.T) {
	is := is.New(t)
	fsys := virtual.List{}
	des, err := fs.ReadDir(fsys, ".")
	is.True(errors.Is(err, fs.ErrNotExist))
	is.Equal(des, nil)
	fsys = virtual.List{
		&virtual.File{
			Path: ".",
			Mode: fs.ModeDir,
		},
	}
	des, err = fs.ReadDir(fsys, ".")
	is.NoErr(err)
	is.Equal(len(des), 0)
}

func TestListReadDirWithChild(t *testing.T) {
	is := is.New(t)
	fsys := virtual.List{
		&virtual.File{
			Path: "a/b.txt",
			Data: []byte("b"),
		},
		&virtual.File{
			Path: "a",
			Mode: fs.ModeDir,
		},
	}
	des, err := fs.ReadDir(fsys, "a")
	is.NoErr(err)
	is.Equal(len(des), 1)
	is.Equal(des[0].Name(), "b.txt")
}

func TestListOpenParentInvalid(t *testing.T) {
	is := is.New(t)
	fsys := virtual.List{
		&virtual.File{
			Path: "a/b.txt",
			Data: []byte("b"),
		},
		&virtual.File{
			Path: "a",
			Mode: fs.ModeDir,
		},
	}
	file, err := fsys.Open("../a/b.txt")
	is.True(errors.Is(err, fs.ErrInvalid))
	is.Equal(file, nil)
}

func TestEmptyPath(t *testing.T) {
	is := is.New(t)
	fsys := virtual.List{
		&virtual.File{
			Data: []byte(`<h1>index</h1>`),
		},
	}
	code, err := fs.ReadFile(fsys, "")
	is.True(err != nil)
	is.True(errors.Is(err, fs.ErrInvalid))
	is.Equal(code, nil)
	code, err = fs.ReadFile(fsys, ".")
	is.True(err != nil)
	is.True(errors.Is(err, fs.ErrNotExist))
	is.Equal(code, nil)
}
