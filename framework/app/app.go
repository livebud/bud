package app

import (
	_ "embed"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/package/budfs"
	"github.com/livebud/bud/package/di"

	"github.com/livebud/bud/internal/gotemplate"
	"github.com/livebud/bud/package/gomod"
)

//go:embed app.gotext
var template string

var generator = gotemplate.MustParse("framework/app/app.gotext", template)

func Generate(state *State) ([]byte, error) {
	return generator.Generate(state)
}

func New(injector *di.Injector, module *gomod.Module, flag *framework.Flag) *Generator {
	return &Generator{flag, injector, module}
}

type Generator struct {
	flag     *framework.Flag
	injector *di.Injector
	module   *gomod.Module
}

func (g *Generator) GenerateFile(fsys budfs.FS, file *budfs.File) error {
	state, err := Load(fsys, g.injector, g.module, g.flag)
	if err != nil {
		return err
	}
	code, err := Generate(state)
	if err != nil {
		return err
	}
	file.Data = code
	return nil
}
