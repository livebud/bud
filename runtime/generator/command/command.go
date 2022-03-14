package command

import (
	"context"
	_ "embed"
	"io/fs"

	"gitlab.com/mnm/bud/internal/gotemplate"
	"gitlab.com/mnm/bud/package/overlay"
	"gitlab.com/mnm/bud/pkg/gomod"
	"gitlab.com/mnm/bud/pkg/parser"
)

//go:embed command.gotext
var template string

var generator = gotemplate.MustParse("command.gotext", template)

// TODO: rename to Command
type Generator struct {
	FS     fs.FS
	Module *gomod.Module
	Parser *parser.Parser
}

func (g *Generator) GenerateFile(ctx context.Context, f overlay.F, file *overlay.File) error {
	// Load command state
	state, err := Load(g.FS, g.Module, g.Parser)
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
