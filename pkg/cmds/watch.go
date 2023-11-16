package cmds

import (
	"context"
	"encoding/json"

	"github.com/livebud/bud/pkg/cli"
	"github.com/livebud/bud/pkg/logs"
	"github.com/livebud/bud/pkg/mod"
	"github.com/livebud/bud/pkg/sse"
	"github.com/livebud/bud/pkg/watcher"
)

func Watch(log logs.Log, mod *mod.Module, sse sse.Publisher) *watchCmd {
	return &watchCmd{log, mod, sse}
}

type watchCmd struct {
	log logs.Log
	mod *mod.Module
	sse sse.Publisher
}

func (c *watchCmd) Usage(cmd cli.Command) {
}

func (c *watchCmd) Run(ctx context.Context) error {
	return watcher.Watch(ctx, ".", func(events *watcher.Events) error {
		eventData, err := json.Marshal(events)
		if err != nil {
			c.log.Error(err)
		}
		if err := c.sse.Publish(ctx, &sse.Event{Data: []byte(eventData)}); err != nil {
			c.log.Error(err)
		}
		return nil
	})
}
