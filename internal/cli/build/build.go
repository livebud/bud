package build

import (
	"context"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/cli/bud"
	"github.com/livebud/bud/internal/dsync"
	"github.com/livebud/bud/internal/gobuild"
	"github.com/livebud/bud/internal/versions"
	"github.com/livebud/bud/package/goplugin"
	"github.com/livebud/bud/package/remotefs"
	"github.com/livebud/bud/package/vfs"
)

// New command for bud build
func New(bud *bud.Command, in *bud.Input) *Command {
	return &Command{
		bud:  bud,
		in:   in,
		Flag: new(framework.Flag),
	}
}

// Command for running bud build
type Command struct {
	bud  *bud.Command
	in   *bud.Input
	Flag *framework.Flag
}

// Run the build command
func (c *Command) Run(ctx context.Context) error {
	// Find go.mod
	module, err := bud.Module(c.bud.Dir)
	if err != nil {
		return err
	}
	// Ensure we have version alignment between the CLI and the runtime
	if err := bud.EnsureVersionAlignment(ctx, module, versions.Bud); err != nil {
		return err
	}
	// Setup the logger
	log, err := bud.Log(c.in.Stderr, c.bud.Log)
	if err != nil {
		return err
	}
	genfs, err := bud.FileSystem(log, module, c.Flag)
	if err != nil {
		return err
	}
	// Sync generate now to support custom generators, if any
	if err := genfs.Sync("bud/internal/generate"); err != nil {
		return err
	}
	if err := vfs.Exist(module, "bud/internal/generate/main.go"); nil == err {
		conn, err := goplugin.Start(module.Directory(), "go", "run", "-mod=mod", "bud/internal/generate/main.go")
		if err != nil {
			return err
		}
		remotefs := remotefs.NewClient(conn)
		defer remotefs.Close()
		if err := dsync.Dir(remotefs, "bud/internal/generator", module.DirFS("bud/internal/generator"), "."); err != nil {
			return err
		}
	}
	if err := genfs.Sync("bud/internal/app"); err != nil {
		return err
	}
	builder := gobuild.New(module)
	return builder.Build(ctx, "bud/internal/app/main.go", "bud/app")
}
