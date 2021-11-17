package mod_test

import (
	"errors"
	"go/build"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"gitlab.com/mnm/bud/go/mod"
	"gitlab.com/mnm/bud/internal/modcache"

	"github.com/matryer/is"
)

func TestFindBy(t *testing.T) {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	modCache := modcache.Default()
	modFile, err := mod.FindIn(modCache, wd)
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
	modFile, err := mod.FindIn(modCache, wd)
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
	modCache := modcache.Default()
	modFile, err := mod.FindIn(modCache, wd)
	is.NoErr(err)
	dir, err := modFile.ResolveDirectory("github.com/matryer/is/zargle")
	is.Equal(dir, "")
	is.True(errors.Is(err, os.ErrNotExist))
}

func TestResolveStdDirectory(t *testing.T) {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	modCache := modcache.Default()
	modFile, err := mod.FindIn(modCache, wd)
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
	modCache := modcache.Default()
	modFile, err := mod.FindIn(modCache, wd)
	is.NoErr(err)
	im, err := modFile.ResolveImport(wd)
	is.NoErr(err)
	is.Equal(path.Join(modFile.ModulePath(), "go", "mod"), im)
}

func TestAddRequire(t *testing.T) {
	is := is.New(t)
	modCache := modcache.Default()
	modPath := filepath.Join(t.TempDir(), "go.mod")
	modFile, err := mod.Parse(modCache, modPath, []byte(`module app.test`))
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
	modCache := modcache.Default()
	modPath := filepath.Join(t.TempDir(), "go.mod")
	modFile, err := mod.Parse(modCache, modPath, []byte(`module app.test`))
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
// 	err := vfs.WriteTo(dir, vfs.Map{
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
	var modules = []struct {
		version string
		files   map[string][]byte
	}{
		{
			version: "v1.2.3",
			files: map[string][]byte{
				"go.mod":   []byte("module mod.test/module\n\ngo 1.12"),
				"const.go": []byte("package module\n\nconst Answer = 42"),
			},
		},
		{
			version: "v1.2.4",
			files: map[string][]byte{
				"go.mod":   []byte("module mod.test/module\n\ngo 1.12"),
				"const.go": []byte("package module\n\nconst Answer = 43"),
			},
		},
	}
	for _, m := range modules {
		err := modCache.WriteModule(m.version, m.files)
		is.NoErr(err)
	}
	appDir := t.TempDir()
	modFile, err := mod.Parse(modCache, filepath.Join(appDir, "go.mod"), []byte(`module app.test`))
	is.NoErr(err)
	modFile.AddRequire("mod.test/module", "v1.2.4")
	dir, err := modFile.ResolveDirectory("mod.test/module")
	is.NoErr(err)
	is.Equal(dir, filepath.Join(cacheDir, "mod.test", "module@v1.2.4"))
}

func TestModCacheRead(t *testing.T) {
	is := is.New(t)
	modPath := filepath.Join(t.TempDir(), "go.mod")
	modCache := modcache.Default()
	modFile, err := mod.Parse(modCache, modPath, []byte(`
		module mod.test

		require github.com/matryer/is v1.4.0
	`))
	is.NoErr(err)
	des, err := fs.ReadDir(modFile, "github.com/matryer/is")
	is.NoErr(err)
	is.True(len(des) > 0)
}
