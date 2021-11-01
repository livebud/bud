package build

import (
	"context"
	"fmt"

	"gitlab.com/mnm/bud/bfs"
)

func New() *Command {
	return &Command{
		Hot:   false,
		Embed: true,
	}
}

type Command struct {
	Hot   bool
	Embed bool
}

func (c *Command) Run(ctx context.Context, generators map[string]bfs.Generator) error {
	fmt.Println("building code!")
	// 1. Run the generators
	// 2. go build bud/main.go
	return nil
}
