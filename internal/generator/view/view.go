package view

import (
	_ "embed"
	"io/fs"

	"gitlab.com/mnm/bud/budfs"
	"gitlab.com/mnm/bud/gen"
	"gitlab.com/mnm/bud/internal/gotemplate"
	"gitlab.com/mnm/bud/internal/imports"
	"gitlab.com/mnm/bud/pkg/gomod"
	"gitlab.com/mnm/bud/vfs"
)

//go:embed view.gotext
var template string

var generator = gotemplate.MustParse("view.gotext", template)

type Generator struct {
	BFS    budfs.FS
	Module *gomod.Module
}

type State struct {
	Imports []*imports.Import
}

func (g *Generator) GenerateFile(_ gen.F, file *gen.File) error {
	exist := vfs.SomeExist(g.BFS, "view")
	if len(exist) == 0 {
		return fs.ErrNotExist
	}
	imports := imports.New()
	imports.AddNamed("transform", g.Module.Import("bud/transform"))
	imports.AddNamed("gen", "gitlab.com/mnm/bud/gen")
	imports.AddNamed("mod", "gitlab.com/mnm/bud/pkg/gomod")
	imports.AddNamed("js", "gitlab.com/mnm/bud/js")
	imports.AddNamed("view", "gitlab.com/mnm/bud/view")
	code, err := generator.Generate(State{
		Imports: imports.List(),
	})
	if err != nil {
		return err
	}
	file.Write(code)
	return nil
}
