package commander

import (
	"context"
	"flag"
	"fmt"
	"io"
)

type Command interface {
	Command(name, help string) Command
	Flag(name, help string) *Flag
	Arg(name string) *Arg
	Args(name string) *Args
	Run(fn func(ctx context.Context) error)
	Advanced() Command
	Hidden() Command
}

type subcommand struct {
	parent   *subcommand
	name     string
	full     string
	help     string
	commands map[string]*subcommand
	flags    []*Flag
	args     []*Arg
	run      func(ctx context.Context) error
	restArgs *Args // optional, collects the rest of the args
	advanced bool
	hidden   bool
}

var _ Command = (*subcommand)(nil)

func (c *subcommand) Command(name, help string) Command {
	full := name
	if c.full != "" {
		full = c.full + ":" + name
	}
	command := &subcommand{
		c,
		name,
		full,
		help,
		map[string]*subcommand{},
		nil,
		nil,
		nil,
		nil,
		false,
		false,
	}
	c.commands[name] = command
	return command
}

func (c *subcommand) Flag(name, help string) *Flag {
	flag := &Flag{name, help, 0, nil}
	c.flags = append(c.flags, flag)
	return flag
}

func (c *subcommand) Arg(name string) *Arg {
	arg := &Arg{name, nil}
	c.args = append(c.args, arg)
	return arg
}

func (c *subcommand) Args(name string) *Args {
	args := &Args{name, nil}
	c.restArgs = args
	return args
}

func (c *subcommand) Run(fn func(ctx context.Context) error) {
	c.run = fn
}

func (c *subcommand) Advanced() Command {
	c.advanced = true
	return c
}

func (c *subcommand) Hidden() Command {
	c.hidden = true
	return c
}

func (c *subcommand) extract(fset *flag.FlagSet, arguments []string) (args []string, err error) {
	for len(arguments) > 0 {
		if err := fset.Parse(arguments); err != nil {
			return nil, err
		}
		if fset.NArg() == 0 {
			return args, nil
		}
		args = append(args, fset.Arg(0))
		arguments = fset.Args()[1:]
	}
	return args, nil
}

func (c *subcommand) parse(ctx context.Context, arguments []string) error {
	fset := flag.NewFlagSet(c.full, flag.ContinueOnError)
	fset.SetOutput(io.Discard)
	for _, flag := range c.flags {
		fset.Var(flag.value, flag.name, flag.help)
		if flag.short != 0 {
			fset.Var(flag.value, string(flag.short), flag.help)
		}
	}
	args, err := c.extract(fset, arguments)
	if err != nil {
		return err
	}
	numArgs := len(c.args)
	for i, arg := range args {
		// Handle variadic arguments
		if i >= numArgs {
			if c.restArgs == nil {
				return fmt.Errorf("%w %s", ErrCommandNotFound, arg)
			}
			for _, arg := range args[i:] {
				if err := c.restArgs.value.Set(arg); err != nil {
					return err
				}
			}
			break
		}
		// Otherwise, set the argument as normal
		if err := c.args[i].value.Set(arg); err != nil {
			return err
		}
	}
	// Verify that all the args have been set or have default values
	if err := verifyArgs(c.args); err != nil {
		return err
	}
	// Also verify rest args if we have any
	if c.restArgs != nil {
		if err := c.restArgs.verify(c.restArgs.Name); err != nil {
			return err
		}
	}
	// Print usage if there's no run function defined
	if c.run == nil {
		return flag.ErrHelp
	}
	// Verify that all the flags have been set or have default values
	if err := verifyFlags(c.flags); err != nil {
		return err
	}
	// Run the command
	return c.run(ctx)
}
