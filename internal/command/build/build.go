package build

import (
	"context"

	"github.com/livebud/bud/internal/command"
)

func New(bud *command.Bud) *Command {
	return &Command{
		bud: bud,
	}
}

type Command struct {
	bud *command.Bud

	// Flags
	Embed  bool
	Minify bool
	Hot    string
}

// Run the build command
// 1. Setup
// 2. Compile
func (c *Command) Run(ctx context.Context) error {
	return nil
}
