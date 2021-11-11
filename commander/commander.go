package commander

import (
	"context"
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"text/template"
)

//go:embed usage.gotext
var usage string

var defaultUsage = template.Must(template.New("usage").Funcs(colors).Parse(usage))

func New(name string) *CLI {
	config := &config{"", os.Stderr, defaultUsage, []os.Signal{os.Interrupt}}
	return &CLI{newCommand(config, name, ""), config}
}

type Command struct {
	config *config
	fset   *flag.FlagSet
	run    func(ctx context.Context) error

	// state for the template
	name     string
	usage    string
	commands map[string]*Command
	flags    []*Flag
	args     []*Arg
}

func newCommand(config *config, name, usage string) *Command {
	fset := flag.NewFlagSet(name, flag.ContinueOnError)
	fset.SetOutput(ioutil.Discard)
	return &Command{
		config:   config,
		fset:     fset,
		name:     name,
		usage:    usage,
		commands: map[string]*Command{},
	}
}

type CLI struct {
	root   *Command
	config *config
}

type config struct {
	version  string
	writer   io.Writer
	template *template.Template
	signals  []os.Signal
}

func (c *CLI) Writer(writer io.Writer) *CLI {
	c.config.writer = writer
	return c
}

func (c *CLI) Version(version string) *CLI {
	c.config.version = version
	return c
}

func (c *CLI) Template(template *template.Template) {
	c.config.template = template
}

func (c *CLI) Trap(signals ...os.Signal) {
	c.config.signals = signals
}

func (c *CLI) Parse(args []string) error {
	ctx, cancel := Trap(context.Background(), c.config.signals...)
	defer cancel()
	if err := c.root.parse(ctx, args); err != nil {
		return err
	}
	// Give the caller a chance to handle context cancellations and therefore
	// interrupts specifically.
	return ctx.Err()
}

func (c *CLI) Command(name, usage string) *Command {
	return c.root.Command(name, usage)
}

func (c *CLI) Flag(name, usage string) *Flag {
	return c.root.Flag(name, usage)
}

func (c *CLI) Arg(name, usage string) *Arg {
	return c.root.Arg(name, usage)
}

func (c *CLI) Run(runner func(ctx context.Context) error) {
	c.root.Run(runner)
}

func (c *Command) printUsage() error {
	usage, err := generateUsage(c.config.template, c)
	if err != nil {
		return err
	}
	fmt.Fprint(c.config.writer, usage)
	return nil
}

type value interface {
	flag.Getter
	verify(displayName string) error
}

func (c *Command) parse(ctx context.Context, args []string) error {
	// Set flags
	for _, flag := range c.flags {
		c.fset.Var(flag.value, flag.name, flag.usage)
		if flag.short != 0 {
			c.fset.Var(flag.value, string(flag.short), flag.usage)
		}
	}
	// Parse the arguments
	if err := c.fset.Parse(args); err != nil {
		// Print usage if the developer used -h or --help
		if errors.Is(err, flag.ErrHelp) {
			return c.printUsage()
		}
		return err
	}
	// Verify that all the flags have been set or have default values
	if err := verifyFlags(c.flags); err != nil {
		return err
	}
	// Check if the first argument is a subcommand
	if sub, ok := c.commands[c.fset.Arg(0)]; ok {
		return sub.parse(ctx, c.fset.Args()[1:])
	}
	// Handle the remaining arguments
	numArgs := len(c.args)
	for i, arg := range c.fset.Args() {
		if i >= numArgs {
			return fmt.Errorf("unexpected %s", arg)
		}
		if err := c.args[i].value.Set(arg); err != nil {
			return err
		}
	}
	// Verify that all the args have been set or have default values
	if err := verifyArgs(c.args); err != nil {
		return err
	}
	// Print usage if there's no run function defined
	if c.run == nil {
		if len(c.fset.Args()) == 0 {
			return c.printUsage()
		}
		return fmt.Errorf("unexpected %s", c.fset.Arg(0))
	}
	return c.run(ctx)
}

func (c *Command) Run(runner func(ctx context.Context) error) {
	c.run = runner
}

func (c *Command) Command(name, usage string) *Command {
	if c.commands[name] != nil {
		return c.commands[name]
	}
	cmd := newCommand(c.config, name, usage)
	c.commands[name] = cmd
	return cmd
}

func (c *Command) Arg(name, usage string) *Arg {
	arg := &Arg{
		Name:  name,
		Usage: usage,
	}
	c.args = append(c.args, arg)
	return arg
}

func (c *Command) Flag(name, usage string) *Flag {
	flag := &Flag{
		name:  name,
		usage: usage,
	}
	c.flags = append(c.flags, flag)
	return flag
}
