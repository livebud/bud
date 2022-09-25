package vcache_test

import (
	"io/fs"
	"testing"

	"github.com/livebud/bud/package/budfs/treefs"
	"github.com/livebud/bud/package/log/testlog"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/virtual"

	"github.com/livebud/bud/package/virtual/vcache"
)

func TestStat(t *testing.T) {
	is := is.New(t)
	cache := vcache.New()
	cache.Set("go.mod", &virtual.File{
		Path: "go.mod",
		Data: []byte("module github.com/livebud/bud"),
		Mode: 0644,
	})
	entry, ok := cache.Get("go.mod")
	is.True(ok)
	file := virtual.New(entry)
	info, err := file.Stat()
	is.NoErr(err)
	is.Equal(info.Mode().String(), "-rw-r--r--")
	de := fs.FileInfoToDirEntry(info)
	is.Equal(de.Type().String(), "----------")
	is.Equal(info.Mode().Type().String(), "----------")
}

func TestReadParentNoGenerate(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	cache := vcache.New()
	tfs := treefs.New(".")
	generates := 0
	tfs.FileGenerator("controller/controller.go", treefs.Generate(func(target string) (fs.File, error) {
		generates++
		return virtual.New(&virtual.File{
			Path: target,
			Data: []byte("package controller"),
		}), nil
	}))
	fsys := vcache.Wrap(cache, tfs, log)
	des, err := fs.ReadDir(fsys, "controller")
	is.NoErr(err)
	is.Equal(len(des), 1)
	is.Equal(des[0].Name(), "controller.go")
	is.Equal(generates, 0)
	code, err := fs.ReadFile(fsys, "controller/controller.go")
	is.NoErr(err)
	is.Equal(string(code), "package controller")
	is.Equal(generates, 1)
}
