package controller

import (
	_ "embed"

	"gitlab.com/mnm/bud/gen"
	"gitlab.com/mnm/bud/go/mod"
	"gitlab.com/mnm/bud/internal/gotemplate"
	"gitlab.com/mnm/bud/internal/imports"
)

//go:embed controller.gotext
var template string

var generator = gotemplate.MustParse("controller.gotext", template)

type Generator struct {
	Module *mod.Module
}

// var controllers = glob.MustCompile("controller/{*,**}.go")

func (g *Generator) GenerateFile(f gen.F, file *gen.File) error {

	imports := imports.New()
	imports.AddStd("net/http")
	imports.AddNamed("view", g.Module.Import("bud/view"))
	// TODO: replace with dynamic list
	imports.AddNamed("controller", g.Module.Import("controller"))
	code, err := generator.Generate(State{
		Imports: imports.List(),
	})
	if err != nil {
		return err
	}
	file.Write(code)
	return nil
}
