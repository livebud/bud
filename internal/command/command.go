package command

import (
	"context"
	"fmt"
	"os"
	"time"

	"gitlab.com/mnm/bud/internal/bud"
	"gitlab.com/mnm/bud/package/trace"
	runtime_bud "gitlab.com/mnm/bud/runtime/bud"
)

// Bud command
type Bud struct {
	Flag  runtime_bud.Flag
	Trace bool
	Dir   string
	Args  []string
}

var emptyShutdown = func(*error) {}

func (b *Bud) Tracer(ctx context.Context) (context.Context, func(*error), error) {
	if !b.Trace {
		return ctx, emptyShutdown, nil
	}
	tracer, ctx, err := trace.Serve(ctx)
	if err != nil {
		return nil, nil, err
	}
	shutdown := func(outerError *error) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		trace, err := tracer.Print(ctx)
		if err != nil {
			*outerError = err
			return
		}
		fmt.Fprintf(os.Stderr, "\n%s", trace)
		if err := tracer.Shutdown(ctx); err != nil {
			*outerError = err
			return
		}
	}
	return ctx, shutdown, nil
}

// Run a custom command
func (c *Bud) Run(ctx context.Context) (err error) {
	ctx, shutdown, err := c.Tracer(ctx)
	if err != nil {
		return err
	}
	defer shutdown(&err)
	// Start the trace
	ctx, span := trace.Start(ctx, "running bud")
	defer span.End(&err)
	// Load the compiler
	compiler, err := bud.Find(c.Dir)
	if err != nil {
		return err
	}
	// Compiler the project CLI
	project, err := compiler.Compile(ctx, &c.Flag)
	if err != nil {
		return err
	}
	// Run the custom command
	return project.Execute(ctx, c.Args...)
}
