package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net"
	"path/filepath"
	"strings"

	"github.com/livebud/bud/internal/dsync"
	"github.com/livebud/bud/internal/once"
	"github.com/livebud/bud/internal/prompter"
	"github.com/livebud/bud/internal/sh"

	"github.com/livebud/bud"
	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/framework/afs"
	"github.com/livebud/bud/framework/generator"
	"github.com/livebud/bud/internal/dag"
	"github.com/livebud/bud/package/budhttp/budsvr"
	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/genfs"
	"github.com/livebud/bud/package/gomod"
	v8 "github.com/livebud/bud/package/js/v8"
	"github.com/livebud/bud/package/parser"
	"github.com/livebud/bud/package/remotefs"
	"github.com/livebud/bud/package/socket"
	"github.com/livebud/bud/package/virtual"

	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/log/console"
	"github.com/livebud/bud/package/log/levelfilter"
)

func (c *CLI) module() (*gomod.Module, error) {
	module, err := gomod.Find(c.Dir)
	if err != nil {
		return nil, err
	}
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

func (c *CLI) startDev(ctx context.Context, cfg *bud.Config, ln socket.Listener, log log.Log) (io.Closer, error) {
	var closer once.Closer

	vm, err := v8.Load()
	if err != nil {
		return nil, err
	}
	closer.Add(func() error {
		vm.Close()
		return nil
	})

	// TODO: replace with *bud.Config
	flag := &framework.Flag{
		Embed:  cfg.Embed,
		Minify: cfg.Minify,
		Hot:    cfg.Hot,
	}
	server := budsvr.New(ln, c.Bus, flag, virtual.Map{}, log, vm)
	server.Start(ctx)
	closer.Add(server.Close)
	return &closer, nil
}

func (c *CLI) genfs(cfg *bud.Config, db *dag.DB, log log.Log, module *gomod.Module) *genfs.FileSystem {
	fsys := virtual.Exclude(module, func(path string) bool {
		return path == "bud" || strings.HasPrefix(path, "bud/")
	})
	genfs := genfs.New(db, fsys, log)
	parser := parser.New(genfs, module)
	injector := di.New(genfs, log, module, parser)
	genfs.FileGenerator("bud/internal/generator/generator.go", generator.New(log, module, parser))
	genfs.FileGenerator("bud/cmd/afs/main.go", afs.New(cfg, injector, log, module))
	return genfs
}

func expectEnv(env []string, key string) error {
	for _, keyValue := range env {
		if strings.HasPrefix(keyValue, key+"=") {
			return nil
		}
	}
	return fmt.Errorf("%q is missing from the environment", key)
}

var afsSkips = []func(name string, isDir bool) bool{
	// Skip hidden files and directories
	func(name string, isDir bool) bool {
		base := filepath.Base(name)
		return base[0] == '_' || base[0] == '.'
	},
	// Skip files we want to carry over
	func(name string, isDir bool) bool {
		switch name {
		case "bud/bud.db", "bud/afs", "bud/app":
			return true
		default:
			return false
		}
	},
}

var appSkips = append(afsSkips,
	// Skip over the afs files we just generated
	func(name string, isDir bool) bool {
		return isAFSPath(name)
	},
)

// Generate AFS
func (c *CLI) generateAFS(ctx context.Context, cmd *sh.Command, fsys *genfs.FileSystem, log log.Log, module *gomod.Module) error {
	if err := dsync.To(fsys, module, "bud", dsync.WithSkip(afsSkips...), dsync.WithLog(log)); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}
	// Build the afs binary
	if err := cmd.Run(ctx, "go", "build", "-mod=mod", "-o="+afsBinPath, afsMainPath); err != nil {
		return err
	}
	return nil
}

func (c *CLI) startAFS(ctx context.Context, cmd *sh.Command, module *gomod.Module) (*sh.Process, error) {
	// Ensure the command has the BUD_DEV_URL environment variable
	if err := expectEnv(cmd.Env, "BUD_DEV_URL"); err != nil {
		return nil, err
	}
	// Ensure the command has the AFS_FDS_START environment variable
	if err := expectEnv(cmd.Env, "AFS_FDS_START"); err != nil {
		return nil, err
	}
	// Start afs
	return cmd.Start(ctx, module.Directory("bud", "afs"))
}

func (c *CLI) buildApp(ctx context.Context, remoteClient *remotefs.Client, cmd *sh.Command, log log.Log, module *gomod.Module) error {
	// Sync the app files again with the remote filesystem
	if err := dsync.To(remoteClient, module, "bud", dsync.WithSkip(appSkips...), dsync.WithLog(log)); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}
	// Build the application binary
	if err := cmd.Run(ctx, "go", "build", "-mod=mod", "-o="+appBinPath, appMainPath); err != nil {
		return err
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
