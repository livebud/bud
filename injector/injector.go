package injector

import (
	"context"

	"github.com/livebud/bud/cli"
	"github.com/livebud/bud/cli/routes"
	"github.com/livebud/bud/di"
	"github.com/livebud/bud/internal/mod"
	"github.com/livebud/bud/mux"
)

func New() di.Injector {
	in := di.New()
	di.Loader[*mod.Module](in, loadModule)
	di.Loader[*cli.CLI](in, loadCLI)
	di.Extend[*cli.CLI](in, loadRoutes)
	return in
}

func loadModule(in di.Injector) (*mod.Module, error) {
	return mod.Find()
}

func loadCLI(in di.Injector) (*cli.CLI, error) {
	mod, err := di.Load[*mod.Module](in)
	if err != nil {
		return nil, err
	}
	return cli.Default(mod), nil
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
