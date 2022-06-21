package run

import (
	"context"
	"path/filepath"

	"github.com/livebud/bud/internal/pubsub"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/watcher"
)

type watchfs struct {
	bus pubsub.Client
	dir string
	log log.Interface
}

func (w *watchfs) Run(ctx context.Context) error {
	return watcher.Watch(ctx, w.dir, func(paths []string) error {
		topic := merge(paths)
		w.log.Debug("run: new watch event", "paths", paths, "topic", topic)
		w.bus.Publish(topic, nil)
		return nil
	})
}

// Merge all events into one event to avoid any refresh trashing
func merge(paths []string) string {
	for _, path := range paths {
		if filepath.Ext(path) == ".go" {
			return "watch:backend:update"
		}
	}
	return "watch:frontend:update"
}
