package modcache_test

import (
	"fmt"
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/internal/modcache"
)

func TestWriteModule(t *testing.T) {
	is := is.New(t)
	cacheDir := t.TempDir()
	modCache := modcache.New(cacheDir)
	err := modCache.Write(modcache.Versions{
		"v0.0.1": modcache.Files{
			"go.mod":   "module mod.test/one\n\ngo 1.12",
			"const.go": "package module\n\nconst Answer = 42",
		},
		"v0.0.2": modcache.Files{
			"go.mod":   "module mod.test/one\n\ngo 1.12",
			"const.go": "package module\n\nconst Answer = 43",
		},
	})
	is.NoErr(err)
	dir, err := modCache.ResolveDirectory("mod.test/one", "v0.0.2")
	is.NoErr(err)
	fmt.Println(dir)
	is.Equal(dir, modCache.Directory("mod.test", "one@v0.0.2"))
}

func TestResolveDirectoryFromCache(t *testing.T) {
	is := is.New(t)
	modCache := modcache.Default()
	is.True(modCache.Directory() != "")
	dir, err := modCache.ResolveDirectory("github.com/matryer/is", "v1.4.0")
	is.NoErr(err)
	is.Equal(dir, modCache.Directory("github.com", "matryer", "is@v1.4.0"))
}
