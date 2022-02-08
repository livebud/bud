package expand_test

import (
	"context"
	"testing"

	"gitlab.com/mnm/bud/pkg/di"

	"gitlab.com/mnm/bud/pkg/gen"
	"gitlab.com/mnm/bud/pkg/parser"
	"gitlab.com/mnm/bud/pkg/pluginfs"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/internal/buddy/expand"
	"gitlab.com/mnm/bud/internal/testdir"
	"gitlab.com/mnm/bud/pkg/gomod"
)

func TestEmpty(t *testing.T) {
	t.SkipNow()
	is := is.New(t)
	td := testdir.New()
	dir := t.TempDir()
	err := td.Write(dir)
	is.NoErr(err)
	module, err := gomod.Find(dir)
	is.NoErr(err)
	pluginFS, err := pluginfs.Load(module)
	is.NoErr(err)
	genFS := gen.New(pluginFS)
	parser := parser.New(genFS, module)
	injector := di.New(genFS, module, parser)
	expander := expand.New(genFS, injector, module, parser)
	ctx := context.Background()
	err = expander.Expand(ctx, &expand.Input{})
	is.NoErr(err)
}
