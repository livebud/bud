package mainfile_test

import (
	"context"
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/framework2/mainfile"
	"gitlab.com/mnm/bud/internal/testdir"
	"gitlab.com/mnm/bud/pkg/gomod"
)

func parseForCLI(dir string, files map[string]string) (*mainfile.State, error) {
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
	app := mainfile.ForCLI(module)
	return app.Parse(ctx)
}

func TestEmptyCLI(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	state, err := parseForCLI(dir, map[string]string{})
	is.NoErr(err)
	_ = state
}
