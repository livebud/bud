package public

import (
	"context"
	_ "embed"

	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/overlay"

	"github.com/livebud/bud/internal/gotemplate"
	"github.com/livebud/bud/runtime/bud"
)

//go:embed public.gotext
var template string

var generator = gotemplate.MustParse("public", template)

func New(flag *bud.Flag, module *gomod.Module) *Generator {
	return &Generator{
		flag:   flag,
		module: module,
	}
}

type Generator struct {
	flag   *bud.Flag
	module *gomod.Module
}

func (g *Generator) GenerateFile(ctx context.Context, fsys overlay.F, file *overlay.File) error {
	state, err := Load(g.flag, fsys, g.module)
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
