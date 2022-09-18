package virtual_test

import (
	"io"
	"io/fs"
	"testing"
	"time"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/virtual"
)

var now = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

func TestFile(t *testing.T) {
	is := is.New(t)
	expect := &virtual.File{
		Path:    "a/b.txt",
		Data:    []byte("c"),
		ModTime: now,
		Mode:    0644,
	}
	result, err := virtual.MarshalJSON(virtual.New(expect))
	is.NoErr(err)
	actual, err := virtual.UnmarshalJSON(result)
	is.NoErr(err)
	stat, err := actual.Stat()
	is.NoErr(err)
	is.Equal(stat.Name(), "b.txt")
	is.Equal(stat.Size(), int64(1))
	is.Equal(stat.Mode(), fs.FileMode(0644))
	is.Equal(stat.ModTime(), now)
	is.Equal(stat.IsDir(), false)
	is.Equal(stat.Sys(), nil)
	code, err := io.ReadAll(actual)
	is.NoErr(err)
	is.Equal(string(code), "c")
}

func TestDir(t *testing.T) {
	is := is.New(t)
	expect := &virtual.Dir{
		Path:    "a/b.txt",
		ModTime: now,
		Mode:    0755,
		Entries: []fs.DirEntry{
			&virtual.DirEntry{
				Path:    "c.txt",
				ModTime: now,
				Mode:    0644,
				Size:    10,
			},
			&virtual.DirEntry{
				Path:    "d.txt",
				ModTime: now,
				Mode:    0644,
				Size:    20,
			},
		},
	}
	result, err := virtual.MarshalJSON(virtual.New(expect))
	is.NoErr(err)
	actual, err := virtual.UnmarshalJSON(result)
	is.NoErr(err)
	stat, err := actual.Stat()
	is.NoErr(err)
	is.Equal(stat.Name(), "b.txt")
	is.Equal(stat.Size(), int64(0))
	is.Equal(stat.Mode(), fs.FileMode(0755|fs.ModeDir))
	is.Equal(stat.ModTime(), now)
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
	is.Equal(stat.Size(), int64(10))
	is.Equal(stat.Mode(), fs.FileMode(0644))
	is.Equal(stat.ModTime(), now)
	is.Equal(stat.Sys(), nil)
	is.Equal(entries[1].Name(), "d.txt")
	is.Equal(entries[1].IsDir(), false)
	stat, err = entries[1].Info()
	is.NoErr(err)
	is.Equal(stat.Size(), int64(20))
	is.Equal(stat.Mode(), fs.FileMode(0644))
	is.Equal(stat.ModTime(), now)
	is.Equal(stat.Sys(), nil)
}
