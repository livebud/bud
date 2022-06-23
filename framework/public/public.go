package public

import (
	"context"
	_ "embed"

	"github.com/livebud/bud/framework"
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
func New(flag *framework.Flag) *Generator {
	return &Generator{
		flag: flag,
	}
}

type Generator struct {
	flag *framework.Flag
}

// Generate the public file
func (g *Generator) GenerateFile(ctx context.Context, fsys overlay.F, file *overlay.File) error {
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
}
