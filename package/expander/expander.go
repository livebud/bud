package expander

import (
	"context"
	"io"
	"io/fs"
	"os"
	"os/exec"

	"gitlab.com/mnm/bud/generator2/cli/command"
	"gitlab.com/mnm/bud/generator2/cli/generator"
	"gitlab.com/mnm/bud/generator2/cli/mainfile"
	"gitlab.com/mnm/bud/generator2/cli/program"
	"gitlab.com/mnm/bud/internal/dsync"
	"gitlab.com/mnm/bud/internal/gobin"
	"gitlab.com/mnm/bud/package/overlay"
	"gitlab.com/mnm/bud/pkg/di"
	"gitlab.com/mnm/bud/pkg/gomod"
	"gitlab.com/mnm/bud/pkg/parser"
)

func Load(dir string) (*Expander, error) {
	module, err := gomod.Find(dir)
	if err != nil {
		return nil, err
	}
	ofs, err := overlay.Load(module)
	if err != nil {
		return nil, err
	}
	parser := parser.New(ofs, module)
	injector := di.New(ofs, module, parser)
	ofs.FileGenerator("bud/.cli/main.go", mainfile.New(module))
	ofs.FileGenerator("bud/.cli/program/program.go", program.New(injector, module))
	ofs.FileGenerator("bud/.cli/command/command.go", command.New(module))
	ofs.FileGenerator("bud/.cli/generator/generator.go", generator.New())
	if err := dsync.Dir(ofs, ".", module.DirFS("."), "."); err != nil {
		return nil, err
	}
	return &Expander{
		dir:     dir,
		overlay: ofs,
		module:  module,
		Env:     os.Environ(),
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
	}, nil
}

type Expander struct {
	dir     string
	overlay fs.FS
	module  *gomod.Module
	Env     []string
	Stdout  io.Writer
	Stderr  io.Writer
}

func (e *Expander) Expand(ctx context.Context, args ...string) error {
	// Generate the CLI
	if err := dsync.Dir(e.overlay, ".", e.module.DirFS("."), "."); err != nil {
		return err
	}
	// Build the CLI
	if err := gobin.Build(ctx, e.dir, "bud/.cli/main.go", "bud/cli"); err != nil {
		return err
	}
	// Run the CLI
	cmd := exec.CommandContext(ctx, "./bud/cli", args...)
	cmd.Dir = e.dir
	cmd.Stdout = e.Stdout
	cmd.Stderr = e.Stderr
	cmd.Env = e.Env
	return cmd.Run()
}
