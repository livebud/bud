package singleflight

import (
	"io/fs"

	"golang.org/x/sync/singleflight"
)

func New() *Loader {
	return &Loader{}
}

type Loader struct {
	g singleflight.Group
}

func (l *Loader) Load(fsys fs.FS, name string) (fs.File, error) {
	file, err, _ := l.g.Do(name, func() (interface{}, error) {
		return fsys.Open(name)
	})
	if err != nil {
		return nil, err
	}
	return file.(fs.File), nil
}
