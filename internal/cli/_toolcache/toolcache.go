package toolcache

import (
	"context"
	"os"

	"github.com/livebud/bud/internal/cli/bud"
)

func New(bud *bud.Command) *Command {
	return &Command{bud: bud}
}

type Command struct {
	bud *bud.Command
}

func (c *Command) Run(ctx context.Context) error {
	module, err := c.bud.Module()
	if err != nil {
		return err
	}
	return os.RemoveAll(module.Directory("bud", "cache"))
}
