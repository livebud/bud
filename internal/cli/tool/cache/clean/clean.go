package clean

import (
	"context"
	"os"
	"path/filepath"

	"github.com/livebud/bud/package/gomod"
)

type Command struct {
	Dir string
}

func (c *Command) Run(ctx context.Context) error {
	module, err := gomod.Find(c.Dir)
	if err != nil {
		return err
	}
	cacheDir := filepath.Join(module.Directory(), "bud", "cache")
	return os.RemoveAll(cacheDir)
}
