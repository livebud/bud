package overlay

import "gitlab.com/mnm/bud/internal/pubsub"

func (f *FileSystem) Subscribe(name string) pubsub.Subscription {
	return f.ps.Subscribe(name)
}
