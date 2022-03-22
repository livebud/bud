package bud

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"

	"gitlab.com/mnm/bud/generator/cli/command"
	"gitlab.com/mnm/bud/generator/cli/generator"
	"gitlab.com/mnm/bud/generator/cli/mainfile"
	"gitlab.com/mnm/bud/generator/cli/program"
	"gitlab.com/mnm/bud/internal/dirhash"
	"gitlab.com/mnm/bud/internal/dsync"
	"gitlab.com/mnm/bud/internal/gitignore"
	"gitlab.com/mnm/bud/package/overlay"
	"gitlab.com/mnm/bud/package/trace"
	"gitlab.com/mnm/bud/pkg/di"
	"gitlab.com/mnm/bud/pkg/gomod"
	"gitlab.com/mnm/bud/pkg/parser"
)

func defaultEnv(module *gomod.Module) Env {
	return Env{
		"HOME":       os.Getenv("HOME"),
		"PATH":       os.Getenv("PATH"),
		"GOPATH":     os.Getenv("GOPATH"),
		"GOMODCACHE": module.ModCache(),
		"TMPDIR":     os.TempDir(),
	}
}

func Find(dir string) (*Compiler, error) {
	module, err := gomod.Find(dir)
	if err != nil {
		return nil, err
	}
	return &Compiler{
		module: module,
		Env:    defaultEnv(module),
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}, nil
}

func New(module *gomod.Module) *Compiler {
	return &Compiler{
		module: module,
		Env:    defaultEnv(module),
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

type Compiler struct {
	module     *gomod.Module
	Env        Env
	Stdout     io.Writer
	Stderr     io.Writer
	ModCacheRW bool
}

func (c *Compiler) cachePath(module *gomod.Module) (string, error) {
	hash, err := dirhash.Hash(module, dirhash.WithSkip(gitignore.FromFS(module)))
	if err != nil {
		return "", err
	}
	return filepath.Join(os.TempDir(), "bud-compiler", hash), nil
}

// Load the overlay
func (c *Compiler) loadOverlay(ctx context.Context, module *gomod.Module) (fsys *overlay.FileSystem, err error) {
	_, span := trace.Start(ctx, "load the overlay")
	defer span.End(&err)
	return overlay.Load(module)
}

// Sync the generators to bud/.cli
func (c *Compiler) sync(ctx context.Context, overlay *overlay.FileSystem) (err error) {
	_, span := trace.Start(ctx, "sync cli", "dir", "bud/.cli")
	defer span.End(&err)
	return dsync.Dir(overlay, "bud/.cli", c.module.DirFS("bud/.cli"), ".")
}

// Build the CLI
func (c *Compiler) goBuild(ctx context.Context, module *gomod.Module, outPath string) (err error) {
	_, span := trace.Start(ctx, "build cli", "from", "bud/.cli/main.go", "to", "bud/cli")
	defer span.End(&err)
	// Ensure that main.go exists
	if _, err := fs.Stat(c.module, "bud/.cli/main.go"); err != nil {
		return err
	}
	// Compile the args
	args := []string{
		"build",
		"-mod=mod",
		"-o=" + outPath,
	}
	if c.ModCacheRW {
		args = append(args, "-modcacherw")
	}
	args = append(args, "bud/.cli/main.go")
	// Run go build
	cmd := exec.CommandContext(ctx, "go", args...)
	cmd.Env = c.Env.List()
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr
	cmd.Dir = module.Directory()
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (c *Compiler) Compile(ctx context.Context, flag Flag) (p *Project, err error) {
	// Start the trace
	ctx, span := trace.Start(ctx, "compile project cli")
	defer span.End(&err)
	// Load the overlay
	overlay, err := c.loadOverlay(ctx, c.module)
	if err != nil {
		return nil, err
	}
	// Initialize dependencies
	parser := parser.New(overlay, c.module)
	injector := di.New(overlay, c.module, parser)
	// Setup the generators
	overlay.FileGenerator("bud/.cli/main.go", mainfile.New(c.module))
	overlay.FileGenerator("bud/.cli/program/program.go", program.New(injector, c.module))
	overlay.FileGenerator("bud/.cli/command/command.go", command.New(overlay, c.module, parser))
	overlay.FileGenerator("bud/.cli/generator/generator.go", generator.New(overlay, c.module, parser))
	// Sync the generators
	if err := c.sync(ctx, overlay); err != nil {
		return nil, err
	}
	cachePath := ""
	cliPath := filepath.Join("bud", "cli")
	// Cached build
	if flag.Cache {
		cachePath, err = c.cachePath(c.module)
		if err != nil {
			return nil, err
		}
		cliPath = filepath.Join(cachePath, "cli")
		if _, err := os.Stat(cliPath); errors.Is(err, fs.ErrNotExist) {
			// Build the binary
			if err := c.goBuild(ctx, c.module, cliPath); err != nil {
				return nil, err
			}
		} else if err != nil {
			return nil, err
		}
	} else {
		// Build the binary
		if err := c.goBuild(ctx, c.module, cliPath); err != nil {
			return nil, err
		}
	}
	return &Project{
		Path:   cliPath,
		Cache:  cachePath,
		Module: c.module,
		Flag:   flag,
		Env:    c.Env,
		Stdout: c.Stdout,
		Stderr: c.Stderr,
	}, nil
}
