package cli

import (
	"context"
	"errors"
	"runtime"

	"github.com/livebud/bud/internal/cli/bud"
	"github.com/livebud/bud/internal/cli/build"
	"github.com/livebud/bud/internal/cli/create"
	"github.com/livebud/bud/internal/cli/newcontroller"
	"github.com/livebud/bud/internal/cli/run"
	"github.com/livebud/bud/internal/cli/toolbs"
	"github.com/livebud/bud/internal/cli/toolcache"
	"github.com/livebud/bud/internal/cli/tooldi"
	"github.com/livebud/bud/internal/cli/toolfscat"
	"github.com/livebud/bud/internal/cli/toolfsls"
	"github.com/livebud/bud/internal/cli/toolfstree"
	"github.com/livebud/bud/internal/cli/toolfstxtar"
	"github.com/livebud/bud/internal/cli/toolv8"
	"github.com/livebud/bud/internal/cli/version"
	"github.com/livebud/bud/internal/versions"
	"github.com/livebud/bud/package/commander"
)

func New(in *bud.Input) *CLI {
	return &CLI{in: in}
}

// CLI is the Bud CLI. It should not be instantiated directly.
type CLI struct {
	in *bud.Input
}

func (c *CLI) Run(ctx context.Context, args ...string) error {
	// Check that we have a valid Go version
	if err := bud.CheckGoVersion(runtime.Version()); err != nil {
		return err
	}
	// $ bud [args...]
	cmd := bud.New(c.in)
	cli := commander.New("bud").Writer(c.in.Stdout)
	cli.Flag("chdir", "change the working directory").Short('C').String(&cmd.Dir).Default(c.in.Dir)
	cli.Flag("help", "show this help message").Short('h').Bool(&cmd.Help).Default(false)
	cli.Flag("log", "filter logs with this pattern").Short('L').String(&cmd.Log).Default("info")
	cli.Args("args").Strings(&cmd.Args)
	cli.Run(cmd.Run)

	{ // $ bud create <dir>
		cmd := create.New(cmd, c.in)
		cli := cli.Command("create", "create a new app")
		cli.Arg("dir").String(&cmd.Dir)
		cli.Flag("dev", "link to the development version").Short('D').Bool(&cmd.Dev).Default(versions.Bud == "latest")
		cli.Flag("log", "filter logs with this pattern").Short('L').String(&cmd.Log).Default("info")
		cli.Flag("module", "module path for go.mod").String(&cmd.Module).Optional()
		cli.Run(cmd.Run)
	}

	{ // $ bud run
		cmd := run.New(cmd, c.in)
		cli := cli.Command("run", "run the dev server")
		cli.Flag("embed", "embed assets").Bool(&cmd.Flag.Embed).Default(false)
		cli.Flag("hot", "hot reloading").Bool(&cmd.Flag.Hot).Default(true)
		cli.Flag("minify", "minify assets").Bool(&cmd.Flag.Minify).Default(false)
		cli.Flag("listen", "address to listen to").String(&cmd.Listen).Default(":3000")
		cli.Run(cmd.Run)
	}

	{ // $ bud build
		cmd := build.New(cmd, c.in)
		cli := cli.Command("build", "build your app into a single binary")
		cli.Flag("embed", "embed assets").Bool(&cmd.Flag.Embed).Default(true)
		cli.Flag("minify", "minify assets").Bool(&cmd.Flag.Minify).Default(true)
		cli.Run(cmd.Run)
	}

	{ // $ bud new
		cli := cli.Command("new", "scaffold code for your app")

		{ // $ bud new controller <name> [actions...]
			cmd := newcontroller.New(cmd, c.in)
			cli := cli.Command("controller", "scaffold a new controller")
			cli.Arg("path").String(&cmd.Path)
			cli.Args("actions").Strings(&cmd.Actions)
			cli.Run(cmd.Run)
		}

	}

	{ // $ bud tool
		cli := cli.Command("tool", "extra tools")

		{ // $ bud tool bs
			cmd := toolbs.New(cmd, c.in)
			cli := cli.Command("bs", "run the bud server")
			cli.Flag("embed", "embed assets").Bool(&cmd.Flag.Embed).Default(false)
			cli.Flag("hot", "hot reloading").Bool(&cmd.Flag.Hot).Default(true)
			cli.Flag("minify", "minify assets").Bool(&cmd.Flag.Minify).Default(false)
			cli.Run(cmd.Run)
		}

		{ // $ bud tool di
			cmd := tooldi.New(cmd, c.in)
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
				cmd := toolfsls.New(cmd, c.in)
				cli := cli.Command("ls", "list a directory")
				cli.Flag("embed", "embed assets").Bool(&cmd.Flag.Embed).Default(false)
				cli.Flag("hot", "hot reloading").Bool(&cmd.Flag.Hot).Default(true)
				cli.Flag("minify", "minify assets").Bool(&cmd.Flag.Minify).Default(false)
				cli.Arg("dir").String(&cmd.Dir).Default(".")
				cli.Run(cmd.Run)
			}

			{ // $ bud tool fs cat [path]
				// TODO: better align with the unix `cat` command
				cmd := toolfscat.New(cmd, c.in)
				cli := cli.Command("cat", "print a file")
				cli.Flag("embed", "embed assets").Bool(&cmd.Flag.Embed).Default(false)
				cli.Flag("hot", "hot reloading").Bool(&cmd.Flag.Hot).Default(true)
				cli.Flag("minify", "minify assets").Bool(&cmd.Flag.Minify).Default(false)
				cli.Arg("path").String(&cmd.Path)
				cli.Run(cmd.Run)
			}

			{ // $ bud tool fs tree [dir]
				cmd := toolfstree.New(cmd, c.in)
				cli := cli.Command("tree", "list the file tree")
				cli.Flag("embed", "embed assets").Bool(&cmd.Flag.Embed).Default(false)
				cli.Flag("hot", "hot reloading").Bool(&cmd.Flag.Hot).Default(true)
				cli.Flag("minify", "minify assets").Bool(&cmd.Flag.Minify).Default(false)
				cli.Arg("dir").String(&cmd.Dir).Default(".")
				cli.Run(cmd.Run)
			}

			{ // $ bud tool fs txtar [dir]
				cmd := toolfstxtar.New(cmd, c.in)
				cli := cli.Command("txtar", "generate and print a txtar archive to stdout")
				cli.Arg("dir").String(&cmd.Dir).Default("bud")
				cli.Flag("embed", "embed assets").Bool(&cmd.Flag.Embed).Default(false)
				cli.Flag("hot", "hot reloading").Bool(&cmd.Flag.Hot).Default(true)
				cli.Flag("minify", "minify assets").Bool(&cmd.Flag.Minify).Default(false)
				cli.Run(cmd.Run)
			}
		}

		{ // $ bud tool v8
			cmd := toolv8.New(c.in.Stdin, c.in.Stdout)
			cli := cli.Command("v8", "execute Javascript with V8 from stdin")
			cli.Run(cmd.Run)
		}

		{ // $ bud tool cache
			cli := cli.Command("cache", "manage the build cache")

			{ // $ bud tool cache clean
				cmd := toolcache.New(cmd)
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
