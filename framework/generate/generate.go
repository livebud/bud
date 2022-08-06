package generate

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/livebud/bud/internal/gotemplate"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/overlay"
)

//go:embed main.gotext
var template string

var generator = gotemplate.MustParse("framework/generate/main.gotext", template)

func Generate(state *State) ([]byte, error) {
	return generator.Generate(state)
}

func New(module *gomod.Module) *Generator {
	return &Generator{module}
}

type Generator struct {
	module *gomod.Module
}

func (g *Generator) GenerateFile(ctx context.Context, fsys overlay.F, file *overlay.File) error {
	state, err := Load(fsys, g.module)
	if err != nil {
		return err
	}
	code, err := Generate(state)
	if err != nil {
		return err
	}
	fmt.Println(string(code))
	file.Data = code
	return nil
}
