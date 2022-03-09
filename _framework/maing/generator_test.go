package maing_test

import (
	"testing"

	"gitlab.com/mnm/bud/framework/maing"
	"gitlab.com/mnm/bud/package/overlay"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/internal/testdir"
	"gitlab.com/mnm/bud/pkg/gomod"
	"gitlab.com/mnm/bud/pkg/parser"
)

func TestParse(t *testing.T) {
	is := is.New(t)
	td := testdir.New()
	dir := t.TempDir()
	err := td.Write(dir)
	is.NoErr(err)
	module, err := gomod.Find(dir)
	is.NoErr(err)
	overlay, err := overlay.Load(module)
	is.NoErr(err)
	overlay.GenerateFile("bud/.cli/main.go", maing.New("bud/.cli/program"))
	parser := parser.New(overlay, module)
	pkg, err := parser.Parse("bud/.cli")
	is.NoErr(err)
	is.Equal(pkg.Name(), "main")
}
