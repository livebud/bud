package run

import (
	"context"
	"fmt"

	"github.com/livebud/bud/internal/cli3/generate"
	"github.com/livebud/bud/package/config"
)

func New(
	config *config.Config,
	generate *generate.Command,
) *Command {
	return &Command{
		config,
		generate,
		true,
	}
}

type Command struct {
	config   *config.Config
	generate *generate.Command
	Watch    bool
}

func (c *Command) Run(ctx context.Context) error {
	fmt.Println("running run", c.config.Hot)
	return nil
	// return c.generate.Run(ctx)
}
