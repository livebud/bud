package view

import (
	"context"
	_ "embed"
	"io/fs"

	"gitlab.com/mnm/bud/internal/gotemplate"
	"gitlab.com/mnm/bud/internal/imports"
	"gitlab.com/mnm/bud/package/overlay"
	"gitlab.com/mnm/bud/pkg/gomod"
	"gitlab.com/mnm/bud/pkg/vfs"
)

//go:embed view.gotext
var template string

var generator = gotemplate.MustParse("view.gotext", template)

type Generator struct {
	FS     fs.FS
	Module *gomod.Module
}

type State struct {
	Imports []*imports.Import
}

func (g *Generator) GenerateFile(ctx context.Context, _ overlay.F, file *overlay.File) error {
	exist := vfs.SomeExist(g.FS, "view")
	if len(exist) == 0 {
		return fs.ErrNotExist
	}
	imports := imports.New()
	imports.AddNamed("transform", g.Module.Import("bud/.app/transform"))
	imports.AddNamed("overlay", "gitlab.com/mnm/bud/package/overlay")
	imports.AddNamed("mod", "gitlab.com/mnm/bud/pkg/gomod")
	imports.AddNamed("js", "gitlab.com/mnm/bud/pkg/js")
	imports.AddNamed("view", "gitlab.com/mnm/bud/runtime/view")
	code, err := generator.Generate(State{
		Imports: imports.List(),
	})
	if err != nil {
		return err
	}
	file.Data = code
	return nil
}
