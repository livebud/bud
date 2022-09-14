package toolbs

import (
	"context"
	"fmt"

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
	return fmt.Errorf("not implemented")
	// log, err := bud.Log(c.in.Stdout, c.bud.Log)
	// if err != nil {
	// 	return err
	// }
	// module, err := bud.Module(c.bud.Dir)
	// if err != nil {
	// 	return err
	// }
	// vm, err := v8.Load()
	// if err != nil {
	// 	return err
	// }
	// // Load the file server
	// servefs, err := bud.FileSystem(log, module, vm, c.Flag)
	// if err != nil {
	// 	return err
	// }
	// bus := pubsub.New()
	// server := budsvr.New(servefs, bus, log, vm)
	// budln, err := socket.Listen(":35729")
	// if err != nil {
	// 	return err
	// }
	// defer budln.Close()
	// log.Info("Listening on http://127.0.0.1:35729")
	// return webrt.Serve(ctx, budln, server)
}
