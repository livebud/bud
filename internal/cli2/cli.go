package cli

import (
	"context"
	"io"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/livebud/bud/internal/extrafile"
	"github.com/livebud/bud/internal/prompter"

	"github.com/livebud/bud/framework/afs"
	"github.com/livebud/bud/framework/generator"
	"github.com/livebud/bud/package/budhttp/budsvr"
	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/genfs"
	v8 "github.com/livebud/bud/package/js/v8"
	"github.com/livebud/bud/package/parser"
	"github.com/livebud/bud/package/remotefs"
	"github.com/livebud/bud/package/virtual"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/dag"
	"github.com/livebud/bud/internal/pubsub"
	"github.com/livebud/bud/internal/versions"
	"github.com/livebud/bud/package/commander"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/log/console"
	"github.com/livebud/bud/package/log/levelfilter"
	"github.com/livebud/bud/package/socket"
)

func New() *CLI {
	return &CLI{
		".",
		"info",
		os.Stdin,
		os.Stdout,
		os.Stderr,
		os.Environ(),
		nil,
		nil,
		nil,
		nil,
	}
}

type CLI struct {
	Dir string
	Log string

	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
	Env    []string

	// Used for testing
	Bus         pubsub.Client
	WebListener socket.Listener
	DevListener socket.Listener
	AFSListener socket.Listener
}

func (c *CLI) Parse(ctx context.Context, args ...string) error {
	// Check that we have a valid Go version
	if err := versions.CheckGo(runtime.Version()); err != nil {
		return err
	}

	// $ bud [args...]
	cmd := &Custom{}
	cli := commander.New("bud").Writer(c.Stdout)
	cli.Flag("chdir", "change the working directory").Short('C').String(&c.Dir).Default(c.Dir)
	cli.Flag("help", "show this help message").Short('h').Bool(&cmd.Help).Default(false)
	cli.Flag("log", "filter logs with this pattern").Short('L').String(&c.Log).Default("info")
	cli.Args("args").Strings(&cmd.Args)
	cli.Run(func(ctx context.Context) error { return c.Custom(ctx, cmd) })

	{ // $ bud generate [packages...]
		cmd := &Generate{Flag: &framework.Flag{}}
		cli := cli.Command("generate", "generate bud packages")
		cli.Flag("listen-dev", "dev server address").String(&cmd.ListenDev).Default(":35729")
		cli.Flag("listen-afs", "app file server address").String(&cmd.ListenAFS).Default(":0")
		cli.Flag("embed", "embed assets").Bool(&cmd.Flag.Embed).Default(false)
		cli.Flag("hot", "hot reloading").Bool(&cmd.Flag.Hot).Default(true)
		cli.Flag("minify", "minify assets").Bool(&cmd.Flag.Minify).Default(false)
		cli.Args("packages").Strings(&cmd.Packages)
		cli.Run(func(ctx context.Context) error { return c.Generate(ctx, cmd) })
	}

	{ // $ bud run
		cmd := &Run{Flag: &framework.Flag{}}
		cli := cli.Command("run", "run the app in development mode")
		cli.Flag("listen", "address to listen on").String(&cmd.ListenWeb).Default(":3000")
		cli.Flag("listen-dev", "dev server address").String(&cmd.ListenDev).Default(":35729")
		cli.Flag("listen-afs", "app file server address").String(&cmd.ListenAFS).Default(":0")
		cli.Flag("embed", "embed assets").Bool(&cmd.Flag.Embed).Default(false)
		cli.Flag("hot", "hot reloading").Bool(&cmd.Flag.Hot).Default(true)
		cli.Flag("minify", "minify assets").Bool(&cmd.Flag.Minify).Default(false)
		cli.Flag("watch", "watch for file changes").Bool(&cmd.Watch).Default(true)
		cli.Run(func(ctx context.Context) error { return c.Run(ctx, cmd) })
	}

	return cli.Parse(ctx, args)
}

// findModule finds the module for the current working directory.
func (c *CLI) findModule() (*gomod.Module, error) {
	return gomod.Find(c.Dir)
}

// loadLog loads the logger
func (c *CLI) loadLog() (log.Log, error) {
	level, err := log.ParseLevel(c.Log)
	if err != nil {
		return nil, err
	}
	log := log.New(levelfilter.New(console.New(c.Stderr), level))
	return log, nil
}

func (c *CLI) openDB(log log.Log, module *gomod.Module) (*dag.DB, error) {
	db, err := dag.Load(log, module.Directory("bud", "bud.db"))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (c *CLI) genfs(cache genfs.Cache, flag *framework.Flag, log log.Log, module *gomod.Module) *genfs.FileSystem {
	fsys := virtual.Exclude(module, func(path string) bool {
		return path == "bud" || strings.HasPrefix(path, "bud/")
	})
	genfs := genfs.New(cache, fsys, log)
	parser := parser.New(genfs, module)
	injector := di.New(genfs, log, module, parser)
	genfs.FileGenerator("bud/internal/generator/generator.go", generator.New(log, module, parser))
	genfs.FileGenerator("bud/cmd/afs/main.go", afs.New(flag, injector, log, module))
	return genfs
}

func (c *CLI) listenDev(listenDev string) (socket.Listener, error) {
	if c.DevListener != nil {
		return c.DevListener, nil
	}
	return socket.Listen(listenDev)
}

func (c *CLI) devServer(bus pubsub.Client, devLn net.Listener, flag *framework.Flag, log log.Log, v8 *v8.VM) *budsvr.Server {
	return budsvr.New(devLn, bus, flag, virtual.Map{}, log, v8)
}

func (c *CLI) listenAFS(listenAFS string) (socket.Listener, error) {
	if c.AFSListener != nil {
		return c.AFSListener, nil
	}
	return socket.Listen(listenAFS)
}

func (c *CLI) listenFileAFS(afsLn socket.Listener) (*os.File, error) {
	fileAFS, err := afsLn.File()
	if err != nil {
		return nil, err
	}
	return fileAFS, nil
}

func (c *CLI) newBus() pubsub.Client {
	if c.Bus != nil {
		return c.Bus
	}
	c.Bus = pubsub.New()
	return c.Bus
}

func (c *CLI) loadV8() (*v8.VM, error) {
	return v8.Load()
}

func (c *CLI) listenWeb(listenWeb string) (socket.Listener, error) {
	if c.WebListener != nil {
		return c.WebListener, nil
	}
	return socket.ListenUp(listenWeb, 5)
}

func (c *CLI) listenFileWeb(webLn socket.Listener) (*os.File, error) {
	return webLn.File()
}

func (c *CLI) dialAFS(ctx context.Context, afsLn net.Listener) (*remotefs.Client, error) {
	return remotefs.Dial(ctx, afsLn.Addr().String())
}

func (c *CLI) loadCommandAFS(module *gomod.Module, afsFile *os.File, devLn net.Listener) (*exec.Cmd, error) {
	// Initialize the command
	cmd := c.newCommand(module, module.Directory("bud", "afs"))
	// Setup the environment
	cmd.Env = append(cmd.Env,
		"BUD_DEV_URL="+devLn.Addr().String(),
	)
	// Inject the file under the AFS prefix
	extrafile.Inject(&cmd.ExtraFiles, &cmd.Env, "AFS", afsFile)
	return cmd, nil
}

func (c *CLI) loadCommandApp(module *gomod.Module, afsLn, devLn net.Listener, webFile *os.File) (*exec.Cmd, error) {
	// Initialize the command
	cmd := c.newCommand(module, module.Directory("bud", "app"))
	// Setup the environment
	cmd.Env = append(cmd.Env,
		"BUD_AFS_URL="+afsLn.Addr().String(),
		"BUD_DEV_URL="+devLn.Addr().String(),
	)
	// Inject that file under the WEB prefix
	extrafile.Inject(&cmd.ExtraFiles, &cmd.Env, "WEB", webFile)
	return cmd, nil
}

func (c *CLI) newPrompter(webLn net.Listener) *prompter.Prompter {
	var prompter prompter.Prompter
	c.Stdout = io.MultiWriter(c.Stdout, &prompter.StdOut)
	c.Stderr = io.MultiWriter(c.Stderr, &prompter.StdErr)
	prompter.Init(webLn.Addr().String())
	return &prompter
}

func (c *CLI) newCommand(module *gomod.Module, name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.Dir = module.Directory()
	cmd.Env = append([]string{}, c.Env...)
	cmd.Stdin = c.Stdin
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr
	return cmd
}
