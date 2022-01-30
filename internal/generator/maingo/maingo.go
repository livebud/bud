package maingo

import (
	_ "embed"

	"gitlab.com/mnm/bud/budfs"
	"gitlab.com/mnm/bud/gen"
	"gitlab.com/mnm/bud/internal/gotemplate"
	"gitlab.com/mnm/bud/internal/imports"
	"gitlab.com/mnm/bud/pkg/gomod"
)

//go:embed maingo.gotext
var template string

var generator = gotemplate.MustParse("maingo", template)

type Generator struct {
	BFS    budfs.FS
	Module *gomod.Module
}

type State struct {
	Imports []*imports.Import
}

func (g *Generator) GenerateFile(_ gen.F, file *gen.File) error {
	if err := gen.SkipUnless(g.BFS, "bud/program/program.go"); err != nil {
		return err
	}
	imports := imports.New()
	imports.AddStd("os")
	// imports.AddStd("fmt")
	imports.Add(g.Module.Import("bud/program"))
	code, err := generator.Generate(State{
		Imports: imports.List(),
	})
	if err != nil {
		return err
	}
	file.Write(code)
	return nil
}
