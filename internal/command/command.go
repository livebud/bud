package command

import (
	"context"
	"fmt"
	"io"

	"github.com/livebud/bud/package/commander"

	"github.com/livebud/bud/internal/compiler"
	"github.com/livebud/bud/package/log/console"
	"github.com/livebud/bud/package/log/filter"

	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log"
)

// Bud command
type Bud struct {
	// Flags
	Dir  string
	Log  string
	Args []string
	Help bool

	// Passed through the subprocesses
	Env    []string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func (c *Bud) Module() (*gomod.Module, error) {
	return gomod.Find(c.Dir)
}

func (c *Bud) Logger() (log.Interface, error) {
	handler, err := filter.Load(console.New(c.Stderr), c.Log)
	if err != nil {
		return nil, err
	}
	return log.New(handler), nil
}

func (c *Bud) Compiler(log log.Interface, module *gomod.Module) *compiler.Bud {
	return &compiler.Bud{
		Module: module,
		Log:    log,
		Stdout: c.Stdout,
		Stderr: c.Stderr,
		Stdin:  c.Stdin,
	}
}

// Run a custom command
func (c *Bud) Run(ctx context.Context) error {
	fmt.Println("running custom command!")
	if c.Help {
		return commander.Usage()
	}
	return nil
}
