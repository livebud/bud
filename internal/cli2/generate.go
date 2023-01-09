package cli

import (
	"context"

	"github.com/livebud/bud"
	"github.com/livebud/bud/internal/once"
	"github.com/livebud/bud/internal/versions"
)

func (c *CLI) Generate(ctx context.Context, in *bud.Generate) error {
	// Find go.mod
	module, err := c.module()
	if err != nil {
		return err
	}

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

	// Start the file server listener
	afsLn, err := c.listenAFS(":0")
	if err != nil {
		return err
	}
	defer afsLn.Close()
	log.Debug("run: afs server is listening on http://" + afsLn.Addr().String())

	// Listen on the dev address
	devLn, err := c.listenDev(in.DevAddress)
	if err != nil {
		return err
	}
	defer devLn.Close()

	budfs := &budFS{afsLn, c, db, devLn, log, module, once.Closer{}}
	defer budfs.Close()

	// Syncing bud
	if err := budfs.Sync(ctx, &in.Config); err != nil {
		return err
	}

	// if err := c.sync(ctx, cmd, log, module, in); err != nil {
	// 	return err
	// }

	// // Load the database
	// db, err := c.openDatabase(log, module)
	// if err != nil {
	// 	return err
	// }
	// defer db.Close()

	// // Listen on the dev address
	// devLn, err := c.listenDev(in.DevAddress)
	// if err != nil {
	// 	return err
	// }
	// defer devLn.Close()
	// log.Debug("run: dev server is listening on http://" + devLn.Addr().String())
	// // Start the dev server
	// eg, ctx := errgroup.WithContext(ctx)
	// eg.Go(func() error { return c.serveDev(ctx, cfg, devLn, log) })

	// // Generate the AFS
	// genfs := c.genfs(cfg, db, log, module)
	// if err := c.generateAFS(ctx, cmd, genfs, log, module); err != nil {
	// 	return err
	// }

	// // Reset the cache
	// // TODO: optimize later
	// if err := db.Reset(); err != nil {
	// 	return err
	// }

	// // Start the file server listener
	// afsLn, err := c.listenAFS(":0")
	// if err != nil {
	// 	return err
	// }
	// defer afsLn.Close()
	// log.Debug("run: afs server is listening on http://" + afsLn.Addr().String())

	// // Load the *os.File for afsLn
	// afsFile, err := afsLn.File()
	// if err != nil {
	// 	return err
	// }
	// defer afsFile.Close()

	// // Remote client
	// remoteClient, err := remotefs.Dial(ctx, afsLn.Addr().String())
	// if err != nil {
	// 	return err
	// }
	// defer remoteClient.Close()

	// // Inject the file under the AFS prefix
	// cmd.Inject("AFS", afsFile)
	// // Start the application file server
	// cmd.Env = append(cmd.Env, "BUD_DEV_URL="+devLn.Addr().String())
	// afsProcess, err := c.startAFS(ctx, cmd, module)
	// if err != nil {
	// 	return err
	// }
	// defer afsProcess.Close()

	// // Build the application
	// if err := c.buildApp(ctx, remoteClient, cmd, log, module); err != nil {
	// 	return err
	// }

	return nil
}
