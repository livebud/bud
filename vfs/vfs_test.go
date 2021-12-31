package vfs_test

import (
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/vfs"
)

func TestSomeExistOne(t *testing.T) {
	is := is.New(t)
	fsys := vfs.Map{
		"a.go": "",
		"b.go": "",
		"c.go": "",
	}
	exist := vfs.SomeExist(fsys, "a.go", "d.go")
	is.Equal(len(exist), 1)
	is.True(exist["a.go"])
	is.True(!exist["d.go"])
}

func TestSomeExistTwo(t *testing.T) {
	is := is.New(t)
	fsys := vfs.Map{
		"a.go": "",
		"b.go": "",
		"c.go": "",
	}
	exist := vfs.SomeExist(fsys, "a.go", "b.go", "d.go")
	is.Equal(len(exist), 2)
	is.True(exist["a.go"])
	is.True(exist["b.go"])
	is.True(!exist["d.go"])
}

func TestSomeExistZero(t *testing.T) {
	is := is.New(t)
	fsys := vfs.Map{
		"a.go": "",
		"b.go": "",
		"c.go": "",
	}
	exist := vfs.SomeExist(fsys, "d.go", "e.go")
	is.Equal(len(exist), 0)
	is.True(!exist["d.go"])
	is.True(!exist["e.go"])
}
