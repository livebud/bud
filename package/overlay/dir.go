package overlay

import (
	"context"

	"github.com/livebud/bud/package/budfs/genfs"
)

type DirGenerator interface {
	GenerateDir(ctx context.Context, fsys F, dir *Dir) error
}

type Dir struct {
	fsys F
	*genfs.Dir
}

func (d *Dir) GenerateFile(path string, fn func(ctx context.Context, fsys F, file *File) error) {
	d.Dir.GenerateFile(path, func(file *genfs.File) error {
		return fn(context.TODO(), d.fsys, &File{file})
	})
}

func (d *Dir) FileGenerator(path string, generator FileGenerator) {
	d.GenerateFile(path, generator.GenerateFile)
}

func (d *Dir) GenerateDir(path string, fn func(ctx context.Context, fsys F, dir *Dir) error) {
	d.Dir.GenerateDir(path, func(dir *genfs.Dir) error {
		return fn(context.TODO(), d.fsys, &Dir{d.fsys, dir})
	})
}

func (d *Dir) DirGenerator(path string, generator DirGenerator) {
	d.GenerateDir(path, generator.GenerateDir)
}
