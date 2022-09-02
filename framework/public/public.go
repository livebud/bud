package public

import (
	"context"
	_ "embed"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/framework/public/publicrt"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/overlay"

	"github.com/livebud/bud/internal/gotemplate"
)

//go:embed public.gotext
var template string

var generator = gotemplate.MustParse("framework/public/public.gotext", template)

// Generate the public file
func Generate(state *State) ([]byte, error) {
	return generator.Generate(state)
}

// New public generator
func New(flag *framework.Flag, module *gomod.Module) *Generator {
	return &Generator{
		flag:   flag,
		module: module,
	}
}

type Generator struct {
	flag   *framework.Flag
	module *gomod.Module
}

func (g *Generator) GenerateDir(ctx context.Context, _ overlay.F, dir *overlay.Dir) error {
	fsys, err := publicrt.LoadFS(g.module)
	if err != nil {
		return err
	}
	dir.GenerateFile("public.go", func(ctx context.Context, _ overlay.F, file *overlay.File) error {
		state, err := Load(fsys, g.flag)
		if err != nil {
			return err
		}
		code, err := Generate(state)
		if err != nil {
			return err
		}
		file.Data = code
		return nil
	})
	return nil
}
