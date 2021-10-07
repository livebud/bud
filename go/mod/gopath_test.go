package mod_test

import (
	"go/build"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-duo/bud/go/mod"
	"github.com/matryer/is"
)

// GO111MODULE=off GOPATH=/tmp/gopath go get github.com/sirupsen/logrus

func TestGoPathDirectory(t *testing.T) {
	// TODO: check that we need $GOPATH support
	t.Skipf("To run this test, we'll need to `go get` with `GO111MODULE=off`")
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	modfile := mod.Virtual("gowithduo.com/duo", wd)
	dir, err := modfile.ResolveDirectory("github.com/matryer/is")
	is.NoErr(err)
	expect := filepath.Join(build.Default.GOPATH, "src", "github.com", "matryer", "is")
	is.Equal(expect, dir)
}
func TestGoPathStdDirectory(t *testing.T) {
	t.Skipf("To run this test, we'll need to `go get` with `GO111MODULE=off`")
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	modfile := mod.Virtual("gowithduo.com/duo", wd)
	dir, err := modfile.ResolveDirectory("net/http")
	is.NoErr(err)
	expected := filepath.Join(build.Default.GOROOT, "src", "net", "http")
	is.Equal(dir, expected)
}
