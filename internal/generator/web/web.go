package web

import (
	_ "embed"

	"gitlab.com/mnm/bud/gen"
	"gitlab.com/mnm/bud/go/mod"
	"gitlab.com/mnm/bud/internal/gotemplate"
)

//go:embed web.gotext
var template string

var generator = gotemplate.MustParse("web", template)

type Generator struct {
	Module *mod.Module
}

func (g *Generator) GenerateFile(f gen.F, file *gen.File) error {
	if err := gen.SkipUnless(f, "bud/router/router.go"); err != nil {
		return err
	}
	state, err := Load(g.Module, f)
	if err != nil {
		return err
	}
	code, err := generator.Generate(state)
	if err != nil {
		return err
	}
	file.Write(code)
	return nil
}
