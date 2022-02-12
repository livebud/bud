package bud

import (
	"context"

	"gitlab.com/mnm/bud/internal/buddy/build"
	"gitlab.com/mnm/bud/internal/buddy/expand"
	"gitlab.com/mnm/bud/internal/buddy/generate"
	"gitlab.com/mnm/bud/internal/buddy/run"
	"gitlab.com/mnm/bud/pkg/buddy"
)

func New(kit buddy.Kit) *Driver {
	return &Driver{
		expander:  expand.New(kit),
		generator: generate.New(kit),
		builder:   build.New(kit),
		runner:    run.New(kit),
	}
}

// Driver is a single, public entrypoint that generators and commands can use to
// extend Bud. The driver itself, should not do much, rather it should delegate
// to various internal implementations.
type Driver struct {
	expander  *expand.Command
	generator *generate.Command
	builder   *build.Command
	runner    *run.Command
}

// Expand input
type Expand = expand.Input

// Expand commands and user-defined generators and generate a "project CLI"
func (d *Driver) Expand(ctx context.Context, in *Expand) error {
	return d.expander.Expand(ctx, in)
}

// Generate bud files from the project CLI. Depends on Expand. Generate does not
// run go build on the files.
func (d *Driver) Generate(ctx context.Context, options ...generate.Option) error {
	return d.generator.Generate(ctx, options...)
}

// Build an application from the generated files. Depends on Generate.
func (d *Driver) Build(ctx context.Context, options ...build.Option) error {
	return d.builder.Build(ctx, options...)
}

// Run an application from the generated files and watch for changes.
// Depends on Generate.
func (d *Driver) Run(ctx context.Context, options ...run.Option) error {
	return d.runner.Run(ctx, options...)
}
