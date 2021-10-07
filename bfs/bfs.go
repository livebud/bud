package bfs

import "io/fs"

type BFS interface {
	Open(name string) (fs.File, error)
	Add(generators map[string]Generator)
}

type Generator interface {
	open(f FS, key, relative, target string) (fs.File, error)
}

type FS interface {
	fs.FS
	link(from, to string, event Event)
}
