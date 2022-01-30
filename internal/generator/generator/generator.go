package generator

import (
	_ "embed"

	"gitlab.com/mnm/bud/budfs"
	"gitlab.com/mnm/bud/pkg/gomod"

	"gitlab.com/mnm/bud/gen"
	"gitlab.com/mnm/bud/internal/gotemplate"
	"gitlab.com/mnm/bud/internal/imports"
)

//go:embed generator.gotext
var template string

var generator = gotemplate.MustParse("generator", template)

type Generator struct {
	BFS    budfs.FS
	Module *gomod.Module
	Embed  bool
	Hot    bool
	Minify bool
}

type State struct {
	Imports []*imports.Import
	Embed   bool
	Hot     bool
	Minify  bool
}

func (g *Generator) GenerateFile(_ gen.F, file *gen.File) error {
	// Don't create a generate file if custom user-generators don't exist
	// if err := gen.SkipUnless(g.BFS, "bud/generator/generator.go"); err != nil {
	// 	return err
	// }
	imports := imports.New()
	imports.AddNamed("public", "gitlab.com/mnm/bud/generator/public")
	// imports.AddStd("os", "fmt")
	imports.AddNamed("gen", "gitlab.com/mnm/bud/gen")
	// imports.AddNamed("generator", g.Module.Import("bud/generator"))
	code, err := generator.Generate(&State{
		Imports: imports.List(),
		Embed:   g.Embed,
		Hot:     g.Hot,
		Minify:  g.Minify,
	})
	if err != nil {
		return err
	}
	file.Write(code)
	return nil
}
