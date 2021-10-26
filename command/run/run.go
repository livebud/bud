package run

import (
	"context"
	"fmt"

	"gitlab.com/mnm/bud/bfs"
)

type Command struct {
	BFS bfs.BFS

	Addr string `flag:"address to bind to" short:"a" default:":3000"`
}

func (c *Command) Run(ctx context.Context) error {
	fmt.Println("running run command")
	return nil
}
