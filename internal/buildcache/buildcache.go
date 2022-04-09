// Package buildcache is a build cache
package buildcache

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"gitlab.com/mnm/bud/internal/gobin"
	"gitlab.com/mnm/bud/internal/imhash"
	"gitlab.com/mnm/bud/internal/symlink"
	"gitlab.com/mnm/bud/package/gomod"
)

func Default() *Cache {
	return &Cache{
		// TODO: make this configurable
		// TODO: use the user cache, once we have a way to clean up
		Dir: filepath.Join(os.TempDir(), "bud", "cache"),
	}
}

type Cache struct {
	Dir string
}

func (c *Cache) exists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

var _ gobin.Builder = (*Cache)(nil)

func (c *Cache) Build(ctx context.Context, module *gomod.Module, mainPath string, outPath string, flags ...string) error {
	hash, err := imhash.Hash(module, filepath.Dir(mainPath))
	if err != nil {
		return err
	}
	cachePath := filepath.Join(c.Dir, hash)
	exists, err := c.exists(cachePath)
	if err != nil {
		return err
	} else if exists {
		return symlink.Link(cachePath, module.Directory(outPath))
	}
	if err := gobin.Build(ctx, module, mainPath, cachePath, flags...); err != nil {
		return err
	}
	return symlink.Link(cachePath, module.Directory(outPath))
}
