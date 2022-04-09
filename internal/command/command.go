package command

import (
	"context"

	"gitlab.com/mnm/bud/internal/bud"
	"gitlab.com/mnm/bud/package/trace"
	runtime_bud "gitlab.com/mnm/bud/runtime/bud"
)

// Bud command
type Bud struct {
	Flag runtime_bud.Flag
	Dir  string
	Args []string
}

// Run a custom command
func (c *Bud) Run(ctx context.Context) (err error) {
	// Start the trace
	ctx, span := trace.Start(ctx, "running bud")
	defer span.End(&err)
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
