package bfs

import (
	"io/fs"

	"github.com/go-duo/bud/internal/pubsub"
)

type BFS interface {
	Open(name string) (fs.File, error)
	Add(generators map[string]Generator)
	Subscribe(name string) (pubsub.Subscription, error)
}

type Generator interface {
	open(f FS, key, relative, target string) (fs.File, error)
}

type FS interface {
	fs.FS
	link(from, to string, event Event)
}
