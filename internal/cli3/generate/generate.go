package generate

import (
	"context"

	"github.com/livebud/bud/internal/dag"
	"github.com/livebud/bud/package/budhttp/budsvr"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log"
)

func New(
	db *dag.DB,
	ds *budsvr.Server2,
	log log.Log,
	module *gomod.Module,
) *Command {
	return &Command{
		db,
		ds,
		log,
		module,
		nil,
	}
}

type Command struct {
	db       *dag.DB
	ds       *budsvr.Server2
	log      log.Log
	module   *gomod.Module
	Packages []string
}

func (c *Command) Run(ctx context.Context) error {
	c.log.Info("running generate: %s", c.Packages)
	return nil
}
