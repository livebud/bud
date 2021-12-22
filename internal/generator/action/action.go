package action

import (
	_ "embed"

	"gitlab.com/mnm/bud/gen"
	"gitlab.com/mnm/bud/go/mod"
	"gitlab.com/mnm/bud/internal/gotemplate"
	"gitlab.com/mnm/bud/internal/imports"
)

//go:embed action.gotext
var template string

var generator = gotemplate.MustParse("action.gotext", template)

type Generator struct {
	Module *mod.Module
}

func (g *Generator) GenerateFile(f gen.F, file *gen.File) error {
	imports := imports.New()
	imports.AddStd("net/http")
	imports.AddNamed("view", "gitlab.com/mnm/bud/view")
	// TODO: replace with dynamic list
	imports.AddNamed("action", g.Module.Import("action"))
	code, err := generator.Generate(State{
		Imports: imports.List(),
	})
	if err != nil {
		return err
	}
	file.Write(code)
	return nil
}
