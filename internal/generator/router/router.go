package router

import (
	_ "embed"

	"gitlab.com/mnm/bud/gen"
	"gitlab.com/mnm/bud/go/mod"
	"gitlab.com/mnm/bud/internal/gotemplate"
	"gitlab.com/mnm/bud/internal/imports"
)

//go:embed router.gotext
var template string

var generator = gotemplate.MustParse("router.gotext", template)

type Generator struct {
	Module *mod.Module
}

func (g *Generator) GenerateFile(f gen.F, file *gen.File) error {
	if err := gen.SkipUnless(f, "bud/action/action.go"); err != nil {
		return err
	}
	imports := imports.New()
	imports.AddStd("net/http")
	imports.AddNamed("router", "gitlab.com/mnm/bud/router")
	imports.AddNamed("action", g.Module.Import("bud/action"))
	imports.AddNamed("public", g.Module.Import("bud/public"))
	imports.AddNamed("view", g.Module.Import("bud/view"))
	code, err := generator.Generate(State{
		Imports: imports.List(),
	})
	if err != nil {
		return err
	}
	file.Write(code)
	return nil
}
