package typecheck_test

import (
	"strings"
	"testing"

	"github.com/livebud/bud/internal/typecheck"
	"github.com/livebud/bud/package/gomod"

	"github.com/livebud/bud/internal/testdir"
	"github.com/matryer/is"
)

func TestSimpleOk(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	td := testdir.New()
	td.Files["main.go"] = `
		package main
		import "app.com/say"
		func main() {
			println(say.Hello())
		}
	`
	td.Files["say/hello.go"] = `
		package say
		func Hello() string { return "hello" }
	`
	err := td.Write(dir)
	is.NoErr(err)
	module, err := gomod.Find(dir)
	is.NoErr(err)
	tc := typecheck.New(module, module)
	err = tc.Check(".")
	is.NoErr(err)
}

func TestSimpleError(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	td := testdir.New()
	td.Files["main.go"] = `
		package main
		import "app.com/say"
		func main() {
			println(say.Hello())
		}
	`
	td.Files["say/hello.go"] = `
		package say
		func Hello() int { return "hello" }
	`
	err := td.Write(dir)
	is.NoErr(err)
	module, err := gomod.Find(dir)
	is.NoErr(err)
	tc := typecheck.New(module, module)
	err = tc.Check(".")
	is.True(err != nil)
	is.True(strings.Contains(err.Error(), "as int value in return statement"))
}

func TestErrorsPackage(t *testing.T) {
	t.SkipNow()
	is := is.New(t)
	dir := t.TempDir()
	td := testdir.New()
	td.Files["main.go"] = `
		package main
		import "errors"
		func main() {
			println(errors.New("some error").Error())
		}
	`
	err := td.Write(dir)
	is.NoErr(err)
	module, err := gomod.Find(dir)
	is.NoErr(err)
	tc := typecheck.New(module, module)
	err = tc.Check(".")
	is.NoErr(err)
}
