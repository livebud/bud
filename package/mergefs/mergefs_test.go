package mergefs_test

import (
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/matryer/is"

	"gitlab.com/mnm/bud/package/mergefs"
)

func TestMerge(t *testing.T) {
	is := is.New(t)
	a := fstest.MapFS{
		"a.txt": &fstest.MapFile{Data: []byte("a")},
	}
	b := fstest.MapFS{
		"b.txt": &fstest.MapFile{Data: []byte("b")},
	}
	c := fstest.MapFS{
		"c.txt": &fstest.MapFile{Data: []byte("c")},
	}
	fsys := mergefs.Merge(a, b, c)
	des, err := fs.ReadDir(fsys, ".")
	is.NoErr(err)
	is.Equal(len(des), 3)
	is.Equal(des[0].Name(), "a.txt")
	is.Equal(des[1].Name(), "b.txt")
	is.Equal(des[2].Name(), "c.txt")
}

func TestInnerMerge(t *testing.T) {
	is := is.New(t)
	a := fstest.MapFS{
		"a.txt": &fstest.MapFile{Data: []byte("a")},
	}
	b := fstest.MapFS{
		"d/b.txt": &fstest.MapFile{Data: []byte("b")},
	}
	c := fstest.MapFS{
		"d/c.txt": &fstest.MapFile{Data: []byte("c")},
	}
	fsys := mergefs.Merge(a, b, c)
	des, err := fs.ReadDir(fsys, ".")
	is.NoErr(err)
	is.Equal(len(des), 2)
	is.Equal(des[0].Name(), "a.txt")
	is.Equal(des[1].Name(), "d")
	des, err = fs.ReadDir(fsys, "d")
	is.NoErr(err)
	is.Equal(len(des), 2)
	is.Equal(des[0].Name(), "b.txt")
	is.Equal(des[1].Name(), "c.txt")
}

func TestOverride(t *testing.T) {
	is := is.New(t)
	a := fstest.MapFS{
		"a.txt": &fstest.MapFile{Data: []byte("a")},
	}
	b := fstest.MapFS{
		"a.txt": &fstest.MapFile{Data: []byte("b")},
	}
	fsys := mergefs.Merge(a, b)
	des, err := fs.ReadDir(fsys, ".")
	is.NoErr(err)
	is.Equal(len(des), 1)
	is.Equal(des[0].Name(), "a.txt")
	code, err := fs.ReadFile(fsys, "a.txt")
	is.NoErr(err)
	is.Equal(string(code), "a")
}
