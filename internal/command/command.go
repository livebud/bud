package command

import (
	"context"
	"io"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/gobuild"
	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/parser"

	"github.com/livebud/bud/framework/app"
	"github.com/livebud/bud/framework/web"
	"github.com/livebud/bud/package/overlay"

	"github.com/livebud/bud/package/commander"

	"github.com/livebud/bud/package/log/console"
	"github.com/livebud/bud/package/log/filter"

	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log"
)

// Bud command
type Bud struct {
	// Flags
	Dir  string
	Log  string
	Args []string
	Help bool

	// Passed through the subprocesses
	Env    []string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func (c *Bud) Module() (*gomod.Module, error) {
	return gomod.Find(c.Dir)
}

func (c *Bud) Logger() (log.Interface, error) {
	handler, err := filter.Load(console.New(c.Stderr), c.Log)
	if err != nil {
		return nil, err
	}
	return log.New(handler), nil
}

// Filesystem loads the filesystem
// TODO: inline the overlay loader
func (c *Bud) FileSystem(module *gomod.Module, flag *framework.Flag) (*overlay.FileSystem, error) {
	genfs, err := overlay.Load(module)
	if err != nil {
		return nil, err
	}
	parser := parser.New(genfs, module)
	injector := di.New(genfs, module, parser)
	genfs.FileGenerator("bud/internal/app/main.go", app.New(injector, module, flag))
	genfs.FileGenerator("bud/internal/app/web/web.go", web.New(module, parser))
	// genfs.FileGenerator("bud/import.go", importfile.New(module))
	// genfs.FileGenerator("bud/.cli/main.go", mainfile.New(module))
	// genfs.FileGenerator("bud/.cli/program/program.go", program.New(injector, module))
	// genfs.FileGenerator("bud/.cli/command/command.go", command.New(injector, module, parser))
	// genfs.FileGenerator("bud/.cli/generator/generator.go", generator.New(genfs, module, parser))
	// genfs.FileGenerator("bud/.cli/transform/transform.go", transform.New(module))
	return genfs, nil
}

func (c *Bud) Builder(module *gomod.Module) *gobuild.Builder {
	return gobuild.New(module)
}

// Run a custom command
// TODO: finish supporting custom commands
// 1. Compile
//   a. Generate generator (later!)
//   	 i. Generate bud/internal/generator
//     ii. Build bud/generator
//     iii. Run bud/generator
//   b. Generate custom command
//     i. Generate bud/internal/command/${name}/
//     ii. Build bud/command/${name}
// 2. Run bud/command/${name}
func (c *Bud) Run(ctx context.Context) error {
	return commander.Usage()
}
