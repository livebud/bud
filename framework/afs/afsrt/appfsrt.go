package afsrt

import (
	"context"
	"io/fs"
	"os"

	"golang.org/x/sync/errgroup"

	"github.com/livebud/bud/package/budfs"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/log/console"
	"github.com/livebud/bud/package/log/levelfilter"
	"github.com/livebud/bud/package/remotefs"
	"github.com/livebud/bud/package/virtual/vcache"

	"github.com/livebud/bud/internal/extrafile"
	"github.com/livebud/bud/package/socket"
)

func Logger(level string) (log.Log, error) {
	lvl, err := log.ParseLevel(level)
	if err != nil {
		return nil, err
	}
	logger := log.New(levelfilter.New(console.New(os.Stderr), lvl))
	return logger, nil
}

// BudFS creates a new filesystem
func BudFS(module *gomod.Module, log log.Log) *budfs.FileSystem {
	return budfs.New(vcache.New(), module, log)
}

// Serve the remote filesystem
func Serve(ctx context.Context, log log.Log, fsys fs.FS, path string) error {
	// First try to load the listener from the parent process.
	ln, err := listen(log, path)
	if err != nil {
		return err
	}
	eg, ctx := errgroup.WithContext(ctx)
	// Handle any immediate errors from remotefs
	eg.Go(func() error {
		return remotefs.Serve(fsys, ln)
	})
	// Any errors in the group will trigger ctx to be canceled, closing the
	// listener. The listener will also be closed if the outside context is
	// canceled.
	eg.Go(func() error {
		<-ctx.Done()
		return ln.Close()
	})
	// Wait for both goroutines to finish
	return eg.Wait()
}

func listen(log log.Log, path string) (socket.Listener, error) {
	files := extrafile.Load("BUD_REMOTEFS")
	if len(files) > 0 {
		log.Debug("afs: serving from BUD_REMOTEFS file listener passed in from the parent process")
		return socket.From(files[0])
	}
	ln, err := socket.ListenUp(path, 5)
	if err != nil {
		return nil, err
	}
	log.Debug("afs: serving from %s", ln.Addr())
	return ln, nil
}
