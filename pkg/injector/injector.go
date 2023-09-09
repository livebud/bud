package injector

import (
	"github.com/livebud/bud/pkg/bud"
	"github.com/livebud/bud/pkg/cli"
	"github.com/livebud/bud/pkg/di"
	"github.com/livebud/bud/pkg/mod"
)

func New() di.Injector {
	in := di.New()
	di.Provide[*mod.Module](in, mod.Find)
	di.Provide[*cli.CLI](in, cli.Default)
	di.Provide[cli.Parser](in, func(cli *cli.CLI) cli.Parser { return cli })
	di.Provide[*bud.Command](in, bud.New)
	di.Register[*cli.CLI](in, bud.Register)
	return in
}
