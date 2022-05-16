// Package buildcache is a build cache
package buildcache

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/livebud/bud/internal/gobin"
	"github.com/livebud/bud/internal/imhash"
	"github.com/livebud/bud/internal/symlink"
	"github.com/livebud/bud/package/gomod"
)

func Default(module *gomod.Module) *Cache {
	return &Cache{
		Dir: filepath.Join(module.Directory(), "bud", "cache"),
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
