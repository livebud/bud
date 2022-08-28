package virtual_test

import (
	"fmt"
	"io"
	"io/fs"
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/virtual"
)

func TestFile(t *testing.T) {
	is := is.New(t)
	expect := &virtual.File{
		Path: "a/b.txt",
		Data: []byte("c"),
	}
	result, err := virtual.MarshalJSON(expect)
	is.NoErr(err)
	actual, err := virtual.UnmarshalJSON(result)
	is.NoErr(err)
	stat, err := actual.Stat()
	is.NoErr(err)
	is.Equal(stat.Name(), "b.txt")
	is.Equal(stat.Size(), int64(-1))
	is.Equal(stat.Mode(), fs.FileMode(0))
	is.True(stat.ModTime().IsZero())
	is.Equal(stat.IsDir(), false)
	is.Equal(stat.Sys(), nil)
	code, err := io.ReadAll(actual)
	is.NoErr(err)
	is.Equal(string(code), "c")
}

func TestDir(t *testing.T) {
	is := is.New(t)
	expect := &virtual.Dir{
		Path: "a/b",
		Entries: []fs.DirEntry{
			&virtual.DirEntry{
				Path: "c.txt",
			},
			&virtual.DirEntry{
				Path: "d.txt",
			},
		},
	}
	result, err := virtual.MarshalJSON(expect)
	is.NoErr(err)
	actual, err := virtual.UnmarshalJSON(result)
	is.NoErr(err)
	stat, err := actual.Stat()
	is.NoErr(err)
	is.Equal(stat.Name(), "b")
	is.Equal(stat.Size(), int64(-1))
	is.Equal(stat.Mode(), fs.FileMode(fs.ModeDir))
	is.True(stat.ModTime().IsZero())
	is.Equal(stat.IsDir(), true)
	is.Equal(stat.Sys(), nil)
	dir, ok := actual.(fs.ReadDirFile)
	is.True(ok)
	entries, err := dir.ReadDir(-1)
	is.NoErr(err)
	is.Equal(len(entries), 2)
	is.Equal(entries[0].Name(), "c.txt")
	is.Equal(entries[0].IsDir(), false)
	stat, err = entries[0].Info()
	is.NoErr(err)
	is.Equal(stat.Size(), int64(-1))
	is.Equal(stat.Mode(), fs.FileMode(0))
	is.True(stat.ModTime().IsZero())
	is.Equal(stat.Sys(), nil)
	is.Equal(entries[1].Name(), "d.txt")
	is.Equal(entries[1].IsDir(), false)
	stat, err = entries[1].Info()
	is.NoErr(err)
	is.Equal(stat.Size(), int64(-1))
	is.Equal(stat.Mode(), fs.FileMode(0))
	is.True(stat.ModTime().IsZero())
	is.Equal(stat.Sys(), nil)
}

func TestDirNoEntries(t *testing.T) {
	is := is.New(t)
	expect := &virtual.Dir{
		Path: "a/b",
	}
	result, err := virtual.MarshalJSON(expect)
	is.NoErr(err)
	fmt.Println(string(result))
	actual, err := virtual.UnmarshalJSON(result)
	is.NoErr(err)
	stat, err := actual.Stat()
	is.NoErr(err)
	is.Equal(stat.Name(), "b")
	is.Equal(stat.Size(), int64(-1))
	is.Equal(stat.Mode(), fs.FileMode(fs.ModeDir))
	is.True(stat.ModTime().IsZero())
	is.Equal(stat.IsDir(), true)
	is.Equal(stat.Sys(), nil)
	dir, ok := actual.(fs.ReadDirFile)
	is.True(ok)
	entries, err := dir.ReadDir(-1)
	is.NoErr(err)
	is.Equal(len(entries), 0)
}

func TestDirEmptyEntries(t *testing.T) {
	is := is.New(t)
	expect := &virtual.Dir{
		Path:    "a/b",
		Entries: []fs.DirEntry{},
	}
	result, err := virtual.MarshalJSON(expect)
	is.NoErr(err)
	fmt.Println(string(result))
	actual, err := virtual.UnmarshalJSON(result)
	is.NoErr(err)
	stat, err := actual.Stat()
	is.NoErr(err)
	is.Equal(stat.Name(), "b")
	is.Equal(stat.Size(), int64(-1))
	is.Equal(stat.Mode(), fs.FileMode(fs.ModeDir))
	is.True(stat.ModTime().IsZero())
	is.Equal(stat.IsDir(), true)
	is.Equal(stat.Sys(), nil)
	dir, ok := actual.(fs.ReadDirFile)
	is.True(ok)
	entries, err := dir.ReadDir(-1)
	is.NoErr(err)
	is.Equal(len(entries), 0)
}
