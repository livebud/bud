package bud

import (
	"context"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/livebud/bud/internal/buildcache"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/overlay"
)

func New(fsys *overlay.FileSystem, module *gomod.Module) *Project {
	return &Project{
		fsys:   fsys,
		module: module,
		bcache: buildcache.Default(),
		Env:    os.Environ(),
		Stderr: os.Stderr,
		Stdout: os.Stdout,
	}
}

type Project struct {
	fsys   *overlay.FileSystem
	module *gomod.Module
	bcache *buildcache.Cache
	Env    []string
	Stdout io.Writer
	Stderr io.Writer
}

func (c *Project) Compile(ctx context.Context, flag *Flag) (*App, error) {
	// Sync the app
	if err := c.fsys.Sync("bud/.app"); err != nil {
		return nil, err
	}
	// Ensure that main.go exists
	if _, err := fs.Stat(c.module, "bud/.app/main.go"); err != nil {
		return nil, err
	}
	// Build the binary
	if err := c.bcache.Build(ctx, c.module, "bud/.app/main.go", filepath.Join("bud", "app")); err != nil {
		return nil, err
	}
	return &App{
		Module: c.module,
		Env:    c.Env,
		Stderr: c.Stderr,
		Stdout: c.Stdout,
	}, nil
}
