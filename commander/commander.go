package commander

import (
	"bytes"
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
	config := &config{os.Stderr, defaultUsage}
	return &CLI{newCommand(config, name, ""), config}
}

type Command struct {
	config *config
	fset   *flag.FlagSet
	run    func(ctx context.Context) error

	// state for the template
	// TODO: hide
	Name     string
	Usage    string
	Commands map[string]*Command
	Flags    []*Flag
	Args     []*Arg
}

func newCommand(config *config, name, usage string) *Command {
	fset := flag.NewFlagSet(name, flag.ContinueOnError)
	fset.SetOutput(ioutil.Discard)
	return &Command{
		config:   config,
		fset:     fset,
		Name:     name,
		Usage:    usage,
		Commands: map[string]*Command{},
	}
}

type CLI struct {
	root   *Command
	config *config
}

type config struct {
	writer   io.Writer
	template *template.Template
}

func (c *CLI) Writer(writer io.Writer) *CLI {
	c.config.writer = writer
	return c
}

func (c *CLI) Template(template *template.Template) {
	c.config.template = template
}

func (c *CLI) Parse(args []string) error {
	return c.root.parse(context.TODO(), args)
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

func (c *Command) usage() error {
	buf := new(bytes.Buffer)
	if err := c.config.template.Execute(buf, c); err != nil {
		return err
	}
	fmt.Fprint(c.config.writer, buf.String())
	return nil
}

func (c *Command) parse(ctx context.Context, args []string) error {
	// Set flags
	for _, flag := range c.Flags {
		c.fset.Var(flag.value, flag.name, flag.name)
	}
	// Parse the arguments
	if err := c.fset.Parse(args); err != nil {
		// Print usage if the developer used -h or --help
		if errors.Is(err, flag.ErrHelp) {
			return c.usage()
		}
		return err
	}
	// Verify that all the flags have been set or have default values
	if err := verifyFlags(c.Flags); err != nil {
		return err
	}
	// Check if the first argument is a subcommand
	if sub, ok := c.Commands[c.fset.Arg(0)]; ok {
		return sub.parse(ctx, c.fset.Args()[1:])
	}
	// Handle the remaining arguments
	for i, arg := range c.Args {
		if err := arg.value.Set(c.fset.Arg(i)); err != nil {
			return err
		}
	}
	// Print usage if there's no run function defined
	if c.run == nil {
		if len(c.fset.Args()) == 0 {
			return c.usage()
		}
		return fmt.Errorf("unexpected %s", c.fset.Arg(0))
	}
	return c.run(ctx)
}

func (c *Command) Run(runner func(ctx context.Context) error) {
	c.run = runner
}

func (c *Command) Command(name, usage string) *Command {
	cmd := newCommand(c.config, name, usage)
	c.Commands[name] = cmd
	return cmd
}

func (c *Command) Arg(name, usage string) *Arg {
	arg := &Arg{
		name:  name,
		usage: usage,
	}
	c.Args = append(c.Args, arg)
	return arg
}

type Arg struct {
	name  string
	usage string
	value flag.Getter
}

func (a *Arg) Int(target *int) *Int {
	value := &Int{target, nil}
	a.value = &intValue{value}
	return value
}

func (a *Arg) String(target *string) *String {
	value := &String{target, nil}
	a.value = &stringValue{inner: value}
	return value
}

func (c *Command) Flag(name, usage string) *Flag {
	flag := &Flag{
		name:  name,
		usage: usage,
	}
	c.Flags = append(c.Flags, flag)
	return flag
}
