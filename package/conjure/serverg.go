package conjure

import (
	"fmt"
	"io/fs"
)

type ServeFile func(file *File) error

func (fn ServeFile) ServeFile(file *File) error {
	return fn(file)
}

type FileServer interface {
	ServeFile(file *File) error
}

type serverg struct {
	path   string // defined generator path
	server FileServer
}

func (g *serverg) Generate(target string) (fs.File, error) {
	if target == g.path {
		return nil, fs.ErrInvalid
	}
	file := &File{
		path: target,
	}
	if err := g.server.ServeFile(file); err != nil {
		return nil, fmt.Errorf("conjure: serve %q > %w", target, err)
	}
	return file.open()
}
