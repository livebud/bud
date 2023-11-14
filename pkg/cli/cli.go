package cli

import (
	"context"
	"flag"

	"github.com/livebud/bud/internal/oldcli"
)

type Command interface {
	Run(ctx context.Context) error
}

type Help struct {
}

type Helper interface {
	Help(h *Help)
}

type Usage struct{}

func (h *Usage) Run(ctx context.Context) error {
	return flag.ErrHelp
}

func New(name string, cmd Command) *CLI {
	cli := oldcli.New(name, "", cmd)
	return &CLI{cli}
}

type CLI struct {
	cli *oldcli.CLI
}

func (c *CLI) Add(name string, cmd Command) *Subcommand {
	c.cli.Command(name, "", cmd)
	return &Subcommand{}
}

type Subcommand struct {
}

// Help sets the help text for the subcommand.
func (s *Subcommand) Help(help string) *Subcommand {
	return s
}

// Example provides an example for the subcommand.
func (s *Subcommand) Example(help string) *Subcommand {
	return s
}

func (c *CLI) Parse(ctx context.Context, args []string) error {
	return c.cli.Parse(ctx, args...)
}
