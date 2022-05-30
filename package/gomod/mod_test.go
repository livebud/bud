package gomod_test

import (
	"context"
	"errors"
	"go/build"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/livebud/bud/internal/fscache"
	"github.com/livebud/bud/internal/testdir"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/modcache"
	"github.com/livebud/bud/package/vfs"

	"github.com/livebud/bud/internal/is"
)

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

func TestFind(t *testing.T) {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	module, err := gomod.Find(wd)
	is.NoErr(err)
	dir := module.Directory()
	root := filepath.Join(wd, "..", "..")
	is.Equal(root, dir)
}

func TestFindDefault(t *testing.T) {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	module, err := gomod.Find(wd)
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
	module, err := gomod.Find(wd, gomod.WithModCache(modCache))
	is.NoErr(err)
	dir, err := module.ResolveDirectory("github.com/evanw/esbuild")
	is.NoErr(err)
	expected := modCache.Directory("github.com", "evanw", "esbuild")
	is.True(strings.HasPrefix(dir, expected))
}

func TestResolveDirectoryNested(t *testing.T) {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	modCache := modcache.Default()
	module, err := gomod.Find(wd, gomod.WithModCache(modCache))
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
	module, err := gomod.Find(wd, gomod.WithModCache(modCache))
	is.NoErr(err)
	dir, err := module.ResolveDirectory("github.com/livebud/bud/package/modcache")
	is.NoErr(err)
	expected := module.Directory("package/modcache")
	is.Equal(dir, expected)
}

func TestResolveDirectoryNotOk(t *testing.T) {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	module, err := gomod.Find(wd)
	is.NoErr(err)
	dir, err := module.ResolveDirectory("github.com/matryer/is/zargle")
	is.True(errors.Is(err, os.ErrNotExist))
	is.Equal(dir, "")
}

func TestResolveStdDirectory(t *testing.T) {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	module, err := gomod.Find(wd)
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
	module, err := gomod.Find(wd)
	is.NoErr(err)
	im, err := module.ResolveImport(wd)
	is.NoErr(err)
	base := filepath.Base(wd)
	is.Equal(module.Import("package", base), im)
}

func TestModuleFindStdlib(t *testing.T) {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	module, err := gomod.Find(wd)
	is.NoErr(err)
	module, err = module.Find("net/http")
	is.NoErr(err)
	imp, err := module.ResolveImport(module.Directory())
	is.NoErr(err)
	is.Equal(imp, "std")
}

func TestFindNested(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Modules["github.com/livebud/bud-test-plugin"] = "v0.0.8"
	err := td.Write(ctx)
	is.NoErr(err)
	modCache := modcache.Default()
	module1, err := gomod.Find(dir)
	is.NoErr(err)
	data, err := fs.ReadFile(module1, "go.mod")
	is.NoErr(err)
	m1, err := gomod.Parse("go.mod", data)
	is.NoErr(err)
	is.Equal(m1.Import(), "app.com")

	module2, err := module1.Find("github.com/livebud/bud-test-plugin")
	is.NoErr(err)
	is.Equal(module2.Import(), "github.com/livebud/bud-test-plugin")
	is.Equal(module2.Directory(), modCache.Directory("github.com/livebud", "bud-test-plugin@v0.0.8"))
	data, err = fs.ReadFile(module2, "go.mod")
	is.NoErr(err)
	m2, err := gomod.Parse("go.mod", data)
	is.NoErr(err)
	is.Equal(m2.Import(), "github.com/livebud/bud-test-plugin")

	// Find the nested module from bud-test-plugin
	req := module2.File().Require("github.com/livebud/bud-test-nested-plugin")
	is.True(req != nil)
	module3, err := module2.Find("github.com/livebud/bud-test-nested-plugin")
	is.NoErr(err)
	is.Equal(module3.Import(), "github.com/livebud/bud-test-nested-plugin")
	is.Equal(module3.Directory(), modCache.Directory("github.com/livebud", "bud-test-nested-plugin@"+req.Version))
	data, err = fs.ReadFile(module3, "go.mod")
	is.NoErr(err)
	m3, err := gomod.Parse("go.mod", data)
	is.NoErr(err)
	is.Equal(m3.Import(), "github.com/livebud/bud-test-nested-plugin")

	// Ensure module1 is not overriden
	is.Equal(module1.Import(), "app.com")
	is.Equal(module1.Directory(), dir)

	// Ensure module2 is not overriden
	is.Equal(module2.Import(), "github.com/livebud/bud-test-plugin")
	is.Equal(module2.Directory(), modCache.Directory("github.com/livebud", "bud-test-plugin@v0.0.8"))
}

func TestOpen(t *testing.T) {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	module, err := gomod.Find(wd)
	is.NoErr(err)
	des, err := fs.ReadDir(module, ".")
	is.NoErr(err)
	is.Equal(true, contains(des, "go.mod", "main.go"))
}

func TestFileCacheDir(t *testing.T) {
	t.SkipNow()
	is := is.New(t)
	appDir := t.TempDir()
	err := vfs.Write(appDir, vfs.Map{
		"go.mod":  []byte(`module app.com`),
		"main.go": []byte(`package main`),
	})
	is.NoErr(err)
	fmap := fscache.New()
	module, err := gomod.Find(appDir, gomod.WithFSCache(fmap))
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

func TestModuleFindLocal(t *testing.T) {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	module1, err := gomod.Find(wd)
	is.NoErr(err)
	// Find local web directory within module1
	module2, err := module1.Find(module1.Import("runtime", "web"))
	is.NoErr(err)
	is.Equal(module1.Directory(), module2.Directory())
}

func TestModuleFindFromFS(t *testing.T) {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	module1, err := gomod.Find(wd)
	is.NoErr(err)
	// First ensure the package doesn't exist
	module2, err := module1.Find(module1.Import("imagine"))
	is.Equal(nil, module2)
	is.True(errors.Is(err, fs.ErrNotExist))
	// Now find the package using a virtual FS
	vfs := vfs.Map{
		"imagine/imagine.go": []byte(`package imagine`),
	}
	module2, err = module1.FindIn(vfs, module1.Import("imagine"))
	is.NoErr(err)
	is.Equal(module1.Directory(), module2.Directory())
	absDir, err := module2.ResolveDirectoryIn(vfs, module1.Import("imagine"))
	is.NoErr(err)
	is.Equal(module1.Directory("imagine"), absDir)
}

func TestDirFS(t *testing.T) {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	module, err := gomod.Find(wd)
	is.NoErr(err)
	fsys := module.DirFS("internal")
	is.NoErr(err)
	des, err := fs.ReadDir(fsys, ".")
	is.NoErr(err)
	is.Equal(true, contains(des, "bail", "dsync"))
}

func TestHash(t *testing.T) {
	is := is.New(t)
	m1, err := gomod.Parse("go.mod", []byte(`module app.test`))
	is.NoErr(err)
	m2, err := gomod.Parse("go.mod", []byte(`module app.test`))
	is.NoErr(err)
	is.Equal(string(m1.Hash()), string(m2.Hash()))
	m3, err := gomod.Parse("go.mod", []byte(`module apptest`))
	is.NoErr(err)
	is.True(string(m2.Hash()) != string(m3.Hash()))
}

func TestMissingModule(t *testing.T) {
	is := is.New(t)
	module, err := gomod.Parse("go.mod", []byte(`require github.com/evanw/esbuild v0.14.11`))
	is.True(err != nil)
	is.Equal(err.Error(), `mod: missing module statement in "go.mod", received "require github.com/evanw/esbuild v0.14.11\n"`)
	is.Equal(module, nil)
}
