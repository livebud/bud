package mod_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/livebud/bud/pkg/mod"
	"github.com/matryer/is"
)

func TestFind(t *testing.T) {
	is := is.New(t)
	module, err := mod.Find()
	is.NoErr(err)
	wd, err := os.Getwd()
	is.NoErr(err)
	is.Equal(wd, filepath.Join(module.Directory(), "pkg", "mod"))
	is.Equal(module.Import(), "github.com/livebud/bud")
}

func TestFindIn(t *testing.T) {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	module, err := mod.FindIn(".")
	is.NoErr(err)
	is.Equal(wd, filepath.Join(module.Directory(), "pkg", "mod"))
	is.Equal(module.Import(), "github.com/livebud/bud")
}

func TestNew(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	module := mod.New(dir)
	is.Equal(dir, module.Directory())
	is.Equal(module.Import(), "change.me")
}
