package conjure

import (
	"fmt"
	"io/fs"
)

type FileServer interface {
	ServeFile(file *File) error
}

type serverg struct {
	path string // defined generator path
	fn   func(f *File) error
}

func (g *serverg) Generate(target string) (fs.File, error) {
	if target == g.path {
		return nil, fs.ErrInvalid
	}
	file := &File{
		path: target,
	}
	if err := g.fn(file); err != nil {
		return nil, fmt.Errorf("conjure: serve %q . %w", target, err)
	}
	return file.open()
}
