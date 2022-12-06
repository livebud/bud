package generate

import (
	"context"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/buddy"
	"github.com/livebud/bud/internal/cli/bud"
)

// New command for bud generate
func New(bud *bud.Command, in *bud.Input) *Command {
	return &Command{
		bud: bud,
		in:  in,
		Flag: &framework.Flag{
			Env:    in.Env,
			Stderr: in.Stderr,
			Stdin:  in.Stdin,
			Stdout: in.Stdout,
		},
	}
}

// Command for running bud generate
type Command struct {
	bud  *bud.Command
	in   *bud.Input
	Flag *framework.Flag
	Args []string
}

// Run the generate command
func (c *Command) Run(ctx context.Context) error {
	bud, err := buddy.Load(ctx, &buddy.Input{
		Embed:  c.Flag.Embed,
		Minify: c.Flag.Minify,
		Hot:    c.Flag.Hot,
		Log:    c.bud.Log,
		Dir:    c.bud.Dir,
		Stdin:  c.in.Stdin,
		Stdout: c.in.Stdout,
		Stderr: c.in.Stderr,
		Env:    c.in.Env,
	})
	if err != nil {
		return err
	}
	appFS, err := bud.Generate(ctx, c.Args...)
	if err != nil {
		return err
	}
	defer appFS.Close()
	return nil
}
