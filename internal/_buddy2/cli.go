package buddy

import (
	"context"
	"fmt"
	"os/exec"

	"gitlab.com/mnm/bud/pkg/gomod"
)

type CLI struct {
	module *gomod.Module
	path   string
}

func (c *CLI) Command(ctx context.Context, args ...string) *exec.Cmd {
	return nil
}

func (c *CLI) Run(ctx context.Context, options ...runOption) (*Process, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *CLI) Build(ctx context.Context, options ...buildOption) (*App, error) {
	return nil, fmt.Errorf("not implemented")
}
