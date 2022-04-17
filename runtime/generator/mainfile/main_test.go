package mainfile_test

import (
	"context"
	"testing"

	"github.com/matryer/is"
	"github.com/livebud/bud/internal/budtest"
)

// TODO: We should always generate a main, even if empty dir
func TestEmpty(t *testing.T) {
	t.SkipNow()
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	bud := budtest.New(dir)
	project, err := bud.Compile(ctx)
	is.NoErr(err)
	app, err := project.Build(ctx)
	is.NoErr(err)
	is.NoErr(app.Exists("bud/.app/main.go"))
}
