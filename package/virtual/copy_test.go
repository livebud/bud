package virtual_test

import (
	"errors"
	"io/fs"
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/log/testlog"
	"github.com/livebud/bud/package/virtual"
)

func TestCopyEmpty(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	from := virtual.Tree{}
	to := virtual.Tree{}
	is.NoErr(virtual.Copy(log, from, to))
	is.Equal(len(from), 0)
	is.Equal(len(to), 0)
}

func TestCopyEmptySub(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	from := virtual.Tree{}
	to := virtual.Tree{}
	err := virtual.Copy(log, from, to, "sub", "dir")
	is.True(err != nil)
	is.True(errors.Is(err, fs.ErrNotExist))
	is.Equal(len(from), 0)
	is.Equal(len(to), 0)
}

func TestCopyFromTo(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	from := virtual.Tree{
		"sub/dir/a.txt": &virtual.File{Data: []byte("a")},
		"sub/dir/b.txt": &virtual.File{Data: []byte("b")},
		"sub/c.txt":     &virtual.File{Data: []byte("c")},
		"d.txt":         &virtual.File{Data: []byte("d")},
	}
	to := virtual.Tree{
		"e.txt": &virtual.File{Data: []byte("e")},
	}
	err := virtual.Copy(log, from, to)
	is.NoErr(err)
	is.Equal(len(from), 4)
	is.Equal(len(to), 7)
	code, err := fs.ReadFile(to, "sub/dir/a.txt")
	is.NoErr(err)
	is.Equal(string(code), "a")
	code, err = fs.ReadFile(to, "sub/dir/b.txt")
	is.NoErr(err)
	is.Equal(string(code), "b")
	code, err = fs.ReadFile(to, "sub/c.txt")
	is.NoErr(err)
	is.Equal(string(code), "c")
	code, err = fs.ReadFile(to, "d.txt")
	is.NoErr(err)
	is.Equal(string(code), "d")
	code, err = fs.ReadFile(to, "e.txt")
	is.NoErr(err)
	is.Equal(string(code), "e")
}

func TestCopyFromToSub(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	from := virtual.Tree{
		"sub/dir/a.txt": &virtual.File{Data: []byte("a")},
		"sub/dir/b.txt": &virtual.File{Data: []byte("b")},
		"sub/c.txt":     &virtual.File{Data: []byte("c")},
		"d.txt":         &virtual.File{Data: []byte("d")},
	}
	to := virtual.Tree{
		"sub/e.txt": &virtual.File{Data: []byte("e")},
	}
	err := virtual.Copy(log, from, to, "sub")
	is.NoErr(err)
	is.Equal(len(from), 4)
	is.Equal(len(to), 5)
	code, err := fs.ReadFile(to, "sub/dir/a.txt")
	is.NoErr(err)
	is.Equal(string(code), "a")
	code, err = fs.ReadFile(to, "sub/dir/b.txt")
	is.NoErr(err)
	is.Equal(string(code), "b")
	code, err = fs.ReadFile(to, "sub/c.txt")
	is.NoErr(err)
	is.Equal(string(code), "c")
	code, err = fs.ReadFile(to, "sub/e.txt")
	is.NoErr(err)
	is.Equal(string(code), "e")
}
