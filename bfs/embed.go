package bfs

import (
	"io/fs"
	"time"
)

type Embed struct {
	Data    []byte      // file content
	Mode    fs.FileMode // FileInfo.Mode
	ModTime time.Time   // FileInfo.ModTime
}

func (ef *Embed) open(f FS, key, relative, path string) (fs.File, error) {
	return &openFile{
		path:    path,
		data:    ef.Data,
		mode:    ef.Mode,
		modTime: ef.ModTime,
	}, nil
}

// Embedded filesystem. This is conceptually similar to fstest.MapFS,
// but doesn't try synthesizing subdirectories and doesn't support
// reading or walking directories.
type EFS map[string]Generator

var _ BFS = (EFS)(nil)

func (efs EFS) Add(fs map[string]Generator) {
	for path, generator := range fs {
		efs[path] = generator
	}
}

// empty fs that implements FS
type emptyfs struct{}

func (emptyfs) Open(name string) (fs.File, error) { return nil, fs.ErrNotExist }
func (emptyfs) link(from, to string, event Event) {}

func (efs EFS) Open(name string) (fs.File, error) {
	generator, ok := efs[name]
	if !ok {
		return nil, fs.ErrNotExist
	}
	return generator.open(emptyfs{}, "", name, name)
}
