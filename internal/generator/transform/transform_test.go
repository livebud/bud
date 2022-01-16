package transform_test

import (
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/internal/test"
)

func TestEmpty(t *testing.T) {
	t.SkipNow()
	is := is.New(t)
	generator := test.Generator(t)
	app, err := generator.Generate()
	is.NoErr(err)
	is.Equal(false, app.Exists("bud/transform/transform.go"))
}
