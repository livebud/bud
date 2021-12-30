package valid_test

import (
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/internal/valid"
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

func TestActionFile(t *testing.T) {
	is := is.New(t)
	is.True(valid.ActionFile("a.go"))
	is.True(valid.ActionFile("ab.go"))
	is.True(valid.ActionFile("a_.go"))
	is.True(valid.ActionFile("ab_.go"))
	is.True(valid.ActionFile("ab-.go"))
	is.True(valid.ActionFile("aA.go"))
	is.True(valid.ActionFile("aB.go"))
	is.True(valid.ActionFile("A.go"))
	is.True(!valid.ActionFile(""))
	is.True(!valid.ActionFile("_a.go"))
	is.True(!valid.ActionFile("_ab.go"))
	is.True(!valid.ActionFile(".a.go"))
	is.True(!valid.ActionFile(".ab.go"))
	is.True(!valid.ActionFile("a"))
	is.True(!valid.ActionFile("A"))
	is.True(!valid.ActionFile("Ab"))
	is.True(!valid.ActionFile("bud"))
	is.True(!valid.ActionFile("bud.go"))
}
