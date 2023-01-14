package cli

import (
	"context"

	"github.com/livebud/bud"
	"github.com/livebud/bud/internal/dag"
	"github.com/livebud/bud/internal/once"
	"github.com/livebud/bud/internal/sh"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/remotefs"
	"github.com/livebud/bud/package/socket"
)

type budFS struct {
	afsLn  socket.Listener
	cli    *CLI
	db     *dag.DB
	devLn  socket.Listener
	log    log.Log
	module *gomod.Module
	closer once.Closer
}

func (f *budFS) Sync(ctx context.Context, cfg *bud.Config) error {
	cmd := &sh.Command{
		Dir:    f.cli.Dir,
		Env:    f.cli.Env,
		Stdin:  f.cli.Stdin,
		Stdout: f.cli.Stdout,
		Stderr: f.cli.Stderr,
	}

	// Start the dev server
	f.log.Debug("run: dev server is listening on http://" + f.devLn.Addr().String())
	devServer, err := f.cli.startDev(ctx, cfg, f.devLn, f.log)
	if err != nil {
		return err
	}
	f.closer.Add(devServer.Close)

	// Generate the AFS
	genfs := f.cli.genfs(cfg, f.db, f.log, f.module)
	if err := f.cli.generateAFS(ctx, cmd, genfs, f.log, f.module); err != nil {
		return err
	}

	// Reset the cache
	// TODO: optimize later
	if err := f.db.Reset(); err != nil {
		return err
	}

	// Load the *os.File for afsLn
	afsFile, err := f.afsLn.File()
	if err != nil {
		return err
	}
	f.closer.Add(afsFile.Close)

	// Remote client
	remoteClient, err := remotefs.Dial(ctx, f.afsLn.Addr().String())
	if err != nil {
		return err
	}
	f.closer.Add(remoteClient.Close)

	// Inject the file under the AFS prefix
	cmd.Inject("AFS", afsFile)
	// Start the application file server
	cmd.Env = append(cmd.Env, "BUD_DEV_URL="+f.devLn.Addr().String())
	afsProcess, err := f.cli.startAFS(ctx, cmd, f.module)
	if err != nil {
		return err
	}
	f.closer.Add(afsProcess.Close)

	// Build the application
	if err := f.cli.buildApp(ctx, remoteClient, cmd, f.log, f.module); err != nil {
		return err
	}

	return nil
}

func (f *budFS) Refresh(ctx context.Context, paths ...string) error {
	// TODO:  optimize later
	return f.db.Reset()
}

func (f *budFS) Close() error {
	return f.closer.Close()
}
