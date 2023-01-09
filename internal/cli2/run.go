package cli

import (
	"context"
	"errors"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/livebud/bud/internal/once"
	"github.com/livebud/bud/internal/versions"

	"github.com/livebud/bud/internal/prompter"
	"github.com/livebud/bud/package/watcher"

	"github.com/livebud/bud"
)

func (c *CLI) Run(ctx context.Context, in *bud.Run) error {
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

	// Listen on the web address
	webLn, err := c.listenWeb(in.WebAddress)
	if err != nil {
		return err
	}
	defer webLn.Close()
	log.Info("Listening on http://" + webLn.Addr().String())

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

	// Get the file descriptor for the web listener
	webFile, err := webLn.File()
	if err != nil {
		return err
	}
	defer webFile.Close()

	// Inject that file under the WEB prefix
	cmd.Inject("WEB", webFile)

	// Start the app
	cmd.Env = append(cmd.Env, "BUD_AFS_URL="+afsLn.Addr().String())
	if err := expectEnv(cmd.Env, "BUD_AFS_URL"); err != nil {
		return err
	}
	cmd.Env = append(cmd.Env, "BUD_DEV_URL="+devLn.Addr().String())
	if err := expectEnv(cmd.Env, "BUD_DEV_URL"); err != nil {
		return err
	}
	appProcess, err := cmd.Start(ctx, filepath.Join("bud", "app"))
	if err != nil {
		return err
	}
	defer appProcess.Close()

	if !in.Watch {
		return appProcess.Wait()
	}

	// Setup the prompter
	prompter := c.prompter(webLn)

	// Watch for changes
	err = watcher.Watch(ctx, module.Directory(), catchError(prompter, func(events []watcher.Event) error {
		// Trigger reloading
		prompter.Reloading(events)
		// Inform the bud filesystem of the changes
		changes := make([]string, len(events))
		for i, event := range events {
			log.Debug("run: file path changed %q", event.Path)
			changes[i] = event.Path
		}
		// Refresh the cache
		// TODO: optimize later
		if err := budfs.Refresh(ctx, changes...); err != nil {
			return err
		}
		// Check if we can incrementally reload
		if canIncrementallyReload(events) {
			log.Debug("run: incrementally reloading")
			// Publish the frontend:update event
			c.Bus.Publish("frontend:update", nil)
			log.Debug("run: published event %q", "frontend:update")
			// Publish the app:ready event
			c.Bus.Publish("app:ready", nil)
			log.Debug("run: published event %q", "app:ready")
			prompter.SuccessReload()
			return nil
		}
		now := time.Now()
		log.Debug("run: restarting the process")
		if err := appProcess.Close(); err != nil {
			return err
		}
		c.Bus.Publish("backend:update", nil)
		log.Debug("run: published event %q", "backend:update")
		// Generate the app
		if err := budfs.Sync(ctx, &in.Config); err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				return err
			}
		}
		// Restart the process
		p, err := appProcess.Restart(ctx)
		if err != nil {
			c.Bus.Publish("app:error", nil)
			log.Debug("run: published event %q", "app:error")
			return err
		}
		prompter.SuccessReload()
		log.Debug("restarted the process in %s", time.Since(now))
		appProcess = p
		return nil
	}))
	if err != nil {
		return err
	}

	// Close the final process. This process is most likely different than the
	// deferred process.
	if err := appProcess.Close(); err != nil {
		return err
	}

	return nil
}

// logWrap wraps the watch function in a handler that logs the error instead of
// returning the error (and canceling the watcher)
func catchError(prompter *prompter.Prompter, fn func(events []watcher.Event) error) func(events []watcher.Event) error {
	return func(events []watcher.Event) error {
		if err := fn(events); err != nil {
			prompter.FailReload(err.Error())
		}
		return nil
	}
}

// canIncrementallyReload returns true if we can incrementally reload a page
func canIncrementallyReload(events []watcher.Event) bool {
	for _, event := range events {
		if event.Op != watcher.OpUpdate || filepath.Ext(event.Path) == ".go" {
			return false
		}
	}
	return true
}
