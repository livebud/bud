package budfs_test

import (
	"io/fs"
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/budfs"
	"github.com/livebud/bud/package/log/testlog"
	"github.com/livebud/bud/package/virtual"
)

func TestGenerateFile(t *testing.T) {
	is := is.New(t)
	fsys := virtual.Map{}
	log := testlog.New()
	bfs, err := budfs.Load(fsys, log)
	bfs.GenerateFile("a.txt", func(fsys budfs.FS, file *budfs.File) error {
		file.Data = []byte("a")
		return nil
	})
	code, err := fs.ReadFile(bfs, "a.txt")
	is.NoErr(err)
	is.Equal(string(code), "a")
}
