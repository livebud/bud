package project

import (
	"context"
	"io"
	"io/fs"
	"os"

	"gitlab.com/mnm/bud/internal/bud"
	"gitlab.com/mnm/bud/internal/dsync"
	"gitlab.com/mnm/bud/internal/gobin"
	"gitlab.com/mnm/bud/pkg/gomod"
)

func New(fsys fs.FS, module *gomod.Module) *Compiler {
	return &Compiler{
		fsys:   fsys,
		module: module,
		Env:    os.Environ(),
		Stderr: os.Stderr,
		Stdout: os.Stdout,
	}
}

type Compiler struct {
	fsys   fs.FS
	module *gomod.Module
	Env    []string
	Stdout io.Writer
	Stderr io.Writer
}

func (c *Compiler) Compile(ctx context.Context) (*bud.App, error) {
	if err := dsync.Dir(c.fsys, "bud/.app", c.module.DirFS("bud/.app"), "."); err != nil {
		return nil, err
	}
	// Ensure that main.go exists
	if _, err := fs.Stat(c.module, "bud/.app/main.go"); err != nil {
		return nil, err
	}
	if err := gobin.Build(ctx, c.module.Directory(), "bud/.app/main.go", "bud/app"); err != nil {
		return nil, err
	}
	return &bud.App{
		Module: c.module,
		Env:    c.Env,
		Stderr: c.Stderr,
		Stdout: c.Stdout,
	}, nil
}
