package overlay

import "gitlab.com/mnm/bud/package/conjure"

type DirGenerator interface {
	GenerateDir(fsys F, dir *Dir) error
}

type Dir struct {
	fsys F
	*conjure.Dir
}

func (d *Dir) GenerateFile(path string, fn func(fsys F, file *File) error) {
	d.Dir.GenerateFile(path, func(file *conjure.File) error {
		return fn(d.fsys, &File{file})
	})
}

func (d *Dir) FileGenerator(path string, generator FileGenerator) {
	d.GenerateFile(path, generator.GenerateFile)
}

func (d *Dir) GenerateDir(path string, fn func(fsys F, dir *Dir) error) {
	d.Dir.GenerateDir(path, func(dir *conjure.Dir) error {
		return fn(d.fsys, &Dir{d.fsys, dir})
	})
}

func (d *Dir) DirGenerator(path string, generator DirGenerator) {
	d.GenerateDir(path, generator.GenerateDir)
}
