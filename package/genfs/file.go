package genfs

import (
	"io/fs"

	"github.com/livebud/bud/package/virtual"
)

type File struct {
	Data   []byte
	path   string
	target string
}

func (f *File) Target() string {
	return f.target
}

func (f *File) Relative() string {
	return relativePath(f.path, f.target)
}

func (f *File) Path() string {
	return f.path
}

func (f *File) Mode() fs.FileMode {
	return fs.FileMode(0)
}

type FileGenerator interface {
	GenerateFile(fsys FS, file *File) error
}

type GenerateFile func(fsys FS, file *File) error

func (fn GenerateFile) GenerateFile(fsys FS, file *File) error {
	return fn(fsys, file)
}

type fileGenerator struct {
	cg    CacheGraph
	fn    func(fsys FS, file *File) error
	genfs fs.FS
	mode  fs.FileMode
	path  string
}

func (f *fileGenerator) Mode() fs.FileMode {
	return f.mode
}

func (f *fileGenerator) Generate(target string) (fs.File, error) {
	if entry, ok := f.cg.Get(target); ok {
		return virtual.New(entry), nil
	}
	file := &File{nil, f.path, target}
	scoped := &scopedFS{f.cg, f.genfs, target}
	if err := f.fn(scoped, file); err != nil {
		return nil, err
	}
	// TODO: Have File implement virtual.Entry and remove this
	vfile := &virtual.File{
		Path: f.path,
		Mode: f.mode,
		Data: file.Data,
	}
	f.cg.Set(target, vfile)
	return virtual.New(vfile), nil
}
