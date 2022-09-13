package glob_test

import (
	"testing"

	"github.com/gobwas/glob"
	"github.com/livebud/bud/internal/is"
)

func TestMatch(t *testing.T) {
	is := is.New(t)
	matcher, err := glob.Compile("{controller/**.go,view/**}")
	is.NoErr(err)
	is.True(matcher.Match("controller/controller.go"))
	is.True(matcher.Match("view/index.svelte"))
}
