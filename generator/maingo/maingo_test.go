package maingo_test

import (
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/internal/fscache"
	"gitlab.com/mnm/bud/internal/test"
)

func TestEmpty(t *testing.T) {
	is := is.New(t)
	generator := test.Generator(t)
	fsCache := fscache.New()
	app, err := generator.Generate(fsCache)
	is.NoErr(err)
	is.True(!app.Exists("bud/main.go"))
}

func TestGoMod(t *testing.T) {
	is := is.New(t)
	generator := test.Generator(t)
	fsCache := fscache.New()
	app, err := generator.Generate(fsCache)
	is.NoErr(err)
	is.True(!app.Exists("bud/main.go"))
}
