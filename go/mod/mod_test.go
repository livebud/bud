package mod_test

import (
	"errors"
	"go/build"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"gitlab.com/mnm/bud/go/mod"
	"gitlab.com/mnm/bud/internal/modcache"
	"gitlab.com/mnm/bud/vfs"

	"github.com/matryer/is"
)

func TestFind(t *testing.T) {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	module := mod.New(modcache.Default())
	modFile, err := module.Find(wd)
	is.NoErr(err)
	dir := modFile.Directory()
	root := filepath.Join(wd, "..", "..")
	is.Equal(root, dir)
}
func TestFindDefault(t *testing.T) {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	modFile, err := mod.Default().Find(wd)
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
	module := mod.New(modCache)
	modFile, err := module.Find(wd)
	is.NoErr(err)
	dir, err := modFile.ResolveDirectory("github.com/matryer/is")
	is.NoErr(err)
	expected := modCache.Directory("github.com", "matryer", "is")
	is.True(strings.HasPrefix(dir, expected))
}

func TestResolveDirectoryNotOk(t *testing.T) {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	module := mod.New(modcache.Default())
	modFile, err := module.Find(wd)
	is.NoErr(err)
	dir, err := modFile.ResolveDirectory("github.com/matryer/is/zargle")
	is.Equal(dir, "")
	is.True(errors.Is(err, os.ErrNotExist))
}

func TestResolveStdDirectory(t *testing.T) {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	module := mod.New(modcache.Default())
	modFile, err := module.Find(wd)
	is.NoErr(err)
	dir, err := modFile.ResolveDirectory("net/http")
	is.NoErr(err)
	expected := filepath.Join(build.Default.GOROOT, "src", "net", "http")
	is.Equal(dir, expected)
}

func TestResolveImport(t *testing.T) {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	module := mod.New(modcache.Default())
	modFile, err := module.Find(wd)
	is.NoErr(err)
	im, err := modFile.ResolveImport(wd)
	is.NoErr(err)
	is.Equal(path.Join(modFile.ModulePath(), "go", "mod"), im)
}

func TestAddRequire(t *testing.T) {
	is := is.New(t)
	module := mod.New(modcache.Default())
	modPath := filepath.Join(t.TempDir(), "go.mod")
	modFile, err := module.Parse(modPath, []byte(`module app.test`))
	is.NoErr(err)
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
	module := mod.New(modcache.Default())
	modPath := filepath.Join(t.TempDir(), "go.mod")
	modFile, err := module.Parse(modPath, []byte(`module app.test`))
	is.NoErr(err)
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
	err := modCache.Write(modcache.Versions{
		"v1.2.3": modcache.Files{
			"go.mod":   "module mod.test/module\n\ngo 1.12",
			"const.go": "package module\n\nconst Answer = 42",
		},
		"v1.2.4": modcache.Files{
			"go.mod":   "module mod.test/module\n\ngo 1.12",
			"const.go": "package module\n\nconst Answer = 43",
		},
	})
	is.NoErr(err)
	appDir := t.TempDir()
	module := mod.New(modCache)
	modFile, err := module.Parse(filepath.Join(appDir, "go.mod"), []byte(`module app.test`))
	is.NoErr(err)
	modFile.AddRequire("mod.test/module", "v1.2.4")
	dir, err := modFile.ResolveDirectory("mod.test/module")
	is.NoErr(err)
	is.Equal(dir, filepath.Join(cacheDir, "mod.test", "module@v1.2.4"))
}

// func TestModCacheRead(t *testing.T) {
// 	is := is.New(t)
// 	modPath := filepath.Join(t.TempDir(), "go.mod")
// 	module:= mod.New(modcache.Default())
// 	modFile, err := mod.Parse(modCache, modPath, []byte(`
// 		module mod.test

// 		require github.com/matryer/is v1.4.0
// 	`))
// 	is.NoErr(err)
// 	des, err := fs.ReadDir(modFile, "github.com/matryer/is")
// 	is.NoErr(err)
// 	is.True(len(des) > 0)
// }

func TestLoadCustom(t *testing.T) {
	is := is.New(t)
	cacheDir := t.TempDir()
	modCache := modcache.New(cacheDir)
	err := modCache.Write(modcache.Versions{
		"v1.2.3": modcache.Files{
			"go.mod":   "module mod.test/module",
			"const.go": "package module\nconst Answer = 42",
		},
		"v1.2.4": modcache.Files{
			"go.mod":   "module mod.test/module",
			"const.go": "package module\nconst Answer = 43",
		},
	})
	appDir := t.TempDir()
	vfs.Write(appDir, vfs.Map{
		"go.mod": "module app.com\nrequire mod.test/module v1.2.4",
		"app.go": "package app\nimport \"mod.test/module\"\nvar a = module.Answer",
	})
	is.NoErr(err)
	module := mod.New(modCache)
	modfile1, err := module.Find(appDir)
	is.NoErr(err)

	modfile2, err := modfile1.Load("mod.test/module")
	is.NoErr(err)
	is.Equal(modfile2.ModulePath(), "mod.test/module")
	is.Equal(modfile2.Directory(), modCache.Directory("mod.test", "module@v1.2.4"))

	// Ensure modfile1 is not overriden
	is.Equal(modfile1.ModulePath(), "app.com")
	is.Equal(modfile1.Directory(), appDir)
}
