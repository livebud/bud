package command

import (
	"context"
	_ "embed"

	"gitlab.com/mnm/bud/internal/gotemplate"
	"gitlab.com/mnm/bud/internal/parser/command"
	"gitlab.com/mnm/bud/package/overlay"
)

//go:embed command.gotext
var template string

var generator = gotemplate.MustParse("command.gotext", template)

func New(parser *command.Parser) *Generator {
	return &Generator{parser}
}

type Generator struct {
	parser *command.Parser
}

func (g *Generator) GenerateFile(ctx context.Context, f overlay.F, file *overlay.File) error {
	// Load command state
	state, err := g.parser.ParseCLI()
	if err != nil {
		return err
	}
	// Generate our template
	file.Data, err = generator.Generate(state)
	if err != nil {
		return err
	}
	return nil
}
