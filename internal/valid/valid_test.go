package valid_test

import (
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/valid"
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

func TestViewEntry(t *testing.T) {
	is := is.New(t)
	is.True(valid.ViewEntry("a"))
	is.True(valid.ViewEntry("ab"))
	is.True(valid.ViewEntry("a_"))
	is.True(valid.ViewEntry("ab_"))
	is.True(valid.ViewEntry("ab-"))
	is.True(valid.ViewEntry("aA"))
	is.True(valid.ViewEntry("aB"))
	is.True(!valid.ViewEntry(""))
	is.True(!valid.ViewEntry("_a"))
	is.True(!valid.ViewEntry("_ab"))
	is.True(!valid.ViewEntry(".a"))
	is.True(!valid.ViewEntry(".ab"))
	is.True(!valid.ViewEntry("A"))
	is.True(!valid.ViewEntry("Ab"))
	is.True(!valid.ViewEntry("bud"))
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
