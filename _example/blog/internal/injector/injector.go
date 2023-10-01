package injector

import (
	"github.com/livebud/bud/cli"
	"github.com/livebud/bud/di"
	"github.com/livebud/bud/example/blog/internal/env"
	"github.com/livebud/bud/injector"
	"github.com/livebud/bud/mux"
)

func New() di.Injector {
	in := injector.New()
	di.Loader[*env.Env](in, loadEnv)
	// di.Loader[*cli.CLI](in, newCLI)
	// di.Loader[cli.Parser](in, cliParser)
	// di.Loader[cli.Command](in, cliCommand)
	di.Loader[*mux.Router](in, newMux)
	di.Extend[*cli.CLI](in, serveCommand)
	return in
}

func loadEnv(in di.Injector) (*env.Env, error) {
	return env.Load[*env.Env]()
}

func newCLI(in di.Injector) (*cli.CLI, error) {
	return cli.New("blog", "blog app"), nil
}

func cliParser(in di.Injector) (cli.Parser, error) {
	return di.Load[*cli.CLI](in)
}

func cliCommand(in di.Injector) (cli.Command, error) {
	return di.Load[*cli.CLI](in)
}

func newMux(in di.Injector) (*mux.Router, error) {
	return mux.New(), nil
}

func serveCommand(in di.Injector, cmd *cli.CLI) error {
	cmd.Command("serve", "start the server")
	// cmd := cli.New("serve", "start the server")
	// cmd.Run(func(ctx context.Context) error {
	// 	return nil
	// })
	return nil
}
