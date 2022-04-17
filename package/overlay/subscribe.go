package overlay

import "github.com/livebud/bud/internal/pubsub"

func (f *FileSystem) Subscribe(name string) pubsub.Subscription {
	return f.ps.Subscribe(name)
}
