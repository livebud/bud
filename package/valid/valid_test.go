package valid_test

import (
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/valid"
)

func TestDir(t *testing.T) {
	is := is.New(t)
	is.True(valid.Dir("a"))
	is.True(valid.Dir("ab"))
	is.True(valid.Dir("a_"))
	is.True(valid.Dir("ab_"))
	is.True(valid.Dir("ab-"))
	is.True(!valid.Dir("aA"))
	is.True(!valid.Dir("aB"))
	is.True(!valid.Dir(""))
	is.True(!valid.Dir("_a"))
	is.True(!valid.Dir("_ab"))
	is.True(!valid.Dir(".a"))
	is.True(!valid.Dir(".ab"))
	is.True(!valid.Dir("A"))
	is.True(!valid.Dir("Ab"))
	is.True(!valid.Dir("bud"))
}

func TestView(t *testing.T) {
	is := is.New(t)
	is.True(valid.View("a.jsx"))
	is.True(valid.View("a.svelte"))
	is.True(valid.View("a.svg"))
	is.True(valid.View("ab.jsx"))
	is.True(valid.View("ab.svelte"))
	is.True(valid.View("ab.svg"))
	is.True(valid.View("a_.jsx"))
	is.True(valid.View("ab_.jsx"))
	is.True(valid.View("ab-.jsx"))
	is.True(valid.View("aA.jsx"))
	is.True(valid.View("aB.jsx"))
	is.True(!valid.View(""))
	is.True(!valid.View("_a"))
	is.True(!valid.View("_ab"))
	is.True(!valid.View(".a"))
	is.True(!valid.View(".ab"))
	is.True(!valid.View("A"))
	is.True(!valid.View("Ab"))
	is.True(!valid.View("bud"))
	is.True(!valid.View("go.mod"))
	is.True(!valid.View("go.sum"))
	is.True(!valid.View("package.json"))
}

func TestControllerFile(t *testing.T) {
	is := is.New(t)
	is.True(valid.ControllerFile("a.go"))
	is.True(valid.ControllerFile("ab.go"))
	is.True(valid.ControllerFile("a_.go"))
	is.True(valid.ControllerFile("ab_.go"))
	is.True(valid.ControllerFile("ab-.go"))
	is.True(valid.ControllerFile("aA.go"))
	is.True(valid.ControllerFile("aB.go"))
	is.True(valid.ControllerFile("A.go"))
	is.True(!valid.ControllerFile(""))
	is.True(!valid.ControllerFile("_a.go"))
	is.True(!valid.ControllerFile("_ab.go"))
	is.True(!valid.ControllerFile(".a.go"))
	is.True(!valid.ControllerFile(".ab.go"))
	is.True(!valid.ControllerFile("a"))
	is.True(!valid.ControllerFile("A"))
	is.True(!valid.ControllerFile("Ab"))
	is.True(!valid.ControllerFile("bud"))
	is.True(!valid.ControllerFile("bud.go"))
}
