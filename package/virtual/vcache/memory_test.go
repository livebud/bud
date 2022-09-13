package vcache_test

import (
	"io/fs"
	"testing"

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
	file := entry.Open()
	info, err := file.Stat()
	is.NoErr(err)
	is.Equal(info.Mode().String(), "-rw-r--r--")
	de := fs.FileInfoToDirEntry(info)
	is.Equal(de.Type().String(), "-rw-r--r--")
	is.Equal(info.Mode().Type().String(), "-rw-r--r--")
}
