package cli

import (
	"context"
	"embed"
	"errors"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/livebud/bud/internal/shell"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/framework/afs"
	"github.com/livebud/bud/framework/generator"
	"github.com/livebud/bud/framework/transpiler"
	"github.com/livebud/bud/internal/current"
	"github.com/livebud/bud/internal/dag"
	"github.com/livebud/bud/internal/extrafile"
	"github.com/livebud/bud/internal/prompter"
	"github.com/livebud/bud/internal/pubsub"
	"github.com/livebud/bud/internal/versions"
	"github.com/livebud/bud/package/budhttp/budsvr"
	"github.com/livebud/bud/package/commander"
	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/genfs"
	"github.com/livebud/bud/package/gomod"
	v8 "github.com/livebud/bud/package/js/v8"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/log/console"
	"github.com/livebud/bud/package/log/levelfilter"
	"github.com/livebud/bud/package/parser"
	"github.com/livebud/bud/package/remotefs"
	"github.com/livebud/bud/package/socket"
	"github.com/livebud/bud/package/virtual"
)

// Shared list of embedded files
//
//go:embed new_controller/*.gotext create/*.gotext
var embedFS embed.FS

// TODO: figuring out how the resource lifecycle should be managed. Ideally
// the closer is optional, and if it is not passed in, the CLI will manage the
// resources.
func New(closer closer) *CLI {
	return &CLI{
		".",
		"info",
		os.Stdin,
		os.Stdout,
		os.Stderr,
		os.Environ(),
		closer,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
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

type closer interface {
	Add(func() error)
	Close() error
}

type CLI struct {
	Dir string
	Log string

	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
	Env    []string

	// Closer is used to manage resources
	Closer closer

	// Used for testing
	Bus         pubsub.Client
	WebListener socket.Listener
	DevListener socket.Listener
	AFSListener socket.Listener

	// Internal state
	db         *dag.DB
	module     *gomod.Module
	budModule  *gomod.Module
	log        log.Log
	afsFile    *os.File
	webFile    *os.File
	v8         *v8.VM
	ds         *budsvr.Server
	genfs      genfs.FileSystem
	afsClient  *remotefs.Client
	afsProcess *shell.Process
	appProcess *shell.Process
}

func (c *CLI) Parse(ctx context.Context, args ...string) error {
	// Check that we have a valid Go version
	if err := versions.CheckGo(runtime.Version()); err != nil {
		return err
	}

	// $ bud [args...]
	in := &Custom{}
	cli := commander.New("bud").Writer(c.Stdout)
	cli.Flag("chdir", "change the working directory").Short('C').String(&c.Dir).Default(c.Dir)
	cli.Flag("help", "show this help message").Short('h').Bool(&in.Help).Default(false)
	cli.Flag("log", "filter logs with this pattern").Short('L').String(&c.Log).Default("info")
	cli.Args("args").Strings(&in.Args)
	cli.Run(func(ctx context.Context) error { return c.Custom(ctx, in) })

	{ // $ bud create <dir>
		in := &Create{}
		cli := cli.Command("create", "create a new app")
		cli.Flag("dev", "link to the development version").Short('D').Bool(&in.Dev).Default(versions.Bud == "latest")
		cli.Flag("module", "module path for go.mod").String(&in.Module).Optional()
		cli.Arg("dir").String(&c.Dir)
		cli.Run(func(ctx context.Context) error { return c.Create(ctx, in) })
	}

	{ // $ bud generate [packages...]
		in := &Generate{Flag: &framework.Flag{}}
		cli := cli.Command("generate", "generate bud packages")
		cli.Flag("listen-dev", "dev server address").String(&in.ListenDev).Default(":35729")
		cli.Flag("listen-afs", "app file server address").String(&in.ListenAFS).Default(":0")
		cli.Flag("embed", "embed assets").Bool(&in.Flag.Embed).Default(false)
		cli.Flag("hot", "hot reloading").Bool(&in.Flag.Hot).Default(true)
		cli.Flag("minify", "minify assets").Bool(&in.Flag.Minify).Default(false)
		cli.Args("packages").Strings(&in.Packages)
		cli.Run(func(ctx context.Context) error { return c.Generate(ctx, in) })
	}

	{ // $ bud run
		in := &Run{Flag: &framework.Flag{}}
		cli := cli.Command("run", "run the app in development mode")
		cli.Flag("listen", "address to listen on").String(&in.ListenWeb).Default(":3000")
		cli.Flag("listen-dev", "dev server address").String(&in.ListenDev).Default(":35729")
		cli.Flag("listen-afs", "app file server address").String(&in.ListenAFS).Default(":0")
		cli.Flag("embed", "embed assets").Bool(&in.Flag.Embed).Default(false)
		cli.Flag("hot", "hot reloading").Bool(&in.Flag.Hot).Default(true)
		cli.Flag("minify", "minify assets").Bool(&in.Flag.Minify).Default(false)
		cli.Flag("watch", "watch for file changes").Bool(&in.Watch).Default(true)
		cli.Run(func(ctx context.Context) error { return c.Run(ctx, in) })
	}

	{ // $ bud build
		in := &Build{Flag: &framework.Flag{}}
		cli := cli.Command("build", "build your app into a single binary")
		cli.Flag("embed", "embed assets").Bool(&in.Flag.Embed).Default(true)
		cli.Flag("minify", "minify assets").Bool(&in.Flag.Minify).Default(true)
		cli.Run(func(ctx context.Context) error { return c.Build(ctx, in) })
	}

	{ // $ bud new
		cli := cli.Command("new", "scaffold code for your app")

		{ // $ bud new controller <name> [actions...]
			in := &NewController{}
			cli := cli.Command("controller", "scaffold a new controller")
			cli.Arg("path").String(&in.Path)
			cli.Args("actions").Strings(&in.Actions)
			cli.Run(func(ctx context.Context) error { return c.NewController(ctx, in) })
		}

	}

	{ // $ bud tool
		cli := cli.Command("tool", "extra tools")

		{ // $ bud tool ds
			in := &ToolDS{Flag: &framework.Flag{}}
			cli := cli.Command("bs", "run the bud server")
			cli.Flag("embed", "embed assets").Bool(&in.Flag.Embed).Default(false)
			cli.Flag("hot", "hot reloading").Bool(&in.Flag.Hot).Default(true)
			cli.Flag("minify", "minify assets").Bool(&in.Flag.Minify).Default(false)
			cli.Flag("listen-dev", "dev server address").String(&in.ListenDev).Default(":35729")
			cli.Run(func(ctx context.Context) error { return c.ToolDS(ctx, in) })
		}

		{ // $ bud tool di
			in := &ToolDi{}
			cli := cli.Command("di", "dependency injection generator")
			cli.Flag("name", "name of the function").String(&in.Name).Default("Load")
			cli.Flag("dependency", "generate dependency provider").Short('d').Strings(&in.Dependencies)
			cli.Flag("external", "mark dependency as external").Short('e').Strings(&in.Externals).Optional()
			cli.Flag("map", "map interface types to concrete types").Short('m').StringMap(&in.Map).Optional()
			cli.Flag("target", "target import path").Short('t').String(&in.Target)
			cli.Flag("hoist", "hoist dependencies that depend on externals").Bool(&in.Hoist).Default(false)
			cli.Flag("verbose", "verbose logging").Short('v').Bool(&in.Verbose).Default(false)
			cli.Run(func(ctx context.Context) error { return c.ToolDi(ctx, in) })
		}

		{ // $ bud tool fs
			cli := cli.Command("fs", "filesystem tools")

			{ // $ bud tool fs ls [dir]
				in := &ToolFsLs{Flag: &framework.Flag{}}
				cli := cli.Command("ls", "list a directory")
				cli.Flag("embed", "embed assets").Bool(&in.Flag.Embed).Default(false)
				cli.Flag("hot", "hot reloading").Bool(&in.Flag.Hot).Default(false)
				cli.Flag("minify", "minify assets").Bool(&in.Flag.Minify).Default(false)
				cli.Arg("dir").String(&in.Path).Default(".")
				cli.Run(func(ctx context.Context) error { return c.ToolFsLs(ctx, in) })
			}

			{ // $ bud tool fs cat [path]
				// TODO: better align with the unix `cat` command
				in := &ToolFsCat{Flag: &framework.Flag{}}
				cli := cli.Command("cat", "print a file")
				cli.Flag("embed", "embed assets").Bool(&in.Flag.Embed).Default(false)
				cli.Flag("hot", "hot reloading").Bool(&in.Flag.Hot).Default(false)
				cli.Flag("minify", "minify assets").Bool(&in.Flag.Minify).Default(false)
				cli.Arg("path").String(&in.Path)
				cli.Run(func(ctx context.Context) error { return c.ToolFsCat(ctx, in) })
			}

			{ // $ bud tool fs tree [dir]
				in := &ToolFsTree{Flag: &framework.Flag{}}
				cli := cli.Command("tree", "list the file tree")
				cli.Flag("embed", "embed assets").Bool(&in.Flag.Embed).Default(false)
				cli.Flag("hot", "hot reloading").Bool(&in.Flag.Hot).Default(false)
				cli.Flag("minify", "minify assets").Bool(&in.Flag.Minify).Default(false)
				cli.Arg("dir").String(&in.Path).Default(".")
				cli.Run(func(ctx context.Context) error { return c.ToolFsTree(ctx, in) })
			}

			{ // $ bud tool fs txtar [dir]
				in := &ToolFsTxtar{Flag: &framework.Flag{}}
				cli := cli.Command("txtar", "generate and print a txtar archive to stdout")
				cli.Arg("dir").String(&in.Path).Default(".")
				cli.Flag("embed", "embed assets").Bool(&in.Flag.Embed).Default(false)
				cli.Flag("hot", "hot reloading").Bool(&in.Flag.Hot).Default(false)
				cli.Flag("minify", "minify assets").Bool(&in.Flag.Minify).Default(false)
				cli.Run(func(ctx context.Context) error { return c.ToolFsTxtar(ctx, in) })
			}
		}

		{ // $ bud tool v8
			in := &ToolV8{}
			cli := cli.Command("v8", "execute Javascript with V8 from stdin")
			cli.Run(func(ctx context.Context) error { return c.ToolV8(ctx, in) })
		}

		{ // $ bud tool cache
			cli := cli.Command("cache", "manage the build cache")

			{ // $ bud tool cache clean
				in := &ToolCacheClean{}
				cli := cli.Command("clean", "clear the cache directory")
				cli.Run(func(ctx context.Context) error { return c.ToolCacheClean(ctx, in) })
			}
		}
	}

	{ // $ bud version
		in := &Version{}
		cli := cli.Command("version", "show the current version")
		cli.Arg("key").String(&in.Key).Default("")
		cli.Run(func(ctx context.Context) error { return c.Version(ctx, in) })
	}

	if err := cli.Parse(ctx, args); err != nil {
		if !errors.Is(err, context.Canceled) {
			return err
		}
	}
	return nil
}

// findModule finds the module for the current working directory.
func (c *CLI) findModule() (*gomod.Module, error) {
	if c.module != nil {
		return c.module, nil
	}
	module, err := gomod.Find(c.Dir)
	if err != nil {
		return nil, err
	}
	c.module = module
	return c.module, nil
}

func (c *CLI) findBudModule() (*gomod.Module, error) {
	if c.budModule != nil {
		return c.budModule, nil
	}
	dirname, err := current.Directory()
	if err != nil {
		return nil, err
	}
	module, err := gomod.Find(dirname)
	if err != nil {
		return nil, err
	}
	c.budModule = module
	return c.budModule, nil
}

// loadLog loads the logger
func (c *CLI) loadLog() (log.Log, error) {
	if c.log != nil {
		return c.log, nil
	}
	level, err := log.ParseLevel(c.Log)
	if err != nil {
		return nil, err
	}
	c.log = log.New(levelfilter.New(console.New(c.Stderr), level))
	return c.log, nil
}

func (c *CLI) openDB(log log.Log, module *gomod.Module) (*dag.DB, error) {
	if c.db != nil {
		return c.db, nil
	}
	db, err := dag.Load(log, module.Directory("bud", "bud.db"))
	if err != nil {
		return nil, err
	}
	c.Closer.Add(db.Close)
	c.db = db
	return c.db, nil
}

func (c *CLI) genFS(cache genfs.Cache, flag *framework.Flag, log log.Log, module *gomod.Module) genfs.FileSystem {
	if c.genfs != nil {
		return c.genfs
	}
	fsys := virtual.Exclude(module, func(path string) bool {
		return path == "bud" || strings.HasPrefix(path, "bud/")
	})
	genfs := genfs.New(cache, fsys, log)
	parser := parser.New(genfs, module)
	injector := di.New(genfs, log, module, parser)
	genfs.FileGenerator("bud/internal/generator/transpiler/transpiler.go", transpiler.New(flag, log, module, parser))
	generator := generator.New(log, module, parser)
	generator.FileGenerators = fileGenerators
	generator.FileServers = fileServers
	genfs.FileGenerator("bud/internal/generator/generator.go", generator)
	genfs.FileGenerator("bud/cmd/afs/main.go", afs.New(flag, injector, log, module))
	return genfs
}

func (c *CLI) listenDev(listenDev string) (socket.Listener, error) {
	if c.DevListener != nil {
		return c.DevListener, nil
	}
	if listenDev == "" {
		listenDev = ":0"
	}
	devLn, err := socket.Listen(listenDev)
	if err != nil {
		return nil, err
	}
	c.Closer.Add(devLn.Close)
	c.DevListener = devLn
	return devLn, nil
}

func (c *CLI) devServer(bus pubsub.Client, devLn net.Listener, flag *framework.Flag, log log.Log, v8 *v8.VM) *budsvr.Server {
	if c.ds != nil {
		return c.ds
	}
	c.ds = budsvr.New(devLn, bus, flag, virtual.List{}, log, v8)
	c.Closer.Add(c.ds.Close)
	return c.ds
}

func (c *CLI) listenAFS(listenAFS string) (socket.Listener, error) {
	if c.AFSListener != nil {
		return c.AFSListener, nil
	}
	if listenAFS == "" {
		listenAFS = ":0"
	}
	afsLn, err := socket.Listen(listenAFS)
	if err != nil {
		return nil, err
	}
	c.Closer.Add(afsLn.Close)
	c.AFSListener = afsLn
	return afsLn, nil
}

func (c *CLI) listenFileAFS(afsLn socket.Listener) (*os.File, error) {
	if c.afsFile != nil {
		return c.afsFile, nil
	}
	fileAFS, err := afsLn.File()
	if err != nil {
		return nil, err
	}
	c.Closer.Add(fileAFS.Close)
	c.afsFile = fileAFS
	return c.afsFile, nil
}

func (c *CLI) bus() pubsub.Client {
	if c.Bus != nil {
		return c.Bus
	}
	c.Bus = pubsub.New()
	return c.Bus
}

func (c *CLI) loadV8() (*v8.VM, error) {
	if c.v8 != nil {
		return c.v8, nil
	}
	v8, err := v8.Load()
	if err != nil {
		return nil, err
	}
	c.Closer.Add(v8.Close)
	c.v8 = v8
	return v8, nil
}

func (c *CLI) listenWeb(listenWeb string) (socket.Listener, error) {
	if c.WebListener != nil {
		return c.WebListener, nil
	}
	webLn, err := socket.ListenUp(listenWeb, 5)
	if err != nil {
		return nil, err
	}
	c.Closer.Add(webLn.Close)
	c.WebListener = webLn
	return webLn, nil
}

func (c *CLI) listenWebFile(webLn socket.Listener) (*os.File, error) {
	if c.webFile != nil {
		return c.webFile, nil
	}
	webFile, err := webLn.File()
	if err != nil {
		return nil, err
	}
	c.Closer.Add(webFile.Close)
	c.webFile = webFile
	return c.webFile, nil
}

func (c *CLI) dialAFS(ctx context.Context, afsLn net.Listener) (*remotefs.Client, error) {
	if c.afsClient != nil {
		// For some reason the RPC client needs to be recycled when the RPC server
		// shuts down.
		if err := c.afsClient.Close(); err != nil {
			return nil, err
		}
	}
	afsClient, err := remotefs.Dial(ctx, afsLn.Addr().String())
	if err != nil {
		return nil, err
	}
	c.Closer.Add(afsClient.Close)
	c.afsClient = afsClient
	return c.afsClient, nil
}

func (c *CLI) command(dir string, name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Env = append([]string{}, c.Env...)
	cmd.Stdin = c.Stdin
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr
	return cmd
}

func (c *CLI) startAFS(ctx context.Context, flag *framework.Flag, module *gomod.Module, afsFile *os.File, devLn net.Listener) (*shell.Process, error) {
	if c.afsProcess != nil {
		p, err := c.afsProcess.Restart(ctx)
		if err != nil {
			return nil, err
		}
		c.afsProcess = p
		c.Closer.Add(c.afsProcess.Close)
		return c.afsProcess, nil
	}
	// Initialize the command
	cmd := c.command(module.Directory(), filepath.Join("bud", "afs"), flag.Flags()...)
	// Setup the environment
	cmd.Env = append(cmd.Env,
		"BUD_DEV_URL="+devLn.Addr().String(),
	)
	// Inject the file under the AFS prefix
	extrafile.Inject(&cmd.ExtraFiles, &cmd.Env, "AFS", afsFile)
	afsProcess, err := shell.Start(ctx, cmd)
	if err != nil {
		return nil, err
	}
	c.Closer.Add(afsProcess.Close)
	c.afsProcess = afsProcess
	return c.afsProcess, nil
}

func (c *CLI) startApp(ctx context.Context, module *gomod.Module, afsLn, devLn net.Listener, webFile *os.File) (*shell.Process, error) {
	if c.appProcess != nil {
		return c.appProcess, nil
	}
	// Initialize the command
	cmd := c.command(module.Directory(), filepath.Join("bud", "app"))
	// Setup the environment
	cmd.Env = append(cmd.Env,
		"BUD_AFS_URL="+afsLn.Addr().String(),
		"BUD_DEV_URL="+devLn.Addr().String(),
	)
	// Inject that file under the WEB prefix
	extrafile.Inject(&cmd.ExtraFiles, &cmd.Env, "WEB", webFile)
	// Start the command
	appProcess, err := shell.Start(ctx, cmd)
	if err != nil {
		return nil, err
	}
	c.Closer.Add(appProcess.Close)
	c.appProcess = appProcess
	return c.appProcess, nil
}

func (c *CLI) prompter(webLn net.Listener) *prompter.Prompter {
	var prompter prompter.Prompter
	c.Stdout = io.MultiWriter(c.Stdout, &prompter.StdOut)
	c.Stderr = io.MultiWriter(c.Stderr, &prompter.StdErr)
	prompter.Init(webLn.Addr().String())
	return &prompter
}

func (c *CLI) readStdin() (string, error) {
	code, err := io.ReadAll(c.Stdin)
	if err != nil {
		return "", err
	}
	script := string(code)
	if script == "" {
		return "", errors.New("missing script to evaluate")
	}
	return script, nil
}
