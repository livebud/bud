package cli

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/livebud/bud/internal/cli/bud"
	"github.com/livebud/bud/internal/cli/build"
	"github.com/livebud/bud/internal/cli/create"
	"github.com/livebud/bud/internal/cli/newcontroller"
	"github.com/livebud/bud/internal/cli/run"
	"github.com/livebud/bud/internal/cli/toolcache"
	"github.com/livebud/bud/internal/cli/tooldi"
	"github.com/livebud/bud/internal/cli/toolfscat"
	"github.com/livebud/bud/internal/cli/toolfsls"
	"github.com/livebud/bud/internal/cli/toolfstxtar"
	"github.com/livebud/bud/internal/cli/toolv8"
	"github.com/livebud/bud/internal/cli/version"
	"github.com/livebud/bud/internal/pubsub"
	"github.com/livebud/bud/internal/versions"
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

	// Passed in for testing.
	Web socket.Listener // Could be nil
	Hot socket.Listener // Could be nil

}

func (c *CLI) Run(ctx context.Context, args ...string) error {
	// $ bud
	bud := &bud.Command{
		Dir:    c.Dir,
		Env:    c.Env,
		Stdin:  c.Stdin,
		Stdout: c.Stdout,
		Stderr: c.Stderr,
	}
	cli := commander.New("bud").Writer(c.Stdout)
	cli.Flag("chdir", "change the working directory").Short('C').String(&bud.Dir).Default(bud.Dir)
	cli.Flag("help", "show this help message").Short('h').Bool(&bud.Help).Default(false)
	cli.Flag("log", "filter logs with this pattern").Short('L').String(&bud.Log).Default("info")
	cli.Args("args").Strings(&bud.Args).Default(bud.Args...)
	cli.Run(bud.Run)

	{ // $ bud create <dir>
		cmd := create.New(bud)
		cli := cli.Command("create", "create a new app")
		cli.Arg("dir").String(&cmd.Dir)
		cli.Flag("module", "module path for go.mod").String(&cmd.Module).Optional()
		cli.Flag("dev", "link to the development version").Bool(&cmd.Dev).Default(versions.Bud == "latest")
		cli.Run(cmd.Run)
	}

	{ // $ bud run
		cmd := run.New(bud, c.Bus, c.Web, c.Hot)
		cli := cli.Command("run", "run the dev server")
		cli.Flag("embed", "embed assets").Bool(&cmd.Flag.Embed).Default(false)
		cli.Flag("hot", "hot reloading").Bool(&cmd.Flag.Hot).Default(true)
		cli.Flag("minify", "minify assets").Bool(&cmd.Flag.Minify).Default(false)
		cli.Flag("listen", "address to listen to").String(&cmd.Listen).Default(":3000")
		cli.Run(cmd.Run)
	}

	{ // $ bud build
		cmd := build.New(bud)
		cli := cli.Command("build", "build your app into a single binary")
		cli.Flag("embed", "embed assets").Bool(&cmd.Flag.Embed).Default(true)
		cli.Flag("minify", "minify assets").Bool(&cmd.Flag.Minify).Default(true)
		cli.Run(cmd.Run)
	}

	{ // $ bud new
		cli := cli.Command("new", "scaffold code for your app")

		{ // $ bud new controller <name> [actions...]
			cmd := newcontroller.New(bud)
			cli := cli.Command("controller", "scaffold a new controller")
			cli.Arg("path").String(&cmd.Path)
			cli.Args("actions").Strings(&cmd.Actions)
			cli.Run(cmd.Run)
		}

	}

	{ // $ bud tool
		cli := cli.Command("tool", "extra tools")

		{ // $ bud tool di
			cmd := tooldi.New(bud)
			cli := cli.Command("di", "dependency injection generator")
			cli.Flag("dependency", "generate dependency provider").Short('d').Strings(&cmd.Dependencies)
			cli.Flag("external", "mark dependency as external").Short('e').Strings(&cmd.Externals).Optional()
			cli.Flag("map", "map interface types to concrete types").Short('m').StringMap(&cmd.Map).Optional()
			cli.Flag("target", "target import path").Short('t').String(&cmd.Target)
			cli.Flag("hoist", "hoist dependencies that depend on externals").Bool(&cmd.Hoist).Default(false)
			cli.Flag("verbose", "verbose logging").Short('v').Bool(&cmd.Verbose).Default(false)
			cli.Run(cmd.Run)
		}

		{ // $ bud tool fs
			cli := cli.Command("fs", "filesystem tools")

			{ // $ bud tool fs ls [dir]
				cmd := toolfsls.New(bud)
				cli := cli.Command("ls", "list a directory")
				cli.Flag("embed", "embed assets").Bool(&cmd.Flag.Embed).Default(false)
				cli.Flag("hot", "hot reloading").Bool(&cmd.Flag.Hot).Default(true)
				cli.Flag("minify", "minify assets").Bool(&cmd.Flag.Minify).Default(false)
				cli.Arg("dir").String(&cmd.Dir).Default(".")
				cli.Run(cmd.Run)
			}

			{ // $ bud tool fs cat [path]
				// TODO: better align with the unix `cat` command
				cmd := toolfscat.New(bud)
				cli := cli.Command("cat", "print a file")
				cli.Flag("embed", "embed assets").Bool(&cmd.Flag.Embed).Default(false)
				cli.Flag("hot", "hot reloading").Bool(&cmd.Flag.Hot).Default(true)
				cli.Flag("minify", "minify assets").Bool(&cmd.Flag.Minify).Default(false)
				cli.Arg("path").String(&cmd.Path)
				cli.Run(cmd.Run)
			}

			{ // $ bud tool fs txtar [dir]
				cmd := toolfstxtar.New(bud)
				cli := cli.Command("txtar", "generate and print a txtar archive to stdout")
				cli.Arg("dir").String(&cmd.Dir).Default("bud")
				cli.Flag("embed", "embed assets").Bool(&cmd.Flag.Embed).Default(false)
				cli.Flag("hot", "hot reloading").Bool(&cmd.Flag.Hot).Default(true)
				cli.Flag("minify", "minify assets").Bool(&cmd.Flag.Minify).Default(false)
				cli.Run(cmd.Run)
			}
		}

		{ // $ bud tool v8
			cmd := toolv8.New(c.Stdin, c.Stdout)
			cli := cli.Command("v8", "execute Javascript with V8 from stdin")
			cli.Run(cmd.Run)
		}

		{ // $ bud tool cache
			cli := cli.Command("cache", "manage the build cache")

			{ // $ bud tool cache clean
				cmd := toolcache.New(bud)
				cli := cli.Command("clean", "clear the cache directory")
				cli.Run(cmd.Run)
			}
		}
	}

	{ // $ bud version
		cmd := version.New()
		cli := cli.Command("version", "show the current version")
		cli.Arg("key").String(&cmd.Key).Default("")
		cli.Run(cmd.Run)
	}

	// Parse the arguments
	if err := cli.Parse(ctx, args); err != nil {
		// Treat cancellation as a non-error
		if errors.Is(err, context.Canceled) {
			return nil
		}
		return err
	}
	return nil
}
