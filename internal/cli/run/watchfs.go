package run

import (
	"context"
	"path/filepath"

	"github.com/livebud/bud/internal/pubsub"
	"github.com/livebud/bud/package/watcher"
)

type watchfs struct {
	bus pubsub.Client
	dir string
}

func (w *watchfs) Run(ctx context.Context) error {
	return watcher.Watch(ctx, w.dir, func(paths []string) error {
		for _, path := range paths {
			if filepath.Ext(path) == ".go" {

			}
		}
		return nil
	})
}
