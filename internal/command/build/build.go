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
//   a. Generate generator (later!)
//   	 i. Generate bud/internal/generator
//     ii. Build bud/generator
//     iii. Run bud/generator
//   b. Generate app
//     i. Generate bud/internal/app
//     ii. Build into bud/app
func (c *Command) Run(ctx context.Context) error {
	// console, err := c.bud.Console()
	// if err != nil {
	// 	return err
	// }
	module, err := c.bud.Module()
	if err != nil {
		return err
	}
	fsys, err := c.bud.FileSystem(module)
	if err != nil {
		return err
	}
	if err := fsys.Sync("bud/internal/app"); err != nil {
		return err
	}
	// TODO: Build into bud/app
	return nil
}
