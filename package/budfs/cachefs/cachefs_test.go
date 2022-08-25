package cachefs_test

import (
	"io"
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/virtual"
	"github.com/livebud/bud/package/budfs/cachefs"
	"github.com/livebud/bud/package/budfs/mergefs"
	"github.com/livebud/bud/package/log/testlog"
)

func TestCache(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	cache := cachefs.New(log)
	vfile := &virtual.File{
		Data: []byte("a"),
	}
	cache.Set("a.txt", vfile)
	entry, ok := cache.Get("a.txt")
	is.True(ok)
	code, err := io.ReadAll(entry.Open())
	is.NoErr(err)
	is.Equal(string(code), "a")
	// Delete from the cache
	cache.Delete("a.txt")
	// Get from the cache
	entry, ok = cache.Get("a.txt")
	is.Equal(ok, false)
	is.Equal(entry, nil)
}

func TestReadDirFile(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	afs := &fstest.MapFS{
		"a.txt": &fstest.MapFile{Data: []byte("a")},
	}
	bfs := &fstest.MapFS{
		"b.txt": &fstest.MapFile{Data: []byte("b")},
	}
	c1 := cachefs.New(log)
	c2 := cachefs.New(log)
	acfs := c1.Wrap(afs)
	bcfs := c2.Wrap(bfs)
	merge := mergefs.Merge(acfs, bcfs)
	file, err := merge.Open(".")
	is.NoErr(err)
	dir, ok := file.(fs.ReadDirFile)
	is.True(ok)
	des, err := dir.ReadDir(-1)
	is.NoErr(err)
	is.Equal(len(des), 2)
}

func TestSize(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	fsys := fstest.MapFS{
		"a.txt": &fstest.MapFile{Data: []byte("a")},
	}
	cache := cachefs.New(log)
	cfs := cache.Wrap(fsys)
	stat, err := fs.Stat(cfs, "a.txt")
	is.NoErr(err)
	des, err := fs.ReadDir(cfs, ".")
	is.NoErr(err)
	is.Equal(len(des), 1)
	info, err := des[0].Info()
	is.NoErr(err)
	is.Equal(stat.Size(), info.Size())
}

func TestTransparent(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	afs := &fstest.MapFS{
		"a.txt": &fstest.MapFile{Data: []byte("a")},
	}
	bfs := &fstest.MapFS{
		"b.txt": &fstest.MapFile{Data: []byte("b")},
	}
	merge := mergefs.Merge(afs, bfs)
	is.NoErr(fstest.TestFS(merge, "a.txt", "b.txt"))
	c1 := cachefs.New(log)
	c2 := cachefs.New(log)
	acfs := c1.Wrap(afs)
	bcfs := c2.Wrap(bfs)
	cmerge := mergefs.Merge(acfs, bcfs)
	is.NoErr(fstest.TestFS(cmerge, "a.txt", "b.txt"))
}
