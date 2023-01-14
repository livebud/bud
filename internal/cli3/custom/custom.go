package custom

import (
	"context"
	"fmt"

	"github.com/livebud/bud/internal/cli3/generate"
)

func New(generate *generate.Command) *Command {
	return &Command{
		generate: generate,
	}
}

type Command struct {
	generate *generate.Command
	Help     bool
	Args     []string
}

func (c *Command) Run(ctx context.Context) error {
	fmt.Println("Custom!")
	return nil
}
