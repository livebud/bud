package bud

import (
	"context"
	"io"

	"github.com/livebud/bud/internal/config"
	"github.com/livebud/bud/internal/pubsub"

	"github.com/livebud/bud/package/commander"
	"github.com/livebud/bud/package/socket"
)

// Input contains the configuration that gets passed into the commands
type Input struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
	Env    []string

	// Currently passed in only for testing
	Dir   string          // Can be empty
	BudLn socket.Listener // Can be nil
	WebLn socket.Listener // Can be nil
	Bus   pubsub.Client   // Can be nil
}

func New(provide config.Provide) *Command {
	return &Command{provide: provide}
}

type Command struct {
	provide config.Provide
	Args    []string
	Help    bool
}

// Run a custom command
// TODO: finish supporting custom commands
//  1. Compile
//     a. Generate generator (later!)
//     i. Generate bud/internal/generator
//     ii. Build bud/generator
//     iii. Run bud/generator
//     b. Generate custom command
//     i. Generate bud/internal/command/${name}/
//     ii. Build bud/command/${name}
//  2. Run bud/command/${name}
func (c *Command) Run(ctx context.Context) error {
	return commander.Usage()
}
