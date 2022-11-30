package glob_test

import (
	"testing"

	"github.com/livebud/bud/internal/glob"
	"github.com/livebud/bud/internal/is"
)

func TestMatch(t *testing.T) {
	is := is.New(t)
	matcher, err := glob.Compile("{controller/**.go,view/**}")
	is.NoErr(err)
	is.True(matcher.Match("controller/controller.go"))
	is.True(matcher.Match("view/index.svelte"))
}

func TestDirMatch(t *testing.T) {
	is := is.New(t)
	matcher, err := glob.Compile("controller/*/**.go")
	is.NoErr(err)
	is.True(!matcher.Match("controller"))
	is.True(!matcher.Match("controller/controller.go"))
	is.True(matcher.Match("controller/view/view.go"))
	is.True(matcher.Match("controller/public/public.go"))
}

func TestMatchSubdir(t *testing.T) {
	is := is.New(t)
	matcher, err := glob.Compile(`{generator/**.go,bud/internal/generator/*/**.go}`)
	is.NoErr(err)
	is.True(!matcher.Match("bud/internal/generator/generator.go"))
}
