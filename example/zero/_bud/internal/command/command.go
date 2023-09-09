package command

import (
	context "context"
	web "github.com/livebud/bud/example/zero/command/web"
	command "github.com/livebud/bud/example/zero/generator/command"
	log "github.com/livebud/bud/package/log"
)

func New(
	log log.Log,
	webCmd *web.Command,
) *CLI {
	cli := command.New("app")
	webIn := new(web.Serve)
	cli.Run(func(ctx context.Context) error {
		return command.Go(ctx, log,
			func(ctx context.Context) error { return webCmd.GoServe(ctx, webIn) },
		)
	})

	{ // web

		{ // web:serve
			cmd := cli.Command("web:serve", "serve web requests")
			in := new(web.Serve)
			cmd.Run(func(ctx context.Context) error {
				return webCmd.GoServe(ctx, in)
			})
		}
	}

	return cli
}

type CLI = command.CLI
