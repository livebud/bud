package expand

import (
	"context"
	"fmt"
	"io/fs"

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
)

type Command struct {
	Module *gomod.Module
	// TODO: support passing macros in
}

func (c *Command) Expand(ctx context.Context) error {
	if c.Module == nil {
		return fmt.Errorf("expand: requires a module")
	}

	// Initialize generator dependencies
	genfs, err := overlay.Load(c.Module)
	if err != nil {
		return err
	}
	parser := parser.New(genfs, c.Module)
	injector := di.New(genfs, c.Module, parser)

	// Setup the generators
	genfs.FileGenerator("bud/import.go", importfile.New(c.Module))
	genfs.FileGenerator("bud/.cli/main.go", mainfile.New(c.Module))
	genfs.FileGenerator("bud/.cli/program/program.go", program.New(injector, c.Module))
	genfs.FileGenerator("bud/.cli/command/command.go", command.New(injector, c.Module, parser))
	genfs.FileGenerator("bud/.cli/generator/generator.go", generator.New(genfs, c.Module, parser))
	genfs.FileGenerator("bud/.cli/transform/transform.go", transform.New(c.Module))

	// Synchronize the generators with the filesystem
	if err := genfs.Sync("bud/.cli"); err != nil {
		return err
	}

	// Write the import generator
	// TODO: add import writer back in
	// if err := c.writeImporter(ctx, overlay); err != nil {
	// 	return err
	// }

	// Ensure that main.go exists
	if _, err := fs.Stat(c.Module, "bud/.cli/main.go"); err != nil {
		return err
	}

	// Run go build on bud/.cli/main.go outputing to bud/cli
	bcache := buildcache.Default(c.Module)
	if err := bcache.Build(ctx, c.Module, "bud/.cli/main.go", "bud/cli"); err != nil {
		return err
	}

	return nil
}
