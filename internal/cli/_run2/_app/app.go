package app

import (
	"context"
	"io/fs"

	"github.com/livebud/bud/internal/pubsub"
	"github.com/livebud/bud/internal/shell"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/js"
	"github.com/livebud/bud/package/log"
)

type Launcher struct {
	Command   *shell.Command
	FS        fs.ReadDirFS
	Log       log.Log
	Module    *gomod.Module
	Publisher pubsub.Publisher
	VM        js.VM
}

func (l *Launcher) Launch(ctx context.Context) (*Process, error) {
	return nil, nil
}

type Process struct {
}
