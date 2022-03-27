package public

import (
	"context"
	_ "embed"
	"io/fs"

	"gitlab.com/mnm/bud/package/gomod"
	"gitlab.com/mnm/bud/package/overlay"

	"gitlab.com/mnm/bud/internal/gotemplate"
	"gitlab.com/mnm/bud/runtime/bud"
)

//go:embed public.gotext
var template string

var generator = gotemplate.MustParse("public", template)

func New(flag *bud.Flag, fsys fs.FS, module *gomod.Module) *Generator {
	return &Generator{
		flag:   flag,
		fsys:   fsys,
		module: module,
	}
}

type Generator struct {
	flag   *bud.Flag
	fsys   fs.FS
	module *gomod.Module
}

func (g *Generator) GenerateFile(ctx context.Context, _ overlay.F, file *overlay.File) error {
	state, err := Load(g.flag, g.fsys, g.module)
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
