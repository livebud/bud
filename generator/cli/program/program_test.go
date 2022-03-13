package program_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/generator/cli/program"
	"gitlab.com/mnm/bud/pkg/di"
)

func parse(code []byte) (*ast.File, error) {
	fset := token.NewFileSet()
	return parser.ParseFile(fset, "", code, parser.DeclarationErrors)
}

func TestBasic(t *testing.T) {
	is := is.New(t)
	code, err := program.Generate(&program.State{
		Provider: &di.Provider{
			Name: "loadCLI",
		},
	})
	is.NoErr(err)
	file, err := parse(code)
	is.NoErr(err)
	is.Equal(file.Name.Name, "program")
}
