package virtual

import (
	"io/fs"
	"path"
)

func Exclude(fsys FS, fn func(path string) bool) FS {
	return &exclude{fsys, fn}
}

type exclude struct {
	FS
	fn func(path string) bool
}

func (e *exclude) Open(path string) (fs.File, error) {
	if e.fn(path) {
		return nil, fs.ErrNotExist
	}
	return e.FS.Open(path)
}

func (e *exclude) ReadDir(dir string) (results []fs.DirEntry, err error) {
	if e.fn(dir) {
		return nil, fs.ErrNotExist
	}
	des, err := fs.ReadDir(e.FS, dir)
	if err != nil {
		return nil, err
	}
	for _, de := range des {
		if e.fn(path.Join(dir, de.Name())) {
			continue
		}
		results = append(results, de)
	}
	return results, nil
}
