package run

import (
	"context"
	"fmt"
	"io/fs"
	"net"
	"os"

	"github.com/livebud/bud/internal/buildcache"

	"github.com/livebud/bud/internal/generator/command"
	"github.com/livebud/bud/internal/generator/generator"
	"github.com/livebud/bud/internal/generator/importfile"
	"github.com/livebud/bud/internal/generator/mainfile"
	"github.com/livebud/bud/internal/generator/program"
	"github.com/livebud/bud/internal/generator/transform"
	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/exe"
	"github.com/livebud/bud/package/overlay"
	"github.com/livebud/bud/package/parser"

	"github.com/livebud/bud/package/bud"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/socket"
)

type Command struct {
	Dir      string
	Listener net.Listener
	Flag     *bud.Flag
	Env      bud.Env
	Log      log.Interface

	// TODO: support passing in a list of generators
}

func (c *Command) compile(ctx context.Context) (*exe.Cmd, error) {
	// Discard logs by default
	if c.Log == nil {
		c.Log = log.Discard
	}

	// Default flags for running
	if c.Flag == nil {
		c.Flag = &bud.Flag{
			Hot:    ":35729",
			Embed:  false,
			Minify: false,
		}
	}

	// Find go.mod
	module, err := gomod.Find(c.Dir)
	if err != nil {
		return nil, err
	}

	// Initialize generator dependencies
	genfs, err := overlay.Load(module)
	if err != nil {
		return nil, err
	}
	parser := parser.New(genfs, module)
	injector := di.New(genfs, module, parser)

	// Setup the generators
	genfs.FileGenerator("bud/import.go", importfile.New(module))
	genfs.FileGenerator("bud/.cli/main.go", mainfile.New(module))
	genfs.FileGenerator("bud/.cli/program/program.go", program.New(c.Flag, injector, module))
	genfs.FileGenerator("bud/.cli/command/command.go", command.New(module, parser))
	genfs.FileGenerator("bud/.cli/generator/generator.go", generator.New(genfs, module, parser))
	genfs.FileGenerator("bud/.cli/transform/transform.go", transform.New(module))

	// Synchronize the generators with the filesystem
	if err := genfs.Sync("bud/.cli"); err != nil {
		return nil, err
	}

	// Write the import generator
	// TODO: add import writer back in
	// if err := c.writeImporter(ctx, overlay); err != nil {
	// 	return nil, err
	// }

	// Ensure that main.go exists
	if _, err := fs.Stat(module, "bud/.cli/main.go"); err != nil {
		return nil, err
	}

	// Run go build on bud/.cli/main.go outputing to bud/cli
	bcache := buildcache.Default(module)
	if err := bcache.Build(ctx, module, "bud/.cli/main.go", "bud/cli"); err != nil {
		return nil, err
	}

	// Run the project CLI
	// $ bud/cli run
	cmd := exe.Command(ctx, "bud/cli", "run")
	cmd.Dir = module.Directory()
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// Inject the listener into the command to be available in the subprocess
	if err := socket.Inject(cmd, c.Listener); err != nil {
		return nil, err
	}
	return cmd, nil
}

func (c *Command) Run(ctx context.Context) error {
	cmd, err := c.compile(ctx)
	if err != nil {
		return err
	}
	return cmd.Run()
}

func (c *Command) Start(ctx context.Context) (*exe.Cmd, error) {
	cmd, err := c.compile(ctx)
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return cmd, err
}

func formatAddress(l net.Listener) (string, error) {
	switch l.Addr().Network() {
	case "unix":
		return l.Addr().String(), nil
	default:
		host, port, err := net.SplitHostPort(l.Addr().String())
		if err != nil {
			return "", err
		}
		// https://serverfault.com/a/444557
		if host == "::" {
			host = "0.0.0.0"
		}
		return fmt.Sprintf("https://%s:%s", host, port), nil
	}
}
