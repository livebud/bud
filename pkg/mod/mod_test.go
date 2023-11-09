package mod_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/livebud/bud/pkg/mod"
	"github.com/matryer/is"
)

func getBudDir() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(wd, "..", ".."), nil
}

func TestFind(t *testing.T) {
	is := is.New(t)
	budDir, err := getBudDir()
	is.NoErr(err)
	module, err := mod.Find()
	is.NoErr(err)
	is.Equal(budDir, module.Directory())
	is.Equal(module.Import(), "github.com/livebud/bud")
}

func TestMustFind(t *testing.T) {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	budDir, err := getBudDir()
	is.NoErr(err)
	module := mod.MustFind(wd)
	is.Equal(budDir, module.Directory())
	is.Equal(module.Import(), "github.com/livebud/bud")
}

func TestNew(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	module := mod.New(dir)
	is.Equal(dir, module.Directory())
	is.Equal(module.Import(), "change.me")
}
