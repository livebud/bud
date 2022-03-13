package command

import "context"

// CommandExample example command
type CommandExample struct {
	// Dependencies ...
}

type Up struct {
	B  bool   // defaults to lowercase required flag
	A  bool   `arg:"a" help:"..." default:"..." hidden:"true"`
	CC string `short:"c"`
}

func (c *CommandExample) Up(ctx context.Context, in *Up) error {
	return nil
}
