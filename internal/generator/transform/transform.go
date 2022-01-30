package transform

import (
	_ "embed"
	"io/fs"

	"gitlab.com/mnm/bud/internal/entrypoint"

	"gitlab.com/mnm/bud/budfs"
	"gitlab.com/mnm/bud/gen"
	"gitlab.com/mnm/bud/internal/gotemplate"
	"gitlab.com/mnm/bud/internal/imports"
	"gitlab.com/mnm/bud/pkg/gomod"
)

//go:embed transform.gotext
var template string

var generator = gotemplate.MustParse("transform.gotext", template)

type Generator struct {
	BFS    budfs.FS
	Module *gomod.Module
}

type State struct {
	Imports    []*imports.Import
	Transforms []*Transform
}

type Transform struct {
	From     string
	To       string
	Platform string
	Variable string
	Type     string
	Function string
}

func (g *Generator) GenerateFile(_ gen.F, file *gen.File) error {
	views, err := entrypoint.List(g.BFS)
	if err != nil {
		return err
	} else if len(views) == 0 {
		return fs.ErrNotExist
	}
	imports := imports.New()
	imports.AddNamed("transform", "gitlab.com/mnm/bud/transform")
	imports.AddNamed("svelte", "gitlab.com/mnm/bud/svelte")
	code, err := generator.Generate(State{
		Imports: imports.List(),
	})
	if err != nil {
		return err
	}
	file.Write(code)
	return nil
}
