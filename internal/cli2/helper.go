package cli

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"net"
	"path/filepath"
	"strings"

	"github.com/livebud/bud/internal/dsync"
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

func (c *CLI) listenAFS(addr string) (socket.Listener, error) {
	if c.AFSListener != nil {
		return c.AFSListener, nil
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

func (c *CLI) generateAFS(genfs fs.FS, log log.Log, module *gomod.Module) error {

	skips := []func(name string, isDir bool) bool{}
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

	return nil
}

func (c *CLI) generateApp(remotefs fs.FS, log log.Log, module *gomod.Module) error {
	skips := []func(name string, isDir bool) bool{}
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
	// Skip over the afs files we just generated
	skips = append(skips, func(name string, isDir bool) bool {
		return isAFSPath(name)
	})

	// Sync the app files again with the remote filesystem
	if err := dsync.To(remotefs, module, "bud", dsync.WithSkip(skips...), dsync.WithLog(log)); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}

	return nil
}

const (
	// internal/generator
	afsGeneratorDir = "bud/internal/generator"

	// cmd/afs
	afsMainPath = "bud/cmd/afs/main.go"
	afsMainDir  = "bud/cmd/afs"
	afsBinPath  = "bud/afs"

	// cmd/app
	appMainPath = "bud/internal/app/main.go"
	appBinPath  = "bud/app"
)

func isAFSPath(fpath string) bool {
	return fpath == afsBinPath ||
		isWithin(afsGeneratorDir, fpath) ||
		isWithin(afsMainDir, fpath) ||
		fpath == "bud/cmd" // TODO: remove once we move app over to cmd/app/main.go
}

func isWithin(parent, child string) bool {
	if parent == child {
		return true
	}
	return strings.HasPrefix(child, parent)
}
