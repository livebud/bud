package plugin_test

import (
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/internal/test"
)

func TestEmpty(t *testing.T) {
	is := is.New(t)
	generator := test.Generator(t)
	app, err := generator.Generate()
	is.NoErr(err)
}
