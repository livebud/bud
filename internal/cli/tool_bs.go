package cli

import (
	"context"

	"github.com/livebud/bud/framework"
)

type ToolBS struct {
	Flag      *framework.Flag
	ListenDev string
}

func (c *CLI) ToolBS(ctx context.Context, in *ToolBS) error {
	bus := c.bus()

	devLn, err := c.listenDev(in.ListenDev)
	if err != nil {
		return err
	}

	log, err := c.loadLog()
	if err != nil {
		return err
	}
	log.Info("Listening on http://" + devLn.Addr().String())

	v8, err := c.loadV8()
	if err != nil {
		return err
	}

	devServer := c.devServer(bus, devLn, in.Flag, log, v8)
	return devServer.Listen(ctx)
}
