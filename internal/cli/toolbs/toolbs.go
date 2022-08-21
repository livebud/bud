package toolbs

import (
	"context"

	"github.com/livebud/bud/framework/web/webrt"
	"github.com/livebud/bud/package/socket"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/cli/bud"
)

func New(bud *bud.Command, in *bud.Input) *Command {
	return &Command{
		bud:  bud,
		in:   in,
		Flag: new(framework.Flag),
	}
}

type Command struct {
	bud  *bud.Command
	in   *bud.Input
	Flag *framework.Flag
}

func (c *Command) Run(ctx context.Context) error {
	log, err := bud.Log(c.in.Stdout, c.bud.Log)
	if err != nil {
		return err
	}
	module, err := bud.Module(c.bud.Dir)
	if err != nil {
		return err
	}
	// Load the file server
	bfs, err := bud.FileSystem(ctx, log, module, c.Flag, c.in)
	if err != nil {
		return err
	}
	defer bfs.Close()
	budln, err := socket.Listen(":35729")
	if err != nil {
		return err
	}
	defer budln.Close()
	log.Info("Listening on http://127.0.0.1:35729")
	return webrt.Serve(ctx, budln, bfs)
}
