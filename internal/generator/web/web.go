package web

import (
	_ "embed"

	"gitlab.com/mnm/bud/gen"
	"gitlab.com/mnm/bud/go/mod"
	"gitlab.com/mnm/bud/internal/gotemplate"
	"gitlab.com/mnm/bud/internal/imports"
)

//go:embed web.gotext
var template string

var generator = gotemplate.MustParse("web", template)

type Generator struct {
	Module *mod.Module
}

type State struct {
	Imports []*imports.Import
}

func (g *Generator) GenerateFile(f gen.F, file *gen.File) error {
	if err := gen.SkipUnless(f, "bud/router/router.go"); err != nil {
		return err
	}
	imports := imports.New()
	imports.AddStd("net/http", "context", "errors", "os")
	imports.AddNamed("hot", "gitlab.com/mnm/bud/hot")
	imports.AddNamed("middleware", "gitlab.com/mnm/bud/middleware")
	imports.AddNamed("commander", "gitlab.com/mnm/bud/commander")
	imports.AddNamed("router", g.Module.Import("bud/router"))
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
