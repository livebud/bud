package budapi

import (
	"context"
	"io"
	"net"
	"os"

	"github.com/livebud/bud/internal/dag"
	"github.com/livebud/bud/internal/shell"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/log/console"
	"github.com/livebud/bud/package/log/levelfilter"
	"github.com/livebud/bud/package/socket"
)

type Bud struct {
	Dir string
	Log string

	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
	Env    []string

	WebListener net.Listener
	DevListener net.Listener
	AFSListener net.Listener
}

type Interface interface {
	Create(ctx context.Context, config *Create) error
	Generate(ctx context.Context, config *Generate) error
	Build(ctx context.Context, config *Build) error
	Run(ctx context.Context, config *RunInput) error
	Watch(ctx context.Context, config *Watch) error
}

type Flag struct {
	Embed  bool
	Hot    bool
	Minify bool
}

type Create struct {
}

type Generate struct {
	Flag
	Paths []string
}

type Build struct {
	Flag
}

type RunInput struct {
	Flag
	ListenWeb string
	ListenDev string
}

type Watch struct {
	Flag
}

// type Bud struct {
// 	Dir string

// 	Stdin  io.Reader
// 	Stdout io.Writer
// 	Stderr io.Writer
// 	Env    []string

// 	WebLn net.Listener
// 	DevLn net.Listener
// 	AfsLn net.Listener
// }

// Load sets default values for Bud.
func (b *Bud) load() {
	if b.Dir == "" {
		b.Dir = "."
	}
	if b.Stdin == nil {
		b.Stdin = os.Stdin
	}
	if b.Stdout == nil {
		b.Stdout = os.Stdout
	}
	if b.Stderr == nil {
		b.Stderr = os.Stderr
	}
	if b.Env == nil {
		b.Env = os.Environ()
	}
}

func (p *Bud) findModule() (*gomod.Module, error) {
	return gomod.Find(p.Dir)
}

func (p *Bud) loadLog(w io.Writer, level string) (log.Log, error) {
	lvl, err := log.ParseLevel(level)
	if err != nil {
		return nil, err
	}
	log := log.New(levelfilter.New(console.New(w), lvl))
	return log, nil
}

func (p *Bud) listenWeb(listenWeb string) (net.Listener, error) {
	if p.WebListener != nil {
		return p.WebListener, nil
	}
	return socket.Listen(listenWeb)
}

func (p *Bud) listenDev(listenDev string) (net.Listener, error) {
	if p.DevListener != nil {
		return p.DevListener, nil
	}
	return socket.Listen(listenDev)
}

func (p *Bud) openDb(log log.Log, module *gomod.Module) (*dag.DB, error) {
	return dag.Load(log, module.Directory("bud", "bud.db"))
}

func (p *Bud) command(module *gomod.Module) *shell.Command {
	return &shell.Command{
		Dir:    module.Directory(),
		Stdin:  p.Stdin,
		Stdout: p.Stdout,
		Stderr: p.Stderr,
		Env:    p.Env,
	}
}

func (p *Bud) serveDev(ctx context.Context, devLn net.Listener) error {
	// handler := hot.New(log, ps)
	// hot.New
	// budsvr.New()
	return nil
}

func (p *Bud) generateAFS(ctx context.Context, pkgs ...string) error {
	return nil
}

func (p *Bud) serveAFS(ctx context.Context) error {
	return nil
}

func (p *Bud) generateApp(ctx context.Context, pkgs ...string) error {
	return nil
}

func (p *Bud) startApp(ctx context.Context) error {
	return nil
}
