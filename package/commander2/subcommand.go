package commander

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
)

func newSubcommand(config *config, name, usage string) *Subcommand {
	fset := flag.NewFlagSet(name, flag.ContinueOnError)
	fset.SetOutput(io.Discard)
	return &Subcommand{
		config:   config,
		fset:     fset,
		name:     name,
		usage:    usage,
		commands: map[string]*Subcommand{},
	}
}

type Subcommand struct {
	config *config
	fset   *flag.FlagSet
	run    func(ctx context.Context) error
	parsed bool

	// state for the template
	name     string
	usage    string
	commands map[string]*Subcommand
	flags    []*Flag
	args     []*Arg
	restArgs *Args // optional, collects the rest of the args
}

var _ Command = (*Subcommand)(nil)

func (c *Subcommand) printUsage() error {
	usage, err := generateUsage(c.config.template, c)
	if err != nil {
		return err
	}
	fmt.Fprint(c.config.writer, usage)
	return nil
}

type value interface {
	flag.Value
	verify(displayName string) error
}

// Set flags only once
func (c *Subcommand) setFlags() {
	if c.parsed {
		return
	}
	c.parsed = true
	for _, flag := range c.flags {
		c.fset.Var(flag.value, flag.name, flag.usage)
		if flag.short != 0 {
			c.fset.Var(flag.value, string(flag.short), flag.usage)
		}
	}
}

func (c *Subcommand) parse(ctx context.Context, args []string) error {
	// Set flags
	c.setFlags()
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
	restArgs := c.fset.Args()
loop:
	for i, arg := range restArgs {
		if i >= numArgs {
			if c.restArgs == nil {
				return fmt.Errorf("unexpected %s", arg)
			}
			// Loop over the remaining unset args, appending them to restArgs
			for _, arg := range restArgs[i:] {
				c.restArgs.value.Set(arg)
			}
			break loop
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
		if len(restArgs) == 0 {
			return c.printUsage()
		}
		return fmt.Errorf("unexpected %s", c.fset.Arg(0))
	}
	if err := c.run(ctx); err != nil {
		// Support explicitly printing usage
		if errors.Is(err, flag.ErrHelp) {
			return c.printUsage()
		}
		return err
	}
	return nil
}

func (c *Subcommand) Run(runner func(ctx context.Context) error) {
	c.run = runner
}

func (c *Subcommand) Command(name, usage string) Command {
	if c.commands[name] != nil {
		return c.commands[name]
	}
	cmd := newSubcommand(c.config, name, usage)
	c.commands[name] = cmd
	return cmd
}

func (c *Subcommand) Arg(name string) *Arg {
	arg := &Arg{
		Name: name,
	}
	c.args = append(c.args, arg)
	return arg
}

func (c *Subcommand) Args(name string) *Args {
	if c.restArgs != nil {
		// Panic is okay here because settings commands should be done during
		// initialization. We want to fail fast for invalid usage.
		panic("commander: you can only use cmd.Args(name, usage) once per command")
	}
	args := &Args{
		Name: name,
	}
	c.restArgs = args
	return args
}

func (c *Subcommand) Flag(name, usage string) *Flag {
	flag := &Flag{
		name:  name,
		usage: usage,
	}
	c.flags = append(c.flags, flag)
	return flag
}
