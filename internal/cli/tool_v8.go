package cli

import (
	"context"
	"fmt"

	v8 "github.com/livebud/bud/package/js/v8"
)

type ToolV8 struct {
}

func (c *CLI) ToolV8(ctx context.Context, in *ToolV8) error {
	script, err := c.readStdin()
	if err != nil {
		return err
	}
	vm, err := v8.Load()
	if err != nil {
		return err
	}

	result, err := vm.Eval("script.js", script)
	if err != nil {
		return err
	}
	fmt.Fprintln(c.Stdout, result)
	return nil
}
