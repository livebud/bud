package mainfile_test

import (
	"context"
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/framework2/mainfile"
	"gitlab.com/mnm/bud/internal/testdir"
	"gitlab.com/mnm/bud/pkg/gomod"
)

func parseForApp(dir string, files map[string]string) (*mainfile.State, error) {
	ctx := context.Background()
	td := testdir.New()
	err := td.Write(dir)
	if err != nil {
		return nil, err
	}
	module, err := gomod.Find(dir)
	if err != nil {
		return nil, err
	}
	app := mainfile.ForApp(module)
	return app.Parse(ctx)
}

func TestEmptyApp(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	state, err := parseForApp(dir, map[string]string{})
	is.NoErr(err)
	_ = state
}
