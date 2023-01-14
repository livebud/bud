package cli

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/livebud/bud/package/remotefs"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/framework/afs"
	"github.com/livebud/bud/framework/generator"
	"github.com/livebud/bud/internal/dag"
	"github.com/livebud/bud/internal/once"
	"github.com/livebud/bud/internal/prompter"
	"github.com/livebud/bud/internal/pubsub"
	"github.com/livebud/bud/internal/sh"
	"github.com/livebud/bud/package/budhttp/budsvr"
	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/genfs"
	"github.com/livebud/bud/package/gomod"
	v8 "github.com/livebud/bud/package/js/v8"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/log/console"
	"github.com/livebud/bud/package/log/levelfilter"
	"github.com/livebud/bud/package/parser"
	"github.com/livebud/bud/package/socket"
	"github.com/livebud/bud/package/virtual"
)

func newProvider(ctx context.Context, cli *CLI) *provider {
	return &provider{
		cli,
		nil,
		&once.Closer{},
		nil,
		ctx,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
	}
}

type provider struct {
	*CLI

	// Private cached provided values
	afsFile   *os.File
	closer    *once.Closer
	command   *sh.Command
	ctx       context.Context
	db        *dag.DB
	devServer *budsvr.Server
	flag      *framework.Flag
	genfs     *genfs.FileSystem
	logger    *log.Logger
	module    *gomod.Module
	prompter  *prompter.Prompter
	webFile   *os.File
	v8        *v8.VM
}

func (p *provider) Close() error {
	return p.closer.Close()
}

func (p *provider) Module() (*gomod.Module, error) {
	if p.module != nil {
		return p.module, nil
	}
	module, err := gomod.Find(p.Dir)
	if err != nil {
		return nil, err
	}
	p.module = module
	return p.module, nil
}

func (p *provider) Logger() (*log.Logger, error) {
	if p.logger != nil {
		return p.logger, nil
	}
	level, err := log.ParseLevel(p.Log)
	if err != nil {
		return nil, err
	}
	p.logger = log.New(levelfilter.New(console.New(p.Stderr), level))
	return p.logger, nil
}

func (p *provider) Database() (*dag.DB, error) {
	if p.db != nil {
		return p.db, nil
	}
	module, err := p.Module()
	if err != nil {
		return nil, err
	}
	log, err := p.Logger()
	if err != nil {
		return nil, err
	}
	db, err := dag.Load(log, module.Directory("bud", "bud.db"))
	if err != nil {
		return nil, err
	}
	p.closer.Add(db.Close)
	p.db = db
	return p.db, nil
}

func (p *provider) Command() *sh.Command {
	if p.command != nil {
		return p.command
	}
	module, err := p.Module()
	if err != nil {
		return nil
	}
	p.command = &sh.Command{
		Dir:    module.Directory(),
		Env:    p.Env,
		Stdin:  p.Stdin,
		Stdout: p.Stdout,
		Stderr: p.Stderr,
	}
	// Add the environment variable
	p.command.Env = append(p.command.Env, "GOMODCACHE="+module.ModCache())
	return p.command
}

func (p *provider) GenFS() (*genfs.FileSystem, error) {
	if p.genfs != nil {
		return p.genfs, nil
	}
	module, err := p.Module()
	if err != nil {
		return nil, err
	}
	log, err := p.Logger()
	if err != nil {
		return nil, err
	}
	db, err := p.Database()
	if err != nil {
		return nil, err
	}
	fsys := virtual.Exclude(module, func(path string) bool {
		return path == "bud" || strings.HasPrefix(path, "bud/")
	})
	genfs := genfs.New(db, fsys, log)
	parser := parser.New(genfs, module)
	injector := di.New(genfs, log, module, parser)
	genfs.FileGenerator("bud/internal/generator/generator.go", generator.New(log, module, parser))
	genfs.FileGenerator("bud/cmd/afs/main.go", afs.New(p.Config, injector, log, module))
	return genfs, nil
}

func (p *provider) DevListener() (socket.Listener, error) {
	if p.CLI.DevListener != nil {
		return p.CLI.DevListener, nil
	}
	devLn, err := socket.Listen(p.CLI.ListenDev)
	if err != nil {
		return nil, err
	}
	p.closer.Add(devLn.Close)
	p.CLI.DevListener = devLn
	return p.CLI.DevListener, nil
}

func (p *provider) AFSListener() (socket.Listener, error) {
	if p.CLI.AFSListener != nil {
		return p.CLI.AFSListener, nil
	}
	afsLn, err := socket.Listen(p.CLI.ListenAFS)
	if err != nil {
		return nil, err
	}
	p.closer.Add(afsLn.Close)
	p.CLI.AFSListener = afsLn
	return p.CLI.AFSListener, nil
}

func (p *provider) AFSFile() (*os.File, error) {
	if p.afsFile != nil {
		return p.afsFile, nil
	}
	afsLn, err := p.AFSListener()
	if err != nil {
		return nil, err
	}
	afsFile, err := afsLn.File()
	if err != nil {
		return nil, err
	}
	p.afsFile = afsFile
	p.closer.Add(afsFile.Close)
	return afsFile, nil
}

func (p *provider) Bus() pubsub.Client {
	if p.CLI.Bus != nil {
		return p.CLI.Bus
	}
	p.CLI.Bus = pubsub.New()
	return p.CLI.Bus
}

// TODO: replace with *bud.Flag
func (p *provider) Flag() *framework.Flag {
	if p.flag != nil {
		return p.flag
	}
	p.flag = &framework.Flag{
		Embed:  p.Embed,
		Minify: p.Minify,
		Hot:    p.Hot,
	}
	return p.flag
}

func (p *provider) V8() (*v8.VM, error) {
	if p.v8 != nil {
		return p.v8, nil
	}
	vm, err := v8.Load()
	if err != nil {
		return nil, err
	}
	p.v8 = vm
	p.closer.Add(func() error {
		vm.Close()
		return nil
	})
	return vm, nil
}

func (p *provider) DevServer() (*budsvr.Server, error) {
	if p.devServer != nil {
		return p.devServer, nil
	}
	devLn, err := p.DevListener()
	if err != nil {
		return nil, err
	}
	bus := p.Bus()
	// TODO: replace with *bud.Config
	flag := p.Flag()
	log, err := p.Logger()
	if err != nil {
		return nil, err
	}
	v8, err := p.V8()
	if err != nil {
		return nil, err
	}
	p.devServer = budsvr.New(devLn, bus, flag, virtual.Map{}, log, v8)
	return p.devServer, nil
}

func (p *provider) WebListener() (socket.Listener, error) {
	if p.CLI.WebListener != nil {
		return p.CLI.WebListener, nil
	}
	webLn, err := socket.ListenUp(p.ListenWeb, 5)
	if err != nil {
		return nil, err
	}
	p.closer.Add(webLn.Close)
	p.CLI.WebListener = webLn
	return p.CLI.WebListener, nil
}

func (p *provider) WebFile() (*os.File, error) {
	if p.webFile != nil {
		return p.webFile, nil
	}
	webLn, err := p.WebListener()
	if err != nil {
		return nil, err
	}
	webFile, err := webLn.File()
	if err != nil {
		return nil, err
	}
	p.webFile = webFile
	p.closer.Add(webFile.Close)
	return webFile, nil
}

func (p *provider) AFSClient() (*remotefs.Client, error) {
	afsLn, err := p.AFSListener()
	if err != nil {
		return nil, err
	}
	remoteClient, err := remotefs.Dial(p.ctx, afsLn.Addr().String())
	if err != nil {
		return nil, err
	}
	p.closer.Add(remoteClient.Close)
	return remoteClient, nil
}

func (p *provider) AFSCommand() (*sh.Command, error) {
	afsFile, err := p.AFSFile()
	if err != nil {
		return nil, err
	}

	devLn, err := p.DevListener()
	if err != nil {
		return nil, err
	}

	cmd := &sh.Command{
		Dir:    p.Dir,
		Env:    p.Env,
		Stdin:  p.Stdin,
		Stderr: p.Stderr,
		Stdout: p.Stdout,
	}
	// Start the application file server
	cmd.Env = append(cmd.Env, "BUD_DEV_URL="+devLn.Addr().String())
	// Inject the file under the AFS prefix
	cmd.Inject("AFS", afsFile)

	return cmd, nil
}

func (p *provider) AppCommand() (*sh.Command, error) {
	webFile, err := p.WebFile()
	if err != nil {
		return nil, err
	}

	afsLn, err := p.AFSListener()
	if err != nil {
		return nil, err
	}

	devLn, err := p.DevListener()
	if err != nil {
		return nil, err
	}

	cmd := &sh.Command{
		Dir:    p.Dir,
		Env:    p.Env,
		Stdin:  p.Stdin,
		Stderr: p.Stderr,
		Stdout: p.Stdout,
	}

	// Inject that file under the WEB prefix
	cmd.Inject("WEB", webFile)

	// Setup the environment
	cmd.Env = append(cmd.Env,
		"BUD_AFS_URL="+afsLn.Addr().String(),
		"BUD_DEV_URL="+devLn.Addr().String(),
	)

	return cmd, nil
}

func (p *provider) Prompter() (*prompter.Prompter, error) {
	if p.prompter != nil {
		return p.prompter, nil
	}
	webln, err := p.WebListener()
	if err != nil {
		return nil, err
	}
	var prompter prompter.Prompter
	p.Stdout = io.MultiWriter(p.Stdout, &prompter.StdOut)
	p.Stderr = io.MultiWriter(p.Stderr, &prompter.StdErr)
	p.prompter = &prompter
	p.prompter.Init(webln.Addr().String())
	return p.prompter, nil
}
