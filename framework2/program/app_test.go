package program_test

import (
	"context"
	"errors"
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/framework2/program"
	"gitlab.com/mnm/bud/internal/testdir"
	"gitlab.com/mnm/bud/package/overlay"
	"gitlab.com/mnm/bud/pkg/di"
	"gitlab.com/mnm/bud/pkg/gomod"
	"gitlab.com/mnm/bud/pkg/parser"
)

func parseForApp(dir string, files map[string]string) (*program.State, error) {
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
	overlay, err := overlay.Load(module)
	if err != nil {
		return nil, err
	}
	parser := parser.New(overlay, module)
	injector := di.New(overlay, module, parser)
	app := program.ForApp(injector, module)
	return app.Parse(ctx)
}

func TestEmptyApp(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	state, err := parseForApp(dir, map[string]string{})
	is.True(errors.Is(err, program.ErrCantWire))
	is.Equal(state, nil)
}
