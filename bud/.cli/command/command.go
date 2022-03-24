package command

import (
	context "context"
	generator "gitlab.com/mnm/bud/bud/.cli/generator"
	commander "gitlab.com/mnm/bud/package/commander"
	gomod "gitlab.com/mnm/bud/package/gomod"
	build "gitlab.com/mnm/bud/runtime/command/build"
	run "gitlab.com/mnm/bud/runtime/command/run"
	project "gitlab.com/mnm/bud/runtime/project"
)

func New(fsys generator.FS, module *gomod.Module) *CLI {
	return &CLI{fsys, module}
}

type CLI struct {
	fsys generator.FS
	module *gomod.Module
}

func (c *CLI) Parse(ctx context.Context, args ...string) error {
	project := project.New(c.fsys, c.module)
	cli := commander.New("cli")

	{ // cli run
		cmd := &run.Command{Project: project}
		cli := cli.Command("run", "run command")
		cli.Flag("embed", "embed assets").Bool(&cmd.Flag.Embed).Default(false)
		cli.Flag("hot", "hot reload").Bool(&cmd.Flag.Hot).Default(true)
		cli.Flag("minify", "minify assets").Bool(&cmd.Flag.Minify).Default(false)
		cli.Flag("cache", "enable build caching").Bool(&cmd.Flag.Cache).Default(true)
		cli.Flag("port", "port to listen to").String(&cmd.Port).Default(":3000")
		cli.Run(cmd.Run)
	}

	{ // cli build
		cmd := &build.Command{Project: project}
		cli := cli.Command("build", "build command")
		cli.Flag("embed", "embed assets").Bool(&cmd.Flag.Embed).Default(false)
		cli.Flag("hot", "hot reload").Bool(&cmd.Flag.Hot).Default(true)
		cli.Flag("minify", "minify assets").Bool(&cmd.Flag.Minify).Default(false)
		cli.Flag("cache", "enable build caching").Bool(&cmd.Flag.Cache).Default(true)
		cli.Run(cmd.Run)
	}

	return cli.Parse(ctx, args)
}
