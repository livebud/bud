package view

import (
	"context"
	"fmt"
)

type Command struct {
	Name     string `arg:"name" help:"name of the view"`
	WithTest bool   `flag:"with-test" help:"include a view test" default:"true"`
}

func (c *Command) View(ctx context.Context) error {
	fmt.Println("creating new view", c.Name, c.WithTest)
	return nil
}
