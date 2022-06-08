package build

import (
	"context"
	"fmt"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/command"
)

func New(bud *command.Bud) *Command {
	return &Command{
		bud:  bud,
		Flag: new(framework.Flag),
	}
}

type Command struct {
	bud  *command.Bud
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
	// console, err := c.bud.Console()
	// if err != nil {
	// 	return err
	// }
	module, err := c.bud.Module()
	if err != nil {
		return err
	}
	fmt.Println(module.Directory())
	fsys, err := c.bud.FileSystem(module, c.Flag)
	if err != nil {
		return err
	}
	if err := fsys.Sync("bud/internal/app"); err != nil {
		return err
	}
	builder := c.bud.Builder(module)
	return builder.Build(ctx, "bud/internal/app/main.go", "bud/app")
}
