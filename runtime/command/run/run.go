package run

import (
	"context"
	"fmt"
	"time"

	"gitlab.com/mnm/bud/package/watcher"

	"gitlab.com/mnm/bud/package/hot"

	"golang.org/x/sync/errgroup"

	"gitlab.com/mnm/bud/package/socket"
	"gitlab.com/mnm/bud/runtime/project"
)

type Command struct {
	Project *project.Compiler
	Flag    project.Flag
	Port    string
}

func (c *Command) Run(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)
	// Start the web server
	eg.Go(func() error { return c.startApp(ctx) })
	// Start the hot reload server
	if c.Flag.Hot {
		hotServer := hot.New()
		eg.Go(func() error { return c.startWatch(ctx, hotServer) })
		eg.Go(func() error { return c.startHot(ctx, hotServer) })
	}
	return eg.Wait()
}

func (c *Command) startApp(ctx context.Context) error {
	now := time.Now()
	fmt.Println("compiling app", time.Since(now))
	app, err := c.Project.Compile(ctx, &c.Flag)
	if err != nil {
		return err
	}
	fmt.Println("compiled app", time.Since(now))
	listener, err := socket.Load(c.Port)
	if err != nil {
		return err
	}
	now = time.Now()
	fmt.Println("starting app")
	process, err := app.Start(ctx, listener)
	if err != nil {
		return err
	}
	fmt.Println("started app", time.Since(now))
	return process.Wait()
}

func (c *Command) startHot(ctx context.Context, hotServer *hot.Server) error {
	return hotServer.ListenAndServe(":35729")
}

func (c *Command) startWatch(ctx context.Context, hotServer *hot.Server) error {
	return watcher.Watch(ctx, ".", func(path string) error {
		fmt.Println("path changed", path)
		hotServer.Reload("bud/view/index.svelte")
		return nil
	})
}
