package build

import (
	"context"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/cli/bud"
)

func New(bud *bud.Command) *Command {
	return &Command{
		bud:  bud,
		Flag: new(framework.Flag),
	}
}

type Command struct {
	bud  *bud.Command
	Flag *framework.Flag
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
	module, err := c.bud.Module()
	if err != nil {
		return err
	}
	genfs, err := c.bud.FileSystem(module, c.Flag)
	if err != nil {
		return err
	}
	if err := genfs.Sync("bud/internal/app"); err != nil {
		return err
	}
	return c.bud.Build(ctx, module, "bud/internal/app/main.go", "bud/app")
}
