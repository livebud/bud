package bud

import (
	"context"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"

	"gitlab.com/mnm/bud/internal/dsync"
	"gitlab.com/mnm/bud/internal/generator/command"
	"gitlab.com/mnm/bud/internal/generator/generator"
	"gitlab.com/mnm/bud/internal/generator/importfile"
	"gitlab.com/mnm/bud/internal/generator/mainfile"
	"gitlab.com/mnm/bud/internal/generator/program"
	"gitlab.com/mnm/bud/package/di"
	"gitlab.com/mnm/bud/package/gomod"
	"gitlab.com/mnm/bud/package/overlay"
	"gitlab.com/mnm/bud/package/parser"
	"gitlab.com/mnm/bud/package/trace"
	"gitlab.com/mnm/bud/runtime/bud"
)

func defaultEnv(module *gomod.Module) (Env, error) {
	budPath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	return Env{
		"HOME":       os.Getenv("HOME"),
		"PATH":       os.Getenv("PATH"),
		"GOPATH":     os.Getenv("GOPATH"),
		"GOMODCACHE": module.ModCache(),
		"TMPDIR":     os.TempDir(),
		"BUD_PATH":   budPath,
	}, nil
}

func Find(dir string) (*Compiler, error) {
	module, err := gomod.Find(dir)
	if err != nil {
		return nil, err
	}
	return Load(module)
}

func Load(module *gomod.Module) (*Compiler, error) {
	env, err := defaultEnv(module)
	if err != nil {
		return nil, err
	}
	return &Compiler{
		module: module,
		Env:    env,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}, nil
}

type Compiler struct {
	module     *gomod.Module
	Env        Env
	Stdout     io.Writer
	Stderr     io.Writer
	ModCacheRW bool
}

// Load the overlay
func (c *Compiler) loadOverlay(ctx context.Context, module *gomod.Module) (fsys *overlay.FileSystem, err error) {
	_, span := trace.Start(ctx, "load the overlay")
	defer span.End(&err)
	return overlay.Load(module)
}

func (c *Compiler) writeImporter(ctx context.Context, overlay *overlay.FileSystem) error {
	importFile, err := fs.ReadFile(overlay, "bud/import.go")
	if err != nil {
		return err
	}
	if err := c.module.DirFS().WriteFile("bud/import.go", importFile, 0644); err != nil {
		return err
	}
	return nil
}

// Sync the generators to bud/.cli
func (c *Compiler) sync(ctx context.Context, overlay *overlay.FileSystem) (err error) {
	_, span := trace.Start(ctx, "sync cli", "dir", "bud/.cli")
	defer span.End(&err)
	if err := dsync.Dir(overlay, "bud/.cli", c.module.DirFS("bud/.cli"), "."); err != nil {
		return err
	}
	return nil
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

func (c *Compiler) Compile(ctx context.Context, flag *bud.Flag) (p *Project, err error) {
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
	overlay.FileGenerator("bud/import.go", importfile.New(c.module))
	overlay.FileGenerator("bud/.cli/main.go", mainfile.New(c.module))
	overlay.FileGenerator("bud/.cli/program/program.go", program.New(flag, injector, c.module))
	overlay.FileGenerator("bud/.cli/command/command.go", command.New(overlay, c.module, parser))
	overlay.FileGenerator("bud/.cli/generator/generator.go", generator.New(overlay, c.module, parser))
	// Sync the generators
	if err := c.sync(ctx, overlay); err != nil {
		return nil, err
	}
	// Write the import generator
	if err := c.writeImporter(ctx, overlay); err != nil {
		return nil, err
	}
	// Build the binary
	if err := c.goBuild(ctx, c.module, filepath.Join("bud", "cli")); err != nil {
		return nil, err
	}
	return &Project{
		Module: c.module,
		Env:    c.Env,
		Stdout: c.Stdout,
		Stderr: c.Stderr,
	}, nil
}
