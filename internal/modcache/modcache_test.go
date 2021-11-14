package modcache_test

import (
	"path/filepath"
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/internal/modcache"
)

func TestWriteModule(t *testing.T) {
	is := is.New(t)
	cacheDir := t.TempDir()
	err := modcache.WriteModule(cacheDir, "v0.0.1", map[string][]byte{
		"go.mod":   []byte("module mod.test/one\n\ngo 1.12"),
		"const.go": []byte("package module\n\nconst Answer = 42"),
	})
	is.NoErr(err)
	err = modcache.WriteModule(cacheDir, "v0.0.2", map[string][]byte{
		"go.mod":   []byte("module mod.test/one\n\ngo 1.12"),
		"const.go": []byte("package module\n\nconst Answer = 43"),
	})
	is.NoErr(err)
	dir, err := modcache.ResolveDirectory(cacheDir, "mod.test/one", "v0.0.2")
	is.NoErr(err)
	is.Equal(dir, filepath.Join(cacheDir, "mod.test", "one@v0.0.2"))
}

func TestResolveDirectoryFromCache(t *testing.T) {
	is := is.New(t)
	cacheDir := modcache.Directory()
	is.True(cacheDir != "")
	dir, err := modcache.ResolveDirectory(cacheDir, "github.com/matryer/is", "v1.4.0")
	is.NoErr(err)
	is.Equal(dir, filepath.Join(cacheDir, "github.com/matryer", "is@v1.4.0"))
}
