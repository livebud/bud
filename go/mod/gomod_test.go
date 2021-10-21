package mod_test

import (
	"context"
	"errors"
	"go/build"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/go/mod"
	"gitlab.com/mnm/bud/internal/gobin"
	"gitlab.com/mnm/bud/internal/vfs"
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
	expected := filepath.Join(mod.GOMODCACHE(), "github.com", "matryer", "is")
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

func TestRequirePlugin(t *testing.T) {
	is := is.New(t)
	dir := "_tmp"
	is.NoErr(os.RemoveAll(dir))
	is.NoErr(os.MkdirAll(dir, 0755))
	defer func() {
		if !t.Failed() {
			is.NoErr(os.RemoveAll(dir))
		}
	}()
	err := vfs.WriteTo(dir, vfs.Map{
		"go.mod": `module github.com/livebud/test`,
	})
	is.NoErr(err)
	ctx := context.Background()
	err = gobin.GoGet(ctx, dir, "gitlab.com/mnm/testdata/bud-tailwind", "gitlab.com/mnm/testdata/bud-markdown")
	is.NoErr(err)
	modfile, err := mod.Find(dir)
	is.NoErr(err)
	plugins, err := modfile.Plugins()
	is.NoErr(err)
	is.Equal(len(plugins), 2) // expected 2 plugins
	// First plugin
	is.Equal(plugins[0].Import, "gitlab.com/mnm/testdata/bud-markdown")
	expected := filepath.Join(mod.GOMODCACHE(), "gitlab.com", "mnm", "testdata", "bud-markdown")
	is.True(strings.HasPrefix(plugins[0].Dir, expected))
	is.Equal(plugins[0].Name, "bud-markdown")
	// Second plugin
	is.Equal(plugins[1].Import, "gitlab.com/mnm/testdata/bud-tailwind")
	expected = filepath.Join(mod.GOMODCACHE(), "gitlab.com", "mnm", "testdata", "bud-tailwind")
	is.True(strings.HasPrefix(plugins[1].Dir, expected))
	is.Equal(plugins[1].Name, "bud-tailwind")
}

// TODO: test mod.Load(dir)
// TODO: test replaced plugin
