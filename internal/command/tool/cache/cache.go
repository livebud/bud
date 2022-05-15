package cache

import (
	"context"
	"os"
	"path/filepath"

	"github.com/livebud/bud/internal/command"
	"github.com/livebud/bud/package/gomod"
)

func New(bud *command.Bud) *Command {
	return &Command{bud}
}

type Command struct {
	bud *command.Bud
}

func (c *Command) Clean(ctx context.Context) error {
	module, err := gomod.Find(c.bud.Dir)
	if err != nil {
		return err
	}
	cacheDir := filepath.Join(module.Directory(), "bud", "cache")
	return os.RemoveAll(cacheDir)
}
