package command

import (
	"bytes"
	"context"
	_ "embed"
	"flag"
	"io"
	"os"
	"strings"
	"text/tabwriter"
	"text/template"

	"github.com/livebud/bud/internal/sig"
)

type Command struct {
	Name     string
	Help     string
	Flags    Flags
	Args     Args
	Commands Commands
	Run      func(context.Context) error
}

func (c *Command) parse(ctx context.Context, args ...string) error {
	return nil
}

type Commands []*Command

func (commands Commands) Usage() (string, error) {
	buf := new(bytes.Buffer)
	tw := tabwriter.NewWriter(buf, 0, 0, 2, ' ', 0)
	for _, cmd := range commands {
		tw.Write([]byte("\t\t" + cmd.Name))
		if cmd.Help != "" {
			tw.Write([]byte("\t" + dim() + cmd.Help + reset()))
		}
		tw.Write([]byte("\n"))
	}
	if err := tw.Flush(); err != nil {
		return "", err
	}
	return strings.TrimSpace(buf.String()), nil
}

// type String struct {
// 	Value    *string
// 	Default  string
// 	Optional bool
// }

// func (s *String) Set(v string) error {
// 	// *s.v = v
// 	return nil
// }

// func (s *String) String() string {
// 	// return *s.v
// 	return ""
// }

// type Bool struct {
// 	Value    *bool
// 	Default  bool
// 	Optional bool
// }

// func (b *Bool) Set(v string) error {
// 	// *s.v = v
// 	return nil
// }

// func (b *Bool) String() string {
// 	// return *s.v
// 	return ""
// }

//go:embed usage.gotext
var usage string

var usageGenerator = template.Must(template.New("usage").Funcs(colors).Parse(usage))

func New(name string, commands ...*Command) *CLI {
	return &CLI{
		Name:     name,
		Commands: commands,
		Signals:  []os.Signal{os.Interrupt},
		Writer:   os.Stdout,
		Usage:    usageGenerator,
	}
}

type CLI struct {
	Name     string
	Commands []*Command
	Signals  []os.Signal
	Writer   io.Writer
	Usage    *template.Template
}

func (c *CLI) Parse(ctx context.Context, args ...string) error {
	ctx = sig.Trap(ctx, c.Signals...)
	fset := flag.NewFlagSet(c.Name, flag.ContinueOnError)
	fset.SetOutput(io.Discard)
	// // Parse the arguments
	// if err := fset.Parse(args); err != nil {
	// 	// Print usage if the developer used -h or --help
	// 	if errors.Is(err, flag.ErrHelp) {
	// 		buf := new(bytes.Buffer)
	// 		if err := c.Usage.Execute(buf, c.root); err != nil {
	// 			return err
	// 		}
	// 		fmt.Fprint(c.Writer, buf.String())
	// 		return nil
	// 	}
	// 	return err
	// }
	// fmt.Println("hello?", args)
	// return nil

	// return c.schema.parse(ctx, args...)
	return nil
}

// func (c *CLI) printSchemaUsage() error {
// 	buf := new(bytes.Buffer)
// 	if err := c.Usage.Execute(buf, &schemaUsage{c.schema}); err != nil {
// 		return err
// 	}
// 	fmt.Fprint(c.Writer, buf.String())
// 	return nil
// }

// type schemaUsage struct {
// 	s *Schema
// }

// func (s *schemaUsage) Name() string {
// 	return s.s.Name
// }

// type usageFlag struct {
// 	f *Flag
// }

// func (s *schemaUsage) Flags() []*usageFlag {
// 	return nil
// }
