package virt_test

import (
	"errors"
	"io/fs"
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/virt"
)

func TestMap(t *testing.T) {
	is := is.New(t)
	fsys := virt.Map{
		"bud/view/index.svelte": &virt.File{
			Data: []byte(`<h1>index</h1>`),
		},
		"bud/view/about/index.svelte": &virt.File{
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
	fsys := virt.Map{}

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

func TestMapRoot(t *testing.T) {
	is := is.New(t)
	fsys := virt.Map{}
	des, err := fs.ReadDir(fsys, ".")
	is.True(errors.Is(err, fs.ErrNotExist))
	is.Equal(des, nil)
	fsys = virt.Map{
		".": &virt.File{Mode: fs.ModeDir},
	}
	des, err = fs.ReadDir(fsys, ".")
	is.NoErr(err)
	is.Equal(len(des), 0)
}

func TestMapReadDirWithChild(t *testing.T) {
	is := is.New(t)
	fsys := virt.Map{
		"a/b.txt": &virt.File{Data: []byte("b")},
		"a":       &virt.File{Mode: fs.ModeDir},
	}
	des, err := fs.ReadDir(fsys, "a")
	is.NoErr(err)
	is.Equal(len(des), 1)
	is.Equal(des[0].Name(), "b.txt")
}

func TestMapOpenParentInvalid(t *testing.T) {
	is := is.New(t)
	fsys := virt.Map{
		"a/b.txt": &virt.File{Data: []byte("b")},
		"a":       &virt.File{Mode: fs.ModeDir},
	}
	file, err := fsys.Open("../a/b.txt")
	is.True(errors.Is(err, fs.ErrInvalid))
	is.Equal(file, nil)
}
