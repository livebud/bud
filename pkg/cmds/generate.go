package cmds

import (
	"context"
	"io/fs"

	"github.com/livebud/bud/pkg/cli"
	"github.com/livebud/bud/pkg/logs"
	"github.com/livebud/bud/pkg/mod"
	"github.com/livebud/bud/pkg/virt"
)

func Generate(log logs.Log, fsys fs.FS, mod *mod.Module) *generateCmd {
	return &generateCmd{log, fsys, mod, ""}
}

type generateCmd struct {
	log  logs.Log
	fsys fs.FS
	mod  *mod.Module
	dir  string
}

var _ cli.Subcommand = (*generateCmd)(nil)

func (c *generateCmd) Usage(cmd cli.Command) {
	cmd.Arg("dir").String(&c.dir)
}

func (c *generateCmd) Run(ctx context.Context) error {
	return virt.Sync(c.log, c.fsys, c.mod.In(c.dir))
}
