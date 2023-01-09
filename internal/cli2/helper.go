package cli

import (
	"context"
	"io"
	"net"
	"strings"

	"github.com/livebud/bud/internal/prompter"

	"github.com/livebud/bud"
	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/framework/afs"
	"github.com/livebud/bud/framework/generator"
	"github.com/livebud/bud/framework/web/webrt"
	"github.com/livebud/bud/internal/dag"
	"github.com/livebud/bud/package/budhttp/budsvr"
	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/genfs"
	"github.com/livebud/bud/package/gomod"
	v8 "github.com/livebud/bud/package/js/v8"
	"github.com/livebud/bud/package/parser"
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

func (c *CLI) serveDev(ctx context.Context, cfg *bud.Config, ln socket.Listener, log log.Log) error {
	vm, err := v8.Load()
	if err != nil {
		return err
	}
	defer vm.Close()

	// TODO: replace with *bud.Config
	flag := &framework.Flag{
		Embed:  cfg.Embed,
		Minify: cfg.Minify,
		Hot:    cfg.Hot,
	}
	handler := budsvr.NewHandler(flag, virtual.Map{}, c.Bus, log, vm)
	// TODO: replace with something else
	return webrt.Serve(ctx, ln, handler)
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
