package virtual

import (
	"io/fs"
	"path"
)

type Map map[string]string

// Map only implements fs.FS because we can't make directories
// or store permission bits in a map.
var _ fs.FS = (Map)(nil)
var _ fs.SubFS = (Map)(nil)

func (m Map) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{
			Op:   "Open",
			Path: name,
			Err:  fs.ErrInvalid,
		}
	}
	return toTree(m).Open(name)
}

func (m Map) Sub(dir string) (fs.FS, error) {
	if !fs.ValidPath(dir) {
		return nil, &fs.PathError{
			Op:   "Sub",
			Path: dir,
			Err:  fs.ErrInvalid,
		}
	}
	return &subMap{dir, m}, nil
}

func toTree(m map[string]string) Tree {
	tree := Tree{}
	for path, data := range m {
		tree[path] = &File{Data: []byte(data)}
	}
	return tree
}

type subMap struct {
	dir string
	m   Map
}

func (s *subMap) Open(name string) (fs.File, error) {
	return s.m.Open(path.Join(s.dir, name))
}
