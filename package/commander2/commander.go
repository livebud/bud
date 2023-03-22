package commander

import (
	"context"
	_ "embed"
	"flag"
	"io"
	"os"
	"text/template"

	"github.com/livebud/bud/internal/sig"
)

type Command interface {
	Command(name, usage string) Command
	Flag(name, usage string) *Flag
	Arg(name string) *Arg
	Args(name string) *Args
	Run(runner func(ctx context.Context) error)
}

//go:embed usage.gotext
var usage string

var defaultUsage = template.Must(template.New("usage").Funcs(colors).Parse(usage))

func Usage() error {
	return flag.ErrHelp
}

func New(name string) *CLI {
	config := &config{"", os.Stdout, defaultUsage, []os.Signal{os.Interrupt}}
	return &CLI{newSubcommand(config, name, ""), config}
}

type CLI struct {
	root   *Subcommand
	config *config
}

var _ Command = (*CLI)(nil)

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

func (c *CLI) Parse(ctx context.Context, args []string) error {
	ctx = sig.Trap(ctx, c.config.signals...)
	if err := c.root.parse(ctx, args); err != nil {
		return err
	}
	// Give the caller a chance to handle context cancellations and therefore
	// interrupts specifically.
	return ctx.Err()
}

func (c *CLI) Command(name, usage string) Command {
	return c.root.Command(name, usage)
}

func (c *CLI) Flag(name, usage string) *Flag {
	return c.root.Flag(name, usage)
}

func (c *CLI) Arg(name string) *Arg {
	return c.root.Arg(name)
}

func (c *CLI) Args(name string) *Args {
	return c.root.Args(name)
}

func (c *CLI) Run(runner func(ctx context.Context) error) {
	c.root.Run(runner)
}
