package programg_test

import (
	"testing"

	"gitlab.com/mnm/bud/framework/programg"
	"gitlab.com/mnm/bud/pkg/di"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/internal/testdir"
	"gitlab.com/mnm/bud/package/overlay"
	"gitlab.com/mnm/bud/pkg/gomod"
	"gitlab.com/mnm/bud/pkg/parser"
)

func TestEmpty(t *testing.T) {
	is := is.New(t)
	td := testdir.New()
	dir := t.TempDir()
	err := td.Write(dir)
	is.NoErr(err)
	module, err := gomod.Find(dir)
	is.NoErr(err)
	overlay, err := overlay.Load(module)
	is.NoErr(err)
	parser := parser.New(overlay, module)
	injector := di.New(overlay, module, parser)
	overlay.GenerateFile("bud/.cli/program/program.go", programg.New(injector, module, &di.Function{
		Name: "loadCLI",
	}))
	pkg, err := parser.Parse("bud/.cli/program")
	is.NoErr(err)
	is.Equal(pkg.Name(), "program")
}
