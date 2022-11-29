package generate

import (
	"context"
	"strings"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/bfs"
	"github.com/livebud/bud/internal/cli/bud"
	"github.com/livebud/bud/internal/versions"
)

// New command for bud generate
func New(bud *bud.Command, in *bud.Input) *Command {
	return &Command{
		bud: bud,
		in:  in,
		Flag: &framework.Flag{
			Env:    in.Env,
			Stderr: in.Stderr,
			Stdin:  in.Stdin,
			Stdout: in.Stdout,
		},
	}
}

// Command for running bud generate
type Command struct {
	bud  *bud.Command
	in   *bud.Input
	Flag *framework.Flag
	Args []string
}

// Run the generate command
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
	// Load the filesystem
	bfs, err := bfs.Load(c.Flag, log, module)
	if err != nil {
		return err
	}
	defer bfs.Close()
	// Sync either the entire bud directory or the specified files
	return bfs.Sync(selectBudDirs(c.Args)...)
}

func selectBudDirs(patterns []string) (paths []string) {
	for _, pattern := range patterns {
		// Only sync from within the bud directory
		if !strings.HasPrefix(pattern, "bud/") {
			continue
		}
		// Trim the wildcard suffix since SyncDirs is recursive already
		// TODO: support non-recursive syncs
		pattern = strings.TrimSuffix(pattern, "/...")
		// Add the file or directory
		paths = append(paths, pattern)
	}
	return paths
}
