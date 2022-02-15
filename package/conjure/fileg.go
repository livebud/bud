package conjure

import (
	"fmt"
	"io/fs"
)

type GenerateFile func(file *File) error

func (fn GenerateFile) GenerateFile(file *File) error {
	return fn(file)
}

type FileGenerator interface {
	GenerateFile(file *File) error
}

type fileGenerator struct {
	path string
	gen  FileGenerator
}

func (g *fileGenerator) Generate(target string) (fs.File, error) {
	// Prevents prefixes from matching files
	// (e.g. go.mod/go.mod from matching go.mod)
	if g.path != target {
		return nil, fs.ErrNotExist
	}
	file := &File{
		path: target,
	}
	if err := g.gen.GenerateFile(file); err != nil {
		return nil, fmt.Errorf("conjure: generate %q > %w", target, err)
	}
	return file.open()
}
