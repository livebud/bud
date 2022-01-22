package mod_test

import (
	"errors"
	"go/build"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gitlab.com/mnm/bud/2/mod"
	"gitlab.com/mnm/bud/2/virtual"
	"gitlab.com/mnm/bud/internal/modcache"
	"gitlab.com/mnm/bud/vfs"

	"github.com/matryer/is"
)

func TestFind(t *testing.T) {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	module, err := mod.Find(wd)
	is.NoErr(err)
	dir := module.Directory()
	root := filepath.Join(wd, "..", "..")
	is.Equal(root, dir)
}

func TestFindDefault(t *testing.T) {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	module, err := mod.Find(wd)
	is.NoErr(err)
	dir := module.Directory()
	root := filepath.Join(wd, "..", "..")
	is.Equal(root, dir)
}

func TestResolveDirectory(t *testing.T) {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	modCache := modcache.Default()
	module, err := mod.Find(wd, mod.WithModCache(modCache))
	is.NoErr(err)
	dir, err := module.ResolveDirectory("github.com/matryer/is")
	is.NoErr(err)
	expected := modCache.Directory("github.com", "matryer", "is")
	is.True(strings.HasPrefix(dir, expected))
}

func TestResolveDirectoryNested(t *testing.T) {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	modCache := modcache.Default()
	module, err := mod.Find(wd, mod.WithModCache(modCache))
	is.NoErr(err)
	dir, err := module.ResolveDirectory("golang.org/x/mod/modfile")
	is.NoErr(err)
	// $GOMODCACHE/golang.org/x/mod@v0.5.1/modfile
	prefix := modCache.Directory("golang.org", "x", "mod")
	suffix := "/modfile"
	is.True(strings.HasPrefix(dir, prefix))
	is.True(strings.HasSuffix(dir, suffix))
}

func TestResolveDirectoryNestedSame(t *testing.T) {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	modCache := modcache.Default()
	module, err := mod.Find(wd, mod.WithModCache(modCache))
	is.NoErr(err)
	dir, err := module.ResolveDirectory("gitlab.com/mnm/bud/internal/modcache")
	is.NoErr(err)
	expected := module.Directory("internal/modcache")
	is.Equal(dir, expected)
}

func TestResolveDirectoryNotOk(t *testing.T) {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	module, err := mod.Find(wd)
	is.NoErr(err)
	dir, err := module.ResolveDirectory("github.com/matryer/is/zargle")
	is.True(errors.Is(err, os.ErrNotExist))
	is.Equal(dir, "")
}

func TestResolveStdDirectory(t *testing.T) {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	module, err := mod.Find(wd)
	is.NoErr(err)
	dir, err := module.ResolveDirectory("net/http")
	is.NoErr(err)
	expected := filepath.Join(build.Default.GOROOT, "src", "net", "http")
	is.Equal(dir, expected)
}

func TestResolveImport(t *testing.T) {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	module, err := mod.Find(wd)
	is.NoErr(err)
	im, err := module.ResolveImport(wd)
	is.NoErr(err)
	base := filepath.Base(wd)
	is.Equal(module.Import("2", base), im)
}

func TestFindStdlib(t *testing.T) {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	module, err := mod.Find(wd)
	is.NoErr(err)
	module, err = module.Find("net/http")
	is.NoErr(err)
	imp, err := module.ResolveImport(module.Directory())
	is.NoErr(err)
	is.Equal(imp, "std")
}

func containsName(des []fs.DirEntry, name string) bool {
	for _, de := range des {
		if de.Name() == name {
			return true
		}
	}
	return false
}

func contains(des []fs.DirEntry, names ...string) bool {
	for _, name := range names {
		if !containsName(des, name) {
			return false
		}
	}
	return true
}

func TestOpen(t *testing.T) {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	module, err := mod.Find(wd)
	is.NoErr(err)
	des, err := fs.ReadDir(module, ".")
	is.NoErr(err)
	is.Equal(true, contains(des, "go.mod", "main.go"))
}

func TestFileCacheDir(t *testing.T) {
	is := is.New(t)
	appDir := t.TempDir()
	err := vfs.Write(appDir, vfs.Map{
		"go.mod":  []byte(`module app.com`),
		"main.go": []byte(`package main`),
	})
	is.NoErr(err)
	fmap := virtual.FileMap()
	module, err := mod.Find(appDir, mod.WithFileCache(fmap))
	is.NoErr(err)
	// Check initial
	des, err := fs.ReadDir(module, ".")
	is.NoErr(err)
	is.Equal(2, len(des))
	is.Equal("go.mod", des[0].Name())
	is.Equal("main.go", des[1].Name())
	// Delete a file and check again to ensure it's cached
	err = os.RemoveAll(filepath.Join(appDir, "main.go"))
	is.NoErr(err)
	des, err = fs.ReadDir(module, ".")
	is.NoErr(err)
	is.Equal(2, len(des))
	is.Equal("go.mod", des[0].Name())
	is.Equal("main.go", des[1].Name())
	// Delete from the cache to see that it's been updated
	fmap.Delete("main.go")
	des, err = fs.ReadDir(module, ".")
	is.NoErr(err)
	is.Equal(1, len(des))
	is.Equal("go.mod", des[0].Name())
	// Create a file and see stale dir
	err = os.WriteFile(filepath.Join(appDir, "main.go"), []byte(`package main`), 0644)
	is.NoErr(err)
	des, err = fs.ReadDir(module, ".")
	is.NoErr(err)
	is.Equal(1, len(des))
	is.Equal("go.mod", des[0].Name())
	// Mark the cache as creating main.go and see that it's been updated
	fmap.Create("main.go")
	des, err = fs.ReadDir(module, ".")
	is.NoErr(err)
	is.Equal(2, len(des))
	is.Equal("go.mod", des[0].Name())
	is.Equal("main.go", des[1].Name())
}
