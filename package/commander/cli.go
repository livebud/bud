package commander

import (
	"context"
	"errors"
	"flag"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/livebud/bud/internal/sig"
)

var ErrCommandNotFound = errors.New("command not found")

type value interface {
	flag.Value
	verify(displayName string) error
}

func New(name, help string) *CLI {
	return &CLI{
		&subcommand{
			nil,
			name,
			"",
			help,
			map[string]*subcommand{},
			[]*Flag{},
			[]*Arg{},
			nil,
			nil,
			false,
			false,
		},
		[]os.Signal{os.Interrupt},
		defaultUsage,
		os.Stdout,
	}
}

type CLI struct {
	*subcommand
	signals []os.Signal
	usage   *template.Template
	writer  io.Writer
}

var _ Command = (*CLI)(nil)

func (c *CLI) Writer(w io.Writer) *CLI {
	c.writer = w
	return c
}

func (c *CLI) Usage(usage *template.Template) *CLI {
	c.usage = usage
	return c
}

func (c *CLI) Signals(signals ...os.Signal) *CLI {
	c.signals = signals
	return c
}

func (c *CLI) Parse(ctx context.Context, args ...string) error {
	// Setup the context
	ctx = sig.Trap(ctx, c.signals...)
	// Setup the flagset
	fset := flag.NewFlagSet(c.name, flag.ContinueOnError)
	fset.SetOutput(io.Discard)
	// Load the root flags
	for _, flag := range c.flags {
		fset.Var(flag.value, flag.name, flag.help)
		// Handle the short flag too
		if flag.short != 0 {
			fset.Var(flag.value, string(flag.short), flag.help)
		}
	}
	// Parse the flags
	if err := fset.Parse(args); err != nil {
		// Print usage if the developer used -h or --help
		if errors.Is(err, flag.ErrHelp) {
			// Handle subcommand help messages
			cmd := c.findOr(fset.Arg(0), c.subcommand)
			return c.printUsage(cmd)
		}
		return err
	}

	// Find the subcommand
	if cmd, ok := c.find(fset.Arg(0)); ok {
		if err := cmd.parse(ctx, fset.Args()[1:]); err != nil {
			if errors.Is(err, flag.ErrHelp) {
				return c.printUsage(cmd)
			}
			return err
		}
		return nil
	}

	if err := c.parse(ctx, fset.Args()); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return c.printUsage(c.subcommand)
		}
		return err
	}
	return nil
}

func (c *CLI) find(path string) (*subcommand, bool) {
	if path == "" {
		return nil, false
	}
	parts := strings.Split(path, ":")
	cmd := c.subcommand
	for _, part := range parts {
		subcommand, ok := cmd.commands[part]
		if !ok {
			return nil, false
		}
		cmd = subcommand
	}
	return cmd, true
}

func (c *CLI) findOr(path string, fallback *subcommand) *subcommand {
	cmd, ok := c.find(path)
	if !ok {
		return fallback
	}
	return cmd
}

func getRoot(c *subcommand) *subcommand {
	if c.parent == nil {
		return c
	}
	return getRoot(c.parent)
}

func (c *CLI) printUsage(s *subcommand) error {
	return c.usage.Execute(c.writer, &usage{s, getRoot(s)})
}
