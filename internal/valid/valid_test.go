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
	is.True(valid.Dir("aA"))
	is.True(valid.Dir("aB"))
	is.True(!valid.Dir(""))
	is.True(!valid.Dir("_a"))
	is.True(!valid.Dir("_ab"))
	is.True(!valid.Dir(".a"))
	is.True(!valid.Dir(".ab"))
	is.True(!valid.Dir("A"))
	is.True(!valid.Dir("Ab"))
}
