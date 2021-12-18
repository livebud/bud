package mod_test

import (
	"errors"
	"go/build"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gitlab.com/mnm/bud/gen"
	"gitlab.com/mnm/bud/go/mod"
	"gitlab.com/mnm/bud/internal/modcache"
	"gitlab.com/mnm/bud/vfs"

	"github.com/matryer/is"
)

func TestFind(t *testing.T) {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	modFinder := mod.New()
	modFile, err := modFinder.Find(wd)
	is.NoErr(err)
	dir := modFile.Directory()
	root := filepath.Join(wd, "..", "..")
	is.Equal(root, dir)
}
func TestFindDefault(t *testing.T) {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	modFile, err := mod.New().Find(wd)
	is.NoErr(err)
	dir := modFile.Directory()
	root := filepath.Join(wd, "..", "..")
	is.Equal(root, dir)
}

func TestResolveDirectory(t *testing.T) {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	modCache := modcache.Default()
	modFinder := mod.New(mod.WithCache(modCache))
	module, err := modFinder.Find(wd)
	is.NoErr(err)
	dir, err := module.ResolveDirectory("github.com/matryer/is")
	is.NoErr(err)
	expected := modCache.Directory("github.com", "matryer", "is")
	is.True(strings.HasPrefix(dir, expected))
}

func TestResolveDirectoryNotOk(t *testing.T) {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	modFinder := mod.New()
	module, err := modFinder.Find(wd)
	is.NoErr(err)
	dir, err := module.ResolveDirectory("github.com/matryer/is/zargle")
	is.True(errors.Is(err, os.ErrNotExist))
	is.Equal(dir, "")
}

func TestResolveStdDirectory(t *testing.T) {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	modFinder := mod.New()
	module, err := modFinder.Find(wd)
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
	modFinder := mod.New()
	module, err := modFinder.Find(wd)
	is.NoErr(err)
	im, err := module.ResolveImport(wd)
	is.NoErr(err)
	base := filepath.Base(wd)
	is.Equal(module.Import("go", base), im)
}

func TestVirtualResolveImport(t *testing.T) {
	is := is.New(t)
	fsys := vfs.Map{
		"app/go.mod":         `module app.com/app`,
		"app/main.go":        "package app\nfunc main() {}",
		"app/dep/dep.go":     "package dep\nfunc dep() {}",
		"outside/outside.go": "package outside\nfunc outside() {}",
	}
	modFinder := mod.New(mod.WithFS(fsys))
	module, err := modFinder.Find("app")
	is.NoErr(err)
	im, err := module.ResolveImport("app")
	is.NoErr(err)
	is.Equal(im, "app.com/app")
	im, err = module.ResolveImport("app/dep")
	is.NoErr(err)
	is.Equal(im, "app.com/app/dep")
	im, err = module.ResolveImport("outside")
	is.True(err != nil)
	is.Equal(err.Error(), `"outside" can't be outside the module directory "app"`)
	is.Equal(im, "")
}

func TestAddRequire(t *testing.T) {
	is := is.New(t)
	modFinder := mod.New()
	modPath := filepath.Join(t.TempDir(), "go.mod")
	module, err := modFinder.Parse(modPath, []byte(`module app.test`))
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
	modFinder := mod.New()
	modPath := filepath.Join(t.TempDir(), "go.mod")
	module, err := modFinder.Parse(modPath, []byte(`module app.test`))
	is.NoErr(err)
	modFile := module.File()
	modFile.AddReplace("mod.test/two", "", "mod.test/twotwo", "")
	modFile.AddReplace("mod.test/one", "", "mod.test/oneone", "")
	is.Equal(string(modFile.Format()), `module app.test

replace mod.test/two => mod.test/twotwo

replace mod.test/one => mod.test/oneone
`)
}

// func TestPluginRequire(t *testing.T) {
// 	is := is.New(t)
// 	dir := "_tmp"
// 	is.NoErr(os.RemoveAll(dir))
// 	is.NoErr(os.MkdirAll(dir, 0755))
// 	defer func() {
// 		if !t.Failed() {
// 			is.NoErr(os.RemoveAll(dir))
// 		}
// 	}()
// 	err := vfs.Write(dir, vfs.Map{
// 		"go.mod": `module github.com/livebud/test`,
// 	})
// 	is.NoErr(err)
// 	ctx := context.Background()
// 	err = gobin.Get(ctx, dir, "gitlab.com/mnm/testdata/bud-tailwind", "gitlab.com/mnm/testdata/bud-markdown")
// 	is.NoErr(err)
// 	modFile, err := mod.Load(dir)
// 	is.NoErr(err)
// 	plugins, err := modFile.Plugins()
// 	is.NoErr(err)
// 	is.Equal(len(plugins), 2) // expected 2 plugins
// 	// First plugin
// 	is.Equal(plugins[0].Import, "gitlab.com/mnm/testdata/bud-markdown")
// 	expected := filepath.Join(modFile.CacheDir, "gitlab.com", "mnm", "testdata", "bud-markdown")
// 	is.True(strings.HasPrefix(plugins[0].Dir, expected))
// 	is.Equal(plugins[0].Name, "bud-markdown")
// 	// Second plugin
// 	is.Equal(plugins[1].Import, "gitlab.com/mnm/testdata/bud-tailwind")
// 	expected = filepath.Join(modFile.CacheDir, "gitlab.com", "mnm", "testdata", "bud-tailwind")
// 	is.True(strings.HasPrefix(plugins[1].Dir, expected))
// 	is.Equal(plugins[1].Name, "bud-tailwind")
// }

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
	modFinder := mod.New(mod.WithCache(modCache))
	module, err := modFinder.Parse(filepath.Join(appDir, "go.mod"), []byte(`module app.test`))
	is.NoErr(err)
	modFile := module.File()
	modFile.AddRequire("mod.test/module", "v1.2.4")
	dir, err := module.ResolveDirectory("mod.test/module")
	is.NoErr(err)
	is.Equal(dir, filepath.Join(cacheDir, "mod.test", "module@v1.2.4"))
}

// func TestModCacheRead(t *testing.T) {
// 	is := is.New(t)
// 	modPath := filepath.Join(t.TempDir(), "go.mod")
// 	module:= mod.New()
// 	modFile, err := mod.Parse(modCache, modPath, []byte(`
// 		module mod.test

// 		require github.com/matryer/is v1.4.0
// 	`))
// 	is.NoErr(err)
// 	des, err := fs.ReadDir(modFile, "github.com/matryer/is")
// 	is.NoErr(err)
// 	is.True(len(des) > 0)
// }

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
		"go.mod": "module app.com\nrequire mod.test/module v1.2.4",
		"app.go": "package app\nimport \"mod.test/module\"\nvar a = module.Answer",
	})
	is.NoErr(err)
	modFinder := mod.New(mod.WithCache(modCache))
	module1, err := modFinder.Find(appDir)
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
		"go.mod": "module app.com\nrequire mod.test/module v1.2.4",
		"app.go": "package app\nimport \"mod.test/module\"\nvar a = module.Answer",
	})
	is.NoErr(err)
	modFinder := mod.New(mod.WithCache(modCache))
	module1, err := modFinder.Find(appDir)
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

func TestFindFSNesting(t *testing.T) {
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
	fsys := vfs.Map{
		"app.go": "package app\nimport \"mod.test/module\"\nvar a = module.Answer",
	}
	genfs := gen.New(fsys)
	genfs.Add(map[string]gen.Generator{
		"go.mod": gen.GenerateFile(func(f gen.F, file *gen.File) error {
			file.Write([]byte("module app.com\nrequire mod.test/module v1.2.4"))
			return nil
		}),
		"hello/hello.go": gen.GenerateFile(func(f gen.F, file *gen.File) error {
			file.Write([]byte("package hello\nconst Hi = \"hello\""))
			return nil
		}),
	})
	modFinder := mod.New(mod.WithCache(modCache), mod.WithFS(genfs))
	module, err := modFinder.Find(".")
	is.NoErr(err)
	is.Equal(module.Directory(), ".")
	// Resolving a local package
	dir, err := module.ResolveDirectory("app.com/hello")
	is.NoErr(err)
	is.Equal(dir, "hello")
	// Resolve dependencies
	module2, err := module.Find("mod.test/module")
	is.NoErr(err)
	is.Equal(module2.Import(), "mod.test/module")
	is.Equal(module2.Directory(), modCache.Directory("mod.test", "module@v1.2.4"))
	module3, err := module2.Find("mod.test/two")
	is.NoErr(err)
	is.Equal(module3.Import(), "mod.test/two")
	is.Equal(module3.Directory(), modCache.Directory("mod.test", "two@v0.0.2"))
}

func TestOpen(t *testing.T) {
	is := is.New(t)
	fsys := vfs.Map{
		"app.go": "package app\nvar a = 10",
	}
	genfs := gen.New(fsys)
	genfs.Add(map[string]gen.Generator{
		"go.mod": gen.GenerateFile(func(f gen.F, file *gen.File) error {
			file.Write([]byte("module app.com"))
			return nil
		}),
		"hello/hello.go": gen.GenerateFile(func(f gen.F, file *gen.File) error {
			file.Write([]byte("package hello\nconst Hi = \"hello\""))
			return nil
		}),
	})
	modFinder := mod.New(mod.WithFS(genfs))
	module, err := modFinder.Find(".")
	is.NoErr(err)
	is.Equal(module.Import(), "app.com")
	des, err := fs.ReadDir(module, ".")
	is.NoErr(err)
	is.Equal(len(des), 3)
	is.Equal(des[0].Name(), "app.go")
	is.Equal(des[1].Name(), "go.mod")
	is.Equal(des[2].Name(), "hello")
	is.Equal(des[2].IsDir(), true)
	code, err := fs.ReadFile(module, "hello/hello.go")
	is.NoErr(err)
	is.Equal(string(code), "package hello\nconst Hi = \"hello\"")
}
