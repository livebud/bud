package cli

import (
	"context"
	"io"
	"io/fs"
	"net"

	"github.com/livebud/bud/internal/prompter"

	"github.com/livebud/bud/internal/sh"

	"github.com/livebud/bud"
	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/framework/web/webrt"
	"github.com/livebud/bud/internal/dag"
	"github.com/livebud/bud/package/budhttp/budsvr"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/js"
	"github.com/livebud/bud/package/socket"

	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/log/console"
	"github.com/livebud/bud/package/log/levelfilter"
)

func (c *CLI) module(cmd *sh.Command) (*gomod.Module, error) {
	module, err := gomod.Find(c.Dir)
	if err != nil {
		return nil, err
	}
	cmd.Env = append(cmd.Env, "GOMODCACHE="+module.ModCache())
	return module, nil
}

func (c *CLI) logger() (log.Log, error) {
	level, err := log.ParseLevel(c.Log)
	if err != nil {
		return nil, err
	}
	logger := log.New(levelfilter.New(console.New(c.Stderr), level))
	return logger, nil
}

func (c *CLI) listenWeb(addr string) (socket.Listener, error) {
	if c.WebListener != nil {
		return c.WebListener, nil
	}
	return socket.ListenUp(addr, 5)
}

func (c *CLI) prompter(webLn net.Listener) *prompter.Prompter {
	var prompter prompter.Prompter
	c.Stdout = io.MultiWriter(c.Stdout, &prompter.StdOut)
	c.Stderr = io.MultiWriter(c.Stderr, &prompter.StdErr)
	prompter.Init(webLn.Addr().String())
	return &prompter

}

func (c *CLI) listenDev(addr string) (socket.Listener, error) {
	if c.DevListener != nil {
		return c.DevListener, nil
	}
	ln, err := socket.ListenUp(addr, 5)
	if err != nil {
		return nil, err
	}
	return ln, err
}

func (c *CLI) listenFile(addr string) (socket.Listener, error) {
	if c.FileListener != nil {
		return c.FileListener, nil
	}
	return socket.ListenUp(addr, 5)
}

func (c *CLI) openDatabase(log log.Log, module *gomod.Module) (*dag.DB, error) {
	return dag.Load(log, module.Directory("bud", "bud.db"))
}

func (c *CLI) serveDev(ctx context.Context, cfg *bud.Config, remotefs fs.FS, ln socket.Listener, log log.Log, vm js.VM) error {
	// TODO: replace with *bud.Config
	flag := &framework.Flag{
		Embed:  cfg.Embed,
		Minify: cfg.Minify,
		Hot:    cfg.Hot,
	}
	handler := budsvr.NewHandler(flag, remotefs, c.Bus, log, vm)
	// TODO: replace with something else
	return webrt.Serve(ctx, ln, handler)
}
