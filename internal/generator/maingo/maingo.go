package maingo

import (
	_ "embed"

	"gitlab.com/mnm/bud/gen"
	"gitlab.com/mnm/bud/go/mod"
	"gitlab.com/mnm/bud/internal/gotemplate"
	"gitlab.com/mnm/bud/internal/imports"
)

//go:embed maingo.gotext
var template string

var generator = gotemplate.MustParse("maingo", template)

type Generator struct {
	Modfile *mod.File
}

type State struct {
	Imports []*imports.Import
}

func (g *Generator) GenerateFile(f gen.F, file *gen.File) error {
	if err := gen.Exists(f, "bud/command/command.go"); err != nil {
		return err
	}
	imports := imports.New()
	imports.AddStd("os")
	// imports.AddStd("fmt")
	imports.Add(g.Modfile.ModulePath("bud/command"))
	code, err := generator.Generate(State{
		Imports: imports.List(),
	})
	if err != nil {
		return err
	}
	file.Write(code)
	return nil
}
