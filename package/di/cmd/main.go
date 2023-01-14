package main

import (
	"context"
	"os"

	"github.com/livebud/bud/internal/cli/tooldi"
	"github.com/livebud/bud/internal/config"
	"github.com/livebud/bud/package/commander"
	"github.com/livebud/bud/package/log/console"
)

func main() {
	if err := run(context.Background()); err != nil {
		console.Error(err.Error())
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	cfg := config.New()
	cmd := tooldi.New(cfg)
	cli := commander.New("di")
	cli.Flag("name", "name of the function").String(&cmd.Name).Default("Load")
	cli.Flag("dependency", "generate dependency provider").Short('d').Strings(&cmd.Dependencies)
	cli.Flag("external", "mark dependency as external").Short('e').Strings(&cmd.Externals).Optional()
	cli.Flag("map", "map interface types to concrete types").Short('m').StringMap(&cmd.Map).Optional()
	cli.Flag("target", "target import path").Short('t').String(&cmd.Target)
	cli.Flag("hoist", "hoist dependencies that depend on externals").Bool(&cmd.Hoist).Default(false)
	cli.Flag("verbose", "verbose logging").Short('v').Bool(&cmd.Verbose).Default(false)
	cli.Run(cmd.Run)
	return cli.Parse(ctx, os.Args[1:])
}
