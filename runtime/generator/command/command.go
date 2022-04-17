package command

import (
	"context"
	_ "embed"

	"github.com/livebud/bud/internal/gotemplate"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/overlay"
	"github.com/livebud/bud/package/parser"
)

//go:embed command.gotext
var template string

var generator = gotemplate.MustParse("command.gotext", template)

// TODO: rename to Command
type Generator struct {
	Module *gomod.Module
	Parser *parser.Parser
}

func (g *Generator) GenerateFile(ctx context.Context, fsys overlay.F, file *overlay.File) error {
	// Load command state
	state, err := Load(fsys, g.Module, g.Parser)
	if err != nil {
		return err
	}
	// Generate our template
	code, err := generator.Generate(state)
	if err != nil {
		return err
	}
	file.Data = code
	return nil
}
