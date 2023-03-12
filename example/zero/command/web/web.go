package web

import (
	"context"
	"fmt"
)

type Command struct {
}

type Serve struct {
	Listen string
}

func (c *Command) GoServe(ctx context.Context, in *Serve) error {
	fmt.Println("starting the web routine!", in.Listen)
	return nil
}
