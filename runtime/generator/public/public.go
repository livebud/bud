package public

import (
	"context"
	_ "embed"
	"io/fs"

	"gitlab.com/mnm/bud/package/overlay"
	"gitlab.com/mnm/bud/pkg/gomod"

	"gitlab.com/mnm/bud/internal/gotemplate"
)

//go:embed public.gotext
var template string

var generator = gotemplate.MustParse("public", template)

func New(fsys fs.FS, module *gomod.Module) *Generator {
	return &Generator{
		FS:     fsys,
		Module: module,
		// Embed  bool
		// Minify bool
	}
}

type Generator struct {
	FS     fs.FS
	Module *gomod.Module
	Embed  bool
	Minify bool
}

func (g *Generator) GenerateFile(ctx context.Context, _ overlay.F, file *overlay.File) error {
	state, err := Load(g.FS, g.Module, g.Embed, g.Minify)
	if err != nil {
		return err
	}
	code, err := generator.Generate(state)
	if err != nil {
		return err
	}
	file.Data = code
	return nil
}
