package fscache_test

import (
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/livebud/bud/internal/is"

	"github.com/livebud/bud/internal/fscache"
	"github.com/livebud/bud/package/log/testlog"
	"github.com/livebud/bud/package/merged"
)

func TestReadDirFile(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	afs := &fstest.MapFS{
		"a.txt": &fstest.MapFile{Data: []byte("a")},
	}
	bfs := &fstest.MapFS{
		"b.txt": &fstest.MapFile{Data: []byte("b")},
	}
	acfs := fscache.Wrap(afs, log, "a")
	bcfs := fscache.Wrap(bfs, log, "b")
	merge := merged.Merge(acfs, bcfs)
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
	cfs := fscache.Wrap(fsys, log, "a")
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
	merge := merged.Merge(afs, bfs)
	is.NoErr(fstest.TestFS(merge, "a.txt", "b.txt"))
	acfs := fscache.Wrap(afs, log, "a")
	bcfs := fscache.Wrap(bfs, log, "b")
	cmerge := merged.Merge(acfs, bcfs)
	is.NoErr(fstest.TestFS(cmerge, "a.txt", "b.txt"))
}
