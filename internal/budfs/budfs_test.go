package budfs_test

import (
	"strconv"
	"testing"
	"testing/fstest"

	"github.com/livebud/bud/internal/budfs"

	"io/fs"

	"github.com/livebud/bud/package/log/testlog"

	"github.com/livebud/bud/internal/is"
)

func TestReadFsys(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	fsys := fstest.MapFS{
		"a.txt": &fstest.MapFile{Data: []byte("a")},
	}
	bfs := budfs.New(fsys, log)
	des, err := fs.ReadDir(bfs, ".")
	is.NoErr(err)
	is.Equal(len(des), 1)
}

func TestGeneratorPriority(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	fsys := fstest.MapFS{
		"a.txt": &fstest.MapFile{Data: []byte("a")},
	}
	bfs := budfs.New(fsys, log)
	bfs.FileGenerator("a.txt", budfs.GenerateFile(func(fsys budfs.FS, file *budfs.File) error {
		file.Data = []byte("b")
		return nil
	}))
	code, err := fs.ReadFile(bfs, "a.txt")
	is.NoErr(err)
	is.Equal(string(code), "b")
}

func TestCaching(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	fsys := fstest.MapFS{}
	bfs := budfs.New(fsys, log)
	count := 1
	bfs.FileGenerator("a.txt", budfs.GenerateFile(func(fsys budfs.FS, file *budfs.File) error {
		file.Data = []byte(strconv.Itoa(count))
		count++
		return nil
	}))
	code, err := fs.ReadFile(bfs, "a.txt")
	is.NoErr(err)
	is.Equal(string(code), "1")
	code, err = fs.ReadFile(bfs, "a.txt")
	is.NoErr(err)
	is.Equal(string(code), "1")
}

func TestLinking(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	fsys := fstest.MapFS{}
	bfs := budfs.New(fsys, log)
	count := 1
	bfs.FileGenerator("a.txt", budfs.GenerateFile(func(fsys budfs.FS, file *budfs.File) error {
		file.Data = []byte(strconv.Itoa(count))
		count++
		return nil
	}))
	bfs.FileGenerator("b.txt", budfs.GenerateFile(func(fsys budfs.FS, file *budfs.File) error {
		code, err := fs.ReadFile(fsys, "a.txt")
		if err != nil {
			return err
		}
		file.Data = []byte(code)
		return nil
	}))
	code, err := fs.ReadFile(bfs, "b.txt")
	is.NoErr(err)
	is.Equal(string(code), "1")
	// Cached
	code, err = fs.ReadFile(bfs, "b.txt")
	is.NoErr(err)
	is.Equal(string(code), "1")
	bfs.Update("a.txt")
	code, err = fs.ReadFile(bfs, "b.txt")
	is.NoErr(err)
	is.Equal(string(code), "2")
}
