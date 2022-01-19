package mod_test

import (
	"errors"
	"go/build"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gitlab.com/mnm/bud/2/mod"
	"gitlab.com/mnm/bud/internal/modcache"

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
	module, err := mod.Find(wd, mod.WithCache(modCache))
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
	module, err := mod.Find(wd, mod.WithCache(modCache))
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
	module, err := mod.Find(wd, mod.WithCache(modCache))
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
