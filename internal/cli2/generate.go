package cli

import (
	"context"
	"errors"
	"io/fs"
	"path/filepath"

	"github.com/livebud/bud"
	"github.com/livebud/bud/internal/dsync"
	"github.com/livebud/bud/internal/versions"
	"github.com/livebud/bud/package/remotefs"
	"golang.org/x/sync/errgroup"
)

func (c *CLI) Generate(ctx context.Context, in *bud.Generate) error {
	cfg := &in.Config
	cmd := c.Command.Clone()

	// Find go.mod
	module, err := c.module()
	if err != nil {
		return err
	}
	cmd.Env = append(cmd.Env, "GOMODCACHE="+module.ModCache())

	// Align the runtime
	if err := versions.AlignRuntime(ctx, module, versions.Bud); err != nil {
		return err
	}

	// Setup the logger
	log, err := c.logger()
	if err != nil {
		return err
	}

	// Load the database
	db, err := c.openDatabase(log, module)
	if err != nil {
		return err
	}
	defer db.Close()

	// Listen on the dev address
	{
		devLn, err := c.listenDev(":0")
		if err != nil {
			return err
		}
		defer devLn.Close()
		log.Debug("run: dev server is listening on http://" + devLn.Addr().String())
		eg, ctx := errgroup.WithContext(ctx)
		// Start the dev server
		eg.Go(func() error {
			err := c.serveDev(ctx, cfg, devLn, log)
			return err
		})
		cmd.Env = append(cmd.Env, "BUD_DEV_URL="+devLn.Addr().String())
	}

	skips := []func(name string, isDir bool) bool{}
	{
		// Setup genfs
		genfs := c.genfs(cfg, db, log, module)
		// Generate AFS
		// Skip hidden files and directories
		skips = append(skips, func(name string, isDir bool) bool {
			base := filepath.Base(name)
			return base[0] == '_' || base[0] == '.'
		})
		// Skip files we want to carry over
		skips = append(skips, func(name string, isDir bool) bool {
			switch name {
			case "bud/bud.db", "bud/afs", "bud/app":
				return true
			default:
				return false
			}
		})
		if err := dsync.To(genfs, module, "bud", dsync.WithSkip(skips...), dsync.WithLog(log)); err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				return err
			}
		}
		// Build the afs binary
		if err := cmd.Run(ctx, "go", "build", "-mod=mod", "-o="+afsBinPath, afsMainPath); err != nil {
			return err
		}
	}

	// Reset the cache
	// TODO: optimize later
	if err := db.Reset(); err != nil {
		return err
	}

	// Start the file server listener
	afsLn, err := c.listenAFS(":0")
	if err != nil {
		return err
	}
	defer afsLn.Close()

	{
		cmd.Env = append(cmd.Env, "BUD_AFS_URL="+afsLn.Addr().String())
		log.Debug("run: afs server is listening on http://" + afsLn.Addr().String())
		// Load the *os.File for afsLn
		afsFile, err := afsLn.File()
		if err != nil {
			return err
		}
		defer afsFile.Close()
		// Inject the file under the AFS prefix
		cmd.Inject("AFS", afsFile)
		// Start afs
		afsProcess, err := cmd.Start(ctx, module.Directory("bud", "afs"))
		if err != nil {
			return err
		}
		defer afsProcess.Close()
	}

	{
		// Remote client
		remoteClient, err := remotefs.Dial(ctx, afsLn.Addr().String())
		if err != nil {
			return err
		}
		defer remoteClient.Close()
		// Generate the app
		// Skip over the afs files we just generated
		skips = append(skips, func(name string, isDir bool) bool {
			return isAFSPath(name)
		})
		// Sync the app files again with the remote filesystem
		if err := dsync.To(remoteClient, module, "bud", dsync.WithSkip(skips...), dsync.WithLog(log)); err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				return err
			}
		}
		// Build the application binary
		if err := cmd.Run(ctx, "go", "build", "-mod=mod", "-o="+appBinPath, appMainPath); err != nil {
			return err
		}
	}

	return nil
}
