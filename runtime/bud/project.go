package bud

import (
	"context"
	"fmt"
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
		bcache: buildcache.Default(module),
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
	fmt.Println("syncing app")
	// Sync the app
	if err := c.fsys.Sync("bud/.app"); err != nil {
		return nil, err
	}
	fmt.Println("synced app")
	fmt.Println("checking if main exists")
	// Ensure that main.go exists
	if _, err := fs.Stat(c.module, "bud/.app/main.go"); err != nil {
		return nil, err
	}
	fmt.Println("checked if main exists")
	// Build the binary
	fmt.Println("building app")
	if err := c.bcache.Build(ctx, c.module, "bud/.app/main.go", filepath.Join("bud", "app")); err != nil {
		return nil, err
	}
	fmt.Println("built app")
	return &App{
		Module: c.module,
		Env:    c.Env,
		Stderr: c.Stderr,
		Stdout: c.Stdout,
	}, nil
}
