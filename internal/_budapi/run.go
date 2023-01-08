package budapi

import (
	"context"

	"golang.org/x/sync/errgroup"
)

func (b *Bud) Run(ctx context.Context, in *RunInput) error {
	// Find go.mod
	module, err := b.findModule()
	if err != nil {
		return err
	}

	// Load the logger
	log, err := b.loadLog(b.Stderr, b.Log)
	if err != nil {
		return err
	}

	// Setup the web listener
	webLn, err := b.listenWeb(in.ListenWeb)
	if err != nil {
		return err
	}
	defer webLn.Close()
	log.Info("Listening on http://" + webLn.Addr().String())

	// Open the bud database
	db, err := b.openDb(log, module)
	if err != nil {
		return err
	}
	defer db.Close()

	// Setup the bud listener
	devLn, err := b.listenDev(in.ListenDev)
	if err != nil {
		return err
	}
	defer devLn.Close()
	log.Debug("run: bud server is listening on %s", "http://"+devLn.Addr().String())

	eg, ctx := errgroup.WithContext(ctx)

	// Start the dev server
	eg.Go(func() error { return b.serveDev(ctx, devLn) })

	// Generate the afs
	if err := b.generateAFS(ctx); err != nil {
		return err
	}

	// Wait until run finishes
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}
