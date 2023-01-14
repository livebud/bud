package cli

import (
	"context"
	"runtime"

	"github.com/livebud/bud/internal/cli3/custom"
	"github.com/livebud/bud/internal/cli3/generate"
	"github.com/livebud/bud/internal/cli3/run"
	"github.com/livebud/bud/internal/versions"
	"github.com/livebud/bud/package/commander"
	"github.com/livebud/bud/package/config"
)

func New(
	config *config.Config,
	custom *custom.Command,
	generate *generate.Command,
	run *run.Command,
) *CLI {
	return &CLI{
		config,
		custom,
		generate,
		run,
	}
}

type CLI struct {
	config   *config.Config
	custom   *custom.Command
	generate *generate.Command
	run      *run.Command
}

func (c *CLI) Parse(ctx context.Context, args ...string) error {
	// Check that we have a valid Go version
	if err := versions.CheckGo(runtime.Version()); err != nil {
		return err
	}

	// $ bud [args...]
	cli := commander.New("bud").Writer(c.config.Stdout)
	cli.Flag("chdir", "change the working directory").Short('C').String(&c.config.Dir).Default(c.config.Dir)
	cli.Flag("help", "show this help message").Short('h').Bool(&c.custom.Help).Default(false)
	cli.Flag("log", "filter logs with this pattern").Short('L').String(&c.config.Log).Default("info")
	cli.Args("args").Strings(&c.custom.Args)
	cli.Run(c.custom.Run)

	{ // $ bud generate [packages...]
		cmd := cli.Command("generate", "generate the code")
		cmd.Flag("embed", "embed assets").Bool(&c.config.Embed).Default(false)
		cmd.Flag("hot", "hot reloading").Bool(&c.config.Hot).Default(true)
		cmd.Flag("minify", "minify assets").Bool(&c.config.Minify).Default(false)
		cmd.Flag("listen-dev", "dev address to listen to").String(&c.config.ListenDev).Default(c.config.ListenDev)
		cmd.Args("packages").Strings(&c.generate.Packages)
		cmd.Run(c.generate.Run)
	}

	{ // $ bud run
		cmd := cli.Command("run", "run your app in dev")
		cmd.Flag("embed", "embed assets").Bool(&c.config.Embed).Default(false)
		cmd.Flag("hot", "hot reloading").Bool(&c.config.Hot).Default(true)
		cmd.Flag("minify", "minify assets").Bool(&c.config.Minify).Default(false)
		cmd.Flag("watch", "watch for changes").Bool(&c.run.Watch).Default(true)
		cmd.Flag("listen", "address to listen to").String(&c.config.ListenWeb).Default(":3000")
		cmd.Flag("listen-dev", "dev address to listen to").String(&c.config.ListenDev).Default(":35729")
		cmd.Run(c.run.Run)
	}

	return cli.Parse(ctx, args)
}
