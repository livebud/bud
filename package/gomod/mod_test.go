package gomod_test

import (
	"errors"
	"go/build"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/livebud/bud/internal/fscache"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/modcache"
	"github.com/livebud/bud/package/vfs"

	"github.com/matryer/is"
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

func TestAddRequire(t *testing.T) {
	is := is.New(t)
	modPath := filepath.Join(t.TempDir(), "go.mod")
	module, err := gomod.Parse(modPath, []byte(`module app.test`))
	is.NoErr(err)
	modFile := module.File()
	modFile.AddRequire("mod.test/two", "v2")
	modFile.AddRequire("mod.test/one", "v1.2.4")
	is.Equal(string(modFile.Format()), `module app.test

require (
	mod.test/two v2
	mod.test/one v1.2.4
)
`)
}

func TestAddReplace(t *testing.T) {
	is := is.New(t)
	modPath := filepath.Join(t.TempDir(), "go.mod")
	module, err := gomod.Parse(modPath, []byte(`module app.test`))
	is.NoErr(err)
	modFile := module.File()
	modFile.AddReplace("mod.test/two", "", "mod.test/twotwo", "")
	modFile.AddReplace("mod.test/one", "", "mod.test/oneone", "")
	is.Equal(string(modFile.Format()), `module app.test

replace mod.test/two => mod.test/twotwo

replace mod.test/one => mod.test/oneone
`)
}

func TestLocalResolveDirectory(t *testing.T) {
	is := is.New(t)
	cacheDir := t.TempDir()
	modCache := modcache.New(cacheDir)
	err := modCache.Write(modcache.Modules{
		"mod.test/module@v1.2.3": modcache.Files{
			"go.mod":   "module mod.test/module\n\ngo 1.12",
			"const.go": "package module\n\nconst Answer = 42",
		},
		"mod.test/module@v1.2.4": modcache.Files{
			"go.mod":   "module mod.test/module\n\ngo 1.12",
			"const.go": "package module\n\nconst Answer = 43",
		},
	})
	is.NoErr(err)
	appDir := t.TempDir()
	modPath := filepath.Join(appDir, "go.mod")
	modData := []byte(`module app.test`)
	module, err := gomod.Parse(modPath, modData, gomod.WithModCache(modCache))
	is.NoErr(err)
	modFile := module.File()
	modFile.AddRequire("mod.test/module", "v1.2.4")
	dir, err := module.ResolveDirectory("mod.test/module")
	is.NoErr(err)
	is.Equal(dir, filepath.Join(cacheDir, "mod.test", "module@v1.2.4"))
}

func TestFindNested(t *testing.T) {
	is := is.New(t)
	cacheDir := t.TempDir()
	modCache := modcache.New(cacheDir)
	err := modCache.Write(modcache.Modules{
		"mod.test/module@v1.2.3": modcache.Files{
			"go.mod":   "module mod.test/module",
			"const.go": "package module\nconst Answer = 42",
		},
		"mod.test/module@v1.2.4": modcache.Files{
			"go.mod":   "module mod.test/module",
			"const.go": "package module\nconst Answer = 43",
		},
	})
	is.NoErr(err)
	appDir := t.TempDir()
	err = vfs.Write(appDir, vfs.Map{
		"go.mod": []byte("module app.com\nrequire mod.test/module v1.2.4"),
		"app.go": []byte("package app\nimport \"mod.test/module\"\nvar a = module.Answer"),
	})
	is.NoErr(err)
	module1, err := gomod.Find(appDir, gomod.WithModCache(modCache))
	is.NoErr(err)

	module2, err := module1.Find("mod.test/module")
	is.NoErr(err)
	is.Equal(module2.Import(), "mod.test/module")
	is.Equal(module2.Directory(), modCache.Directory("mod.test", "module@v1.2.4"))

	// Ensure module1 is not overriden
	is.Equal(module1.Import(), "app.com")
	is.Equal(module1.Directory(), appDir)
}

func TestFindNestedFS(t *testing.T) {
	is := is.New(t)
	cacheDir := t.TempDir()
	modCache := modcache.New(cacheDir)
	err := modCache.Write(modcache.Modules{
		"mod.test/two@v0.0.1": modcache.Files{
			"go.mod":   "module mod.test/two",
			"const.go": "package two\nconst Answer = 10",
		},
		"mod.test/two@v0.0.2": modcache.Files{
			"go.mod":   "module mod.test/two",
			"const.go": "package two\nconst Answer = 20",
		},
		"mod.test/module@v1.2.3": modcache.Files{
			"go.mod":   "module mod.test/module",
			"const.go": "package module\nconst Answer = 42",
		},
		"mod.test/module@v1.2.4": modcache.Files{
			"go.mod":   "module mod.test/module\nrequire mod.test/two v0.0.2",
			"const.go": "package module\nimport \"mod.test/two\"\nconst Answer = two.Answer",
		},
	})
	is.NoErr(err)
	appDir := t.TempDir()
	err = vfs.Write(appDir, vfs.Map{
		"go.mod": []byte("module app.com\nrequire mod.test/module v1.2.4"),
		"app.go": []byte("package app\nimport \"mod.test/module\"\nvar a = module.Answer"),
	})
	is.NoErr(err)
	module1, err := gomod.Find(appDir, gomod.WithModCache(modCache))
	is.NoErr(err)

	module2, err := module1.Find("mod.test/module")
	is.NoErr(err)
	is.Equal(module2.Import(), "mod.test/module")
	is.Equal(module2.Directory(), modCache.Directory("mod.test", "module@v1.2.4"))

	module3, err := module2.Find("mod.test/two")
	is.NoErr(err)
	is.Equal(module3.Import(), "mod.test/two")
	is.Equal(module3.Directory(), modCache.Directory("mod.test", "two@v0.0.2"))

	// Ensure module1 is not overriden
	is.Equal(module1.Import(), "app.com")
	is.Equal(module1.Directory(), appDir)

	// Ensure module2 is not overriden
	is.Equal(module2.Import(), "mod.test/module")
	is.Equal(module2.Directory(), modCache.Directory("mod.test", "module@v1.2.4"))
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
	modPath := filepath.Join(t.TempDir(), "go.mod")
	m1, err := gomod.Parse(modPath, []byte(`module app.test`))
	is.NoErr(err)
	m2, err := gomod.Parse(modPath, []byte(`module app.test`))
	is.NoErr(err)
	is.Equal(string(m1.Hash()), string(m2.Hash()))
	m3, err := gomod.Parse(modPath, []byte(`module apptest`))
	is.NoErr(err)
	is.True(string(m2.Hash()) != string(m3.Hash()))
}
