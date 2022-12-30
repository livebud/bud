package toolbs

import (
	"context"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/bfs"
	"github.com/livebud/bud/internal/cli/bud"
	"github.com/livebud/bud/internal/pubsub"
	"github.com/livebud/bud/package/budhttp/budsvr"
	v8 "github.com/livebud/bud/package/js/v8"
	"github.com/livebud/bud/package/socket"
)

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
	vm, err := v8.Load()
	if err != nil {
		return err
	}
	// Load the file server
	bfs, err := bfs.Load(c.Flag, log, module)
	if err != nil {
		return err
	}
	bus := pubsub.New()
	budln, err := socket.Listen(":35729")
	if err != nil {
		return err
	}
	defer budln.Close()
	server := budsvr.New(budln, bus, bfs, log, vm)
	log.Info("Listening on http://127.0.0.1:35729")
	return server.Listen(ctx)
}
