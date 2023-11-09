package mod_test

import (
	"io/fs"
	"os"
	"strings"
	"testing"

	"github.com/livebud/bud/pkg/mod"
	"github.com/matryer/is"
)

func TestReadFile(t *testing.T) {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	module, err := mod.Find(wd)
	is.NoErr(err)
	gomod, err := fs.ReadFile(module, "go.mod")
	is.NoErr(err)
	is.True(strings.Contains(string(gomod), "module github.com/livebud/bud"))
}

func TestReadDir(t *testing.T) {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	module, err := mod.Find(wd)
	is.NoErr(err)
	des, err := fs.ReadDir(module, ".")
	is.NoErr(err)
	is.True(len(des) > 0)
}
