package bud

import (
	"context"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/livebud/bud/internal/buildcache"
	"github.com/livebud/bud/internal/generator/command"
	"github.com/livebud/bud/internal/generator/generator"
	"github.com/livebud/bud/internal/generator/importfile"
	"github.com/livebud/bud/internal/generator/mainfile"
	"github.com/livebud/bud/internal/generator/program"
	"github.com/livebud/bud/internal/generator/transform"
	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/overlay"
	"github.com/livebud/bud/package/parser"
	"github.com/livebud/bud/runtime/bud"
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
		bcache: buildcache.Default(),
		Env:    env,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}, nil
}

type Compiler struct {
	module     *gomod.Module
	bcache     *buildcache.Cache
	Env        Env
	Stdout     io.Writer
	Stderr     io.Writer
	ModCacheRW bool
}

// Load the overlay
func (c *Compiler) loadOverlay(ctx context.Context, module *gomod.Module) (fsys *overlay.FileSystem, err error) {
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
	if err := overlay.Sync("bud/.cli"); err != nil {
		return err
	}
	return nil
}

// Build the CLI
func (c *Compiler) goBuild(ctx context.Context, module *gomod.Module, outPath string) (err error) {
	// Ensure that main.go exists
	if _, err := fs.Stat(c.module, "bud/.cli/main.go"); err != nil {
		return err
	}
	flags := []string{}
	if c.ModCacheRW {
		flags = append(flags, "-modcacherw")
	}
	if err := c.bcache.Build(ctx, module, "bud/.cli/main.go", outPath, flags...); err != nil {
		return err
	}
	return nil
}

func (c *Compiler) Compile(ctx context.Context, flag *bud.Flag) (p *Project, err error) {
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
	overlay.FileGenerator("bud/.cli/transform/transform.go", transform.New(c.module))
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
