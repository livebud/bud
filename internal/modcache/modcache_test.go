package modcache_test

import (
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/internal/modcache"
)

func TestWriteModule(t *testing.T) {
	is := is.New(t)
	cacheDir := t.TempDir()
	modcache := modcache.New(cacheDir)
	err := modcache.WriteModule("v0.0.1", map[string][]byte{
		"go.mod":   []byte("module mod.test/one\n\ngo 1.12"),
		"const.go": []byte("package module\n\nconst Answer = 42"),
	})
	is.NoErr(err)
	err = modcache.WriteModule("v0.0.2", map[string][]byte{
		"go.mod":   []byte("module mod.test/one\n\ngo 1.12"),
		"const.go": []byte("package module\n\nconst Answer = 43"),
	})
	is.NoErr(err)
	dir, err := modcache.ResolveDirectory("mod.test/one", "v0.0.2")
	is.NoErr(err)
	is.Equal(dir, modcache.Directory("mod.test", "one@v0.0.2"))
}

func TestResolveDirectoryFromCache(t *testing.T) {
	is := is.New(t)
	modcache := modcache.Default()
	is.True(modcache.Directory() != "")
	dir, err := modcache.ResolveDirectory("github.com/matryer/is", "v1.4.0")
	is.NoErr(err)
	is.Equal(dir, modcache.Directory("github.com", "matryer", "is@v1.4.0"))
}
