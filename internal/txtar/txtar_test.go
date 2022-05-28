package txtar_test

import (
	"io/fs"
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/txtar"
)

func TestParseFile(t *testing.T) {
	is := is.New(t)
	fsys, err := txtar.ParseFile("testdata/one.txt")
	is.NoErr(err)
	des, err := fs.ReadDir(fsys, ".")
	is.NoErr(err)
	is.Equal(len(des), 3)
	is.Equal(des[0].Name(), "a.go")
	is.Equal(des[1].Name(), "b")
	is.Equal(des[2].Name(), "c")
	code, err := fs.ReadFile(fsys, "a.go")
	is.NoErr(err)
	is.Equal(string(code), "package a\n\n")
	code, err = fs.ReadFile(fsys, "b/b.go")
	is.NoErr(err)
	is.Equal(string(code), "package b\n\n")
	code, err = fs.ReadFile(fsys, "c/c/c.txt")
	is.NoErr(err)
	is.Equal(string(code), "c\n")
}
