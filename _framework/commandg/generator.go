package command

import (
	"context"
	_ "embed"

	"gitlab.com/mnm/bud/framework/commandp"
	"gitlab.com/mnm/bud/internal/gotemplate"
)

//go:embed command.gotext
var template string

var generator = gotemplate.MustParse("command.gotext", template)

func Generate(ctx context.Context, state *commandp.State) ([]byte, error) {
	return generator.Generate(state)
}
