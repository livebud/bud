package mod_test

import (
	"errors"
	"go/build"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-duo/bud/go/mod"
	"github.com/matryer/is"
)

func TestDirectory(t *testing.T) {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	dir, err := mod.Directory(wd)
	is.NoErr(err)
	root := filepath.Join(wd, "..", "..")
	is.Equal(root, dir)
}

func TestLoadDirectory(t *testing.T) {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	modfile, err := mod.Load(wd)
	is.NoErr(err)
	dir := modfile.Directory()
	root := filepath.Join(wd, "..", "..")
	is.Equal(root, dir)
}

func TestResolveDirectory(t *testing.T) {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	modfile, err := mod.Load(wd)
	is.NoErr(err)
	dir, err := modfile.ResolveDirectory("github.com/matryer/is")
	is.NoErr(err)
	expected := filepath.Join(mod.GOMODCACHE, "github.com", "matryer", "is")
	is.True(strings.HasPrefix(dir, expected))
}

func TestResolveDirectoryNotOk(t *testing.T) {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	modfile, err := mod.Load(wd)
	is.NoErr(err)
	dir, err := modfile.ResolveDirectory("github.com/matryer/is/zargle")
	is.Equal(dir, "")
	is.True(errors.Is(err, os.ErrNotExist))
}

func TestResolveStdDirectory(t *testing.T) {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	modfile, err := mod.Load(wd)
	is.NoErr(err)
	dir, err := modfile.ResolveDirectory("net/http")
	is.NoErr(err)
	expected := filepath.Join(build.Default.GOROOT, "src", "net", "http")
	is.Equal(dir, expected)
}

func TestResolveImport(t *testing.T) {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	modfile, err := mod.Load(wd)
	is.NoErr(err)
	im, err := modfile.ResolveImport(wd)
	is.NoErr(err)
	is.Equal(path.Join(modfile.ModulePath(), "go", "mod"), im)
}

// TODO: test mod.Load(dir)
