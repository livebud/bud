package command

import (
	"context"

	"github.com/livebud/bud/internal/bud"
	"github.com/livebud/bud/package/commander"
	runtime_bud "github.com/livebud/bud/runtime/bud"
)

func New() *Bud {
	return &Bud{}
}

// Bud command
type Bud struct {
	Flag runtime_bud.Flag
	Dir  string
	Args []string
}

// Run a custom command
func (c *Bud) Run(ctx context.Context) (err error) {
	if len(c.Args) == 0 {
		return commander.Usage()
	}
	// Load the compiler
	compiler, err := bud.Find(c.Dir)
	if err != nil {
		return err
	}
	// Compiler the project CLI
	project, err := compiler.Compile(ctx, &c.Flag)
	if err != nil {
		return err
	}
	// Run the custom command
	return project.Execute(ctx, c.Args...)
}
