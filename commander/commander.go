package commander

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
)

func New(name string) *CLI {
	return &CLI{newCommand(name)}
}

type Command struct {
	fs   *flag.FlagSet
	run  func(ctx context.Context) error
	cmds map[string]*Command
}

func newCommand(name string) *Command {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(ioutil.Discard)
	return &Command{fs, nil, map[string]*Command{}}
}

type CLI struct {
	root *Command
}

func (c *CLI) Parse(args []string) error {
	return c.root.parse(args)
}

func (c *CLI) Command(name, usage string) *Command {
	return c.root.Command(name, usage)
}

func (c *CLI) Flag(name, usage string) *Flag {
	return c.root.Flag(name, usage)
}

func (c *CLI) Run(runner func(ctx context.Context) error) {
	c.root.Run(runner)
}

func (c *Command) parse(args []string) error {
	if err := c.fs.Parse(args); err != nil {
		return err
	}
	switch c.fs.Arg(0) {
	case "":
		return c.run(context.Background())
	default:
		return fmt.Errorf("unhandled command %s", c.fs.Arg(0))
	}
}

func (c *Command) Run(runner func(ctx context.Context) error) {
	c.run = runner
}

func (c *Command) Command(name, usage string) *Command {
	cmd := newCommand(name)
	c.cmds[name] = cmd
	return cmd
}

func (c *Command) Flag(name, usage string) *Flag {
	return &Flag{cmd: c, name: name, usage: usage}
}

type Flag struct {
	cmd   *Command
	name  string
	usage string
}

func (f *Flag) Int(target *int, defaultValue int) {
	f.cmd.fs.IntVar(target, f.name, defaultValue, f.usage)
}

func (f *Flag) String(target *string, defaultValue string) {
	f.cmd.fs.StringVar(target, f.name, defaultValue, f.usage)
}
