package generator_test

import (
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/2/generator"
)

func TestGenerator(t *testing.T) {
	is := is.New(t)
	generator, err := generator.Load(".")
	is.NoErr(err)
}
