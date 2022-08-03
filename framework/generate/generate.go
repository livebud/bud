package generate

import (
	"context"
	_ "embed"

	"github.com/livebud/bud/internal/gotemplate"
	"github.com/livebud/bud/package/overlay"
)

//go:embed generate.gotext
var template string

var generator = gotemplate.MustParse("framework/generate/generate.gotext", template)

func Generate(state *State) ([]byte, error) {
	return generator.Generate(state)
}

func New() *Generator {
	return &Generator{}
}

type Generator struct {
}

func (g *Generator) GenerateFile(ctx context.Context, fsys overlay.F, file *overlay.File) error {
	state, err := Load(fsys)
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
