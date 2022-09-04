package generate

import (
	_ "embed"

	"github.com/livebud/bud/internal/gotemplate"
	"github.com/livebud/bud/package/budfs"
	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/gomod"
)

//go:embed main.gotext
var template string

var generator = gotemplate.MustParse("framework/generate/main.gotext", template)

func Generate(state *State) ([]byte, error) {
	return generator.Generate(state)
}

func New(injector *di.Injector, module *gomod.Module) *Generator {
	return &Generator{injector, module}
}

type Generator struct {
	injector *di.Injector
	module   *gomod.Module
}

func (g *Generator) GenerateFile(fsys *budfs.FS, file *budfs.File) error {
	state, err := Load(fsys, g.injector, g.module)
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
