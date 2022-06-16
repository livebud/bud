package cli

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/livebud/bud/internal/command"
	"github.com/livebud/bud/internal/command/build"
	"github.com/livebud/bud/internal/command/run"
	"github.com/livebud/bud/internal/pubsub"
	"github.com/livebud/bud/package/commander"
	"github.com/livebud/bud/package/log/console"
	"github.com/livebud/bud/package/socket"
)

func Run(ctx context.Context, args ...string) int {
	// Default CLI
	cli := &CLI{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Bus:    pubsub.New(),
		Env: []string{
			"HOME=" + os.Getenv("HOME"),
			"PATH=" + os.Getenv("PATH"),
			"GOPATH=" + os.Getenv("GOPATH"),
			"TMPDIR=" + os.TempDir(),
		},
	}
	// Run the cli
	if err := cli.Run(ctx, args...); err != nil {
		if errors.Is(err, context.Canceled) {
			return 0
		}
		console.Error(err.Error())
		return 1
	}
	return 0
}

// CLI is the Bud CLI. It should not be instantiated directly.
type CLI struct {
	Dir    string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
	Bus    pubsub.Client
	Env    []string

	// Used for testing.
	Web socket.Listener // Could be nil
	Hot socket.Listener // Could be nil
}

func (c *CLI) Run(ctx context.Context, args ...string) error {
	// $ bud
	bud := &command.Bud{
		Dir:    c.Dir,
		Env:    c.Env,
		Stdin:  c.Stdin,
		Stdout: c.Stdout,
		Stderr: c.Stderr,
	}
	cli := commander.New("bud").Writer(c.Stdout)
	cli.Flag("chdir", "change the working directory").Short('C').String(&bud.Dir).Default(bud.Dir)
	cli.Flag("help", "show the help message").Short('h').Bool(&bud.Help).Default(false)
	cli.Flag("log", "log pattern").Short('L').String(&bud.Log).Default("info")
	cli.Args("args").Strings(&bud.Args).Default(bud.Args...)
	cli.Run(bud.Run)

	{ // $ bud run
		cmd := run.New(bud, c.Bus, c.Web, c.Hot)
		cli := cli.Command("run", "run the development server")
		cli.Flag("embed", "embed assets").Bool(&cmd.Flag.Embed).Default(false)
		cli.Flag("hot", "hot reloading").Bool(&cmd.Flag.Hot).Default(true)
		cli.Flag("minify", "minify assets").Bool(&cmd.Flag.Minify).Default(false)
		cli.Flag("listen", "address to listen to").String(&cmd.Listen).Default(":3000")
		cli.Run(cmd.Run)
	}

	{ // $ bud build
		cmd := build.New(bud)
		cli := cli.Command("build", "build the application into a single binary")
		cli.Flag("embed", "embed assets").Bool(&cmd.Flag.Embed).Default(true)
		cli.Flag("minify", "minify assets").Bool(&cmd.Flag.Minify).Default(true)
		cli.Run(cmd.Run)
	}

	return cli.Parse(ctx, args)
}

// // Default command
// var Default = &Command{
// 	Env:    defaultEnv,
// 	Stdin:  os.Stdin,
// 	Stdout: os.Stdout,
// 	Stderr: os.Stderr,
// }

// // Bud CLI
// type Bud struct {
// 	Dir    string
// 	Env    []string
// 	Stdin  io.Reader
// 	Stdout io.Writer
// 	Stderr io.Writer

// 	// Used for testing
// 	Web *os.File // Optional net.Listener
// 	Hot *os.File // Optional net.Listener

// 	// Private CLI flags
// 	help bool
// 	args []string
// }

// func (b *Bud) Start(ctx context.Context, args ...string) error {
// 	cli := commander.New("bud").Writer(c.Stdout)
// 	cli.Flag("chdir", "change the working directory").Short('C').String(&c.Dir).Default(c.Dir)
// 	cli.Flag("help", "show the help message").Short('h').Bool(&c.help).Default(false)
// 	cli.Args("args").Strings(&c.args).Default(c.args...)
// 	cli.Run(c.beforeRun)

// 	{ // bud run
// 		cmd := new(run.Command)
// 		cli := cli.Command("run", "run the development server")
// 		cli.Flag("embed", "embed assets").Bool(&cmd.Embed).Default(false)
// 		cli.Flag("hot", "hot reloading").String(&cmd.Hot).Default(":35729")
// 		cli.Flag("minify", "minify assets").Bool(&cmd.Minify).Default(false)
// 		cli.Flag("listen", "address to listen to").String(&cmd.Listen).Default(":3000")
// 		cli.Run(c.beforeRun(cmd))
// 	}

// 	{ // bud build
// 		cmd := new(run.Command)
// 		cli := cli.Command("run", "run the development server")
// 		cli.Flag("embed", "embed assets").Bool(&cmd.Embed).Default(false)
// 		cli.Flag("hot", "hot reloading").String(&cmd.Hot).Default(":35729")
// 		cli.Flag("minify", "minify assets").Bool(&cmd.Minify).Default(false)
// 		cli.Flag("listen", "address to listen to").String(&cmd.Listen).Default(":3000")
// 		cli.Run(c.beforeRun(cmd))
// 	}

// 	return cli.Parse(ctx, args)

// 	// { // $ bud create <dir>
// 	// 	cmd := &create.Command{}
// 	// 	cli := cli.Command("create", "create a new project")
// 	// 	cli.Arg("dir").String(&cmd.Dir)
// 	// 	cli.Run(c.create(cmd))
// 	// }

// 	// { // $ bud tool
// 	// 	cli := cli.Command("tool", "extra tools")

// 	// 	{ // $ bud tool di
// 	// 		// TODO: move into the project CLI since it depends on being within a Go
// 	// 		// module.
// 	// 		cmd := &tool_di.Command{}
// 	// 		cli := cli.Command("di", "dependency injection generator")
// 	// 		cli.Flag("dependency", "generate dependency provider").Short('d').Strings(&cmd.Dependencies)
// 	// 		cli.Flag("external", "mark dependency as external").Short('e').Strings(&cmd.Externals).Optional()
// 	// 		cli.Flag("map", "map interface types to concrete types").Short('m').StringMap(&cmd.Map).Optional()
// 	// 		cli.Flag("target", "target import path").Short('t').String(&cmd.Target)
// 	// 		cli.Flag("hoist", "hoist dependencies that depend on externals").Bool(&cmd.Hoist).Default(false)
// 	// 		cli.Flag("verbose", "verbose logging").Short('v').Bool(&cmd.Verbose).Default(false)
// 	// 		cli.Run(c.tool_di(cmd))
// 	// 	}

// 	// 	{ // $ bud tool v8
// 	// 		cmd := &tool_v8.Command{Stdin: c.Stdin, Stdout: c.Stdout}
// 	// 		cli := cli.Command("v8", "Execute Javascript with V8 from stdin")
// 	// 		cli.Run(cmd.Run)

// 	// 		{ // $ bud tool v8 serve
// 	// 			cmd := &tool_v8_serve.Command{}
// 	// 			cli := cli.Command("serve", "Serve from a V8 server during development")
// 	// 			cli.Run(cmd.Run)
// 	// 		}
// 	// 	}

// 	// 	{ // $ bud tool cache
// 	// 		cli := cli.Command("cache", "Manage the build cache")

// 	// 		{ // $ bud tool cache clean
// 	// 			// TODO: move into the project CLI since it depends on a project
// 	// 			// existing anyways.
// 	// 			cmd := &tool_cache_clean.Command{}
// 	// 			cli := cli.Command("clean", "Clear the cache directory")
// 	// 			cli.Run(c.tool_cache_clean(cmd))
// 	// 		}
// 	// 	}
// 	// }

// 	// { // $ bud version
// 	// 	cmd := version.Command{}
// 	// 	cli := cli.Command("version", "Show package versions")
// 	// 	cli.Arg("key").String(&cmd.Key).Default("")
// 	// 	cli.Run(cmd.Run)
// 	// }

// 	// return cli.Parse(ctx, args)
// }

// func (b *Bud) beforeRun(cmd *run.Command) func(ctx context.Context) error {
// 	return func(ctx context.Context) error {
// 		return cmd.Run(ctx)
// 	}
// }

// func (b *Bud) runCustom(ctx context.Context) error {
// 	return nil
// }
