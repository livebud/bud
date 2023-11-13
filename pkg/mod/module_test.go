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

func TestSub(t *testing.T) {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	module, err := mod.Find(wd)
	is.NoErr(err)
	subdir := module.Sub("pkg", "mod")
	des, err := fs.ReadDir(subdir, ".")
	is.NoErr(err)
	is.True(len(des) >= 4)
	is.Equal(des[0].Name(), "mod.go")
	is.Equal(des[1].Name(), "mod_test.go")
	is.Equal(des[2].Name(), "module.go")
	is.Equal(des[3].Name(), "module_test.go")
}
