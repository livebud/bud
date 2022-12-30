package generate

import (
	"context"
	"strings"

	"github.com/livebud/bud/framework"
	budfs "github.com/livebud/bud/internal/budfs"
	"github.com/livebud/bud/internal/cli/bud"
	"github.com/livebud/bud/internal/versions"
)

// New command for bud generate
func New(bud *bud.Command, in *bud.Input) *Command {
	return &Command{
		bud:  bud,
		in:   in,
		Flag: new(framework.Flag),
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
	// Setup the listener
	budln, err := bud.BudListener(c.in)
	if err != nil {
		return err
	}
	defer budln.Close()
	// Setup the command shell
	cmd := bud.Shell(c.in, module)
	cmd.Env = append(cmd.Env, "BUD_LISTEN="+budln.Addr().String())
	// Load the budfs
	bfs, err := budfs.Load(cmd, c.Flag, module, log)
	if err != nil {
		return err
	}
	defer bfs.Close(ctx)
	// Start the server
	budServer, err := bud.StartBudServer(ctx, budln, bfs, log)
	if err != nil {
		return err
	}
	defer budServer.Close()
	// Sync the directories
	if err := bfs.Sync(ctx, module, selectBudDirs(c.Args)...); err != nil {
		return err
	}
	return nil
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
