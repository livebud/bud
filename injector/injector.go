package injector

import (
	"context"

	"github.com/livebud/bud/cli"
	"github.com/livebud/bud/cli/routes"
	"github.com/livebud/bud/di"
	"github.com/livebud/bud/mux"
)

func New() di.Injector {
	in := di.New()
	di.Extend[*cli.CLI](in, loadRoutes)
	return in
}

func loadRoutes(in di.Injector, cli *cli.CLI) error {
	routes := routes.New()
	cmd := cli.Command("routes", "list the routes")
	cmd.Run(func(ctx context.Context) error {
		router, err := di.Load[*mux.Router](in)
		if err != nil {
			return err
		}
		return routes.Run(ctx, router)
	})
	return nil
}
