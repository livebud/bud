package generate

import (
	_ "embed"

	"gitlab.com/mnm/bud/budfs"
	"gitlab.com/mnm/bud/mod"

	"gitlab.com/mnm/bud/gen"
	"gitlab.com/mnm/bud/internal/gotemplate"
	"gitlab.com/mnm/bud/internal/imports"
)

//go:embed generate.gotext
var template string

var generator = gotemplate.MustParse("generate", template)

type Generator struct {
	BFS    budfs.FS
	Module *mod.Module
	Embed  bool
	Hot    bool
	Minify bool
}

type State struct {
	*Generator
	Imports []*imports.Import
}

func (g *Generator) GenerateFile(_ gen.F, file *gen.File) error {
	// Don't create a generate file if custom user-generators don't exist
	if err := gen.SkipUnless(g.BFS, "bud/generator/generator.go"); err != nil {
		return err
	}
	imports := imports.New()
	imports.AddStd("os", "fmt")
	imports.AddNamed("gen", "gitlab.com/mnm/bud/gen")
	imports.AddNamed("generator", g.Module.Import("bud/generator"))
	code, err := generator.Generate(&State{
		Imports:   imports.List(),
		Generator: g,
	})
	if err != nil {
		return err
	}
	file.Write(code)
	return nil
}
