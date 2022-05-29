package vfs_test

import (
	"errors"
	"io/fs"
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/vfs"
)

func TestExistAll(t *testing.T) {
	is := is.New(t)
	fsys := vfs.Map{
		"a.go": []byte(""),
		"b.go": []byte(""),
		"c.go": []byte(""),
	}
	err := vfs.Exist(fsys, "a.go", "c.go")
	is.NoErr(err)
}

func TestExistNotAll(t *testing.T) {
	is := is.New(t)
	fsys := vfs.Map{
		"a.go": []byte(""),
		"b.go": []byte(""),
		"c.go": []byte(""),
	}
	err := vfs.Exist(fsys, "a.go", "d.go")
	is.True(errors.Is(err, fs.ErrNotExist))
}

func TestSomeExistOne(t *testing.T) {
	is := is.New(t)
	fsys := vfs.Map{
		"a.go": []byte(""),
		"b.go": []byte(""),
		"c.go": []byte(""),
	}
	exist, err := vfs.SomeExist(fsys, "a.go", "d.go")
	is.NoErr(err)
	is.Equal(len(exist), 1)
	is.True(exist["a.go"])
	is.True(!exist["d.go"])
}

func TestSomeExistTwo(t *testing.T) {
	is := is.New(t)
	fsys := vfs.Map{
		"a.go": []byte(""),
		"b.go": []byte(""),
		"c.go": []byte(""),
	}
	exist, err := vfs.SomeExist(fsys, "a.go", "b.go", "d.go")
	is.NoErr(err)
	is.Equal(len(exist), 2)
	is.True(exist["a.go"])
	is.True(exist["b.go"])
	is.True(!exist["d.go"])
}

func TestSomeExistZero(t *testing.T) {
	is := is.New(t)
	fsys := vfs.Map{
		"a.go": []byte(""),
		"b.go": []byte(""),
		"c.go": []byte(""),
	}
	exist, err := vfs.SomeExist(fsys, "d.go", "e.go")
	is.NoErr(err)
	is.Equal(len(exist), 0)
	is.True(!exist["d.go"])
	is.True(!exist["e.go"])
}
