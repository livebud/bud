package budfs_test

import (
	"testing"
	"testing/fstest"

	"github.com/livebud/bud/internal/budfs"

	"io/fs"

	"github.com/livebud/bud/package/log/testlog"

	"github.com/livebud/bud/internal/is"
)

func TestReadRoot(t *testing.T) {
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
