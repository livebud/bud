package cache

import (
	"context"
	"os"
	"path/filepath"
)

type Command struct {
}

func (c *Command) Clean(ctx context.Context) error {
	cacheDir := filepath.Join(os.TempDir(), "bud-compiler")
	return os.RemoveAll(cacheDir)
}
