package runrt

import (
	"context"
	"io/fs"
	"os"

	"github.com/livebud/bud/package/exe"

	"github.com/livebud/bud/package/socket"

	"github.com/livebud/bud/internal/buildcache"
	"github.com/livebud/bud/package/hot"
	"github.com/livebud/bud/package/overlay"
	"github.com/livebud/bud/runtime/web"
	"golang.org/x/sync/errgroup"

	"github.com/livebud/bud/package/bud"
	"github.com/livebud/bud/package/gomod"
)

type Command struct {
	Module *gomod.Module       // TODO: replace with dir
	FS     *overlay.FileSystem // TODO: construct from list of generators
	Flag   *bud.Flag
	Listen string
}

func (c *Command) Run(ctx context.Context) error {
	// Load the hot server
	hotServer := hot.New()

	// Start up the processes concurrently
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error { return c.startHot(ctx, hotServer) })
	eg.Go(func() error { return c.startApp(ctx, hotServer) })
	return eg.Wait()

	// eg, ctx := errgroup.WithContext(ctx)
	// Initialize the hot server
	// var hotServer *hot.Server
	// if c.Flag.Hot {
	// 	hotServer = hot.New()
	// 	// Start the hot reload server
	// 	eg.Go(func() error { return c.startHot(ctx, hotServer) })
	// }
	// // Start the web server
	// eg.Go(func() error { return c.startApp(ctx, hotServer) })
	// return eg.Wait()
}

func (c *Command) compile(ctx context.Context) error {
	// Sync the application
	if err := c.FS.Sync("bud/.app"); err != nil {
		return err
	}
	// Ensure that main.go exists
	if _, err := fs.Stat(c.Module, "bud/.app/main.go"); err != nil {
		return err
	}
	// Build the application binary
	bcache := buildcache.Default(c.Module)
	if err := bcache.Build(ctx, c.Module, "bud/.app/main.go", "bud/app"); err != nil {
		return err
	}
	return nil
}

func (c *Command) startHot(ctx context.Context, hotServer *hot.Server) error {
	if c.Flag.Hot == "" {
		return nil
	}
	listener, err := socket.Listen(c.Flag.Hot)
	if err != nil {
		return err
	}
	return web.Serve(ctx, listener, hotServer)
}

func (c *Command) startApp(ctx context.Context, hotServer *hot.Server) error {
	if err := c.compile(ctx); err != nil {
		return err
	}
	cmd := exe.Command(ctx, "bud/app")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	cmd.Dir = c.Module.Directory()
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
