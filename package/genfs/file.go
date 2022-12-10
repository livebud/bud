package genfs

import (
	"io/fs"

	"github.com/livebud/bud/package/virtual"
)

type File struct {
	Data []byte

	// Target and path are the same when called within GenerateFile, but not
	// always the same when called within ServeFile
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
	cache Cache
	fn    func(fsys FS, file *File) error
	genfs fs.FS
	path  string
}

func (f *fileGenerator) Generate(target string) (fs.File, error) {
	if target != f.path {
		return nil, formatError(fs.ErrNotExist, "%q path doesn't match %q target", f.path, target)
	}
	if entry, ok := f.cache.Get(target); ok {
		return virtual.New(entry), nil
	}
	file := &File{nil, f.path, target}
	scoped := &scopedFS{f.cache, f.genfs, target}
	if err := f.fn(scoped, file); err != nil {
		return nil, err
	}
	// TODO: Have File implement virtual.Entry and remove this
	vfile := &virtual.File{
		Path: f.path,
		Mode: fs.FileMode(0),
		Data: file.Data,
	}
	f.cache.Set(target, vfile)
	return virtual.New(vfile), nil
}
