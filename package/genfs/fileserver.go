package genfs

import (
	"io/fs"

	"github.com/livebud/bud/package/virtual"
)

type FileServer interface {
	ServeFile(fsys FS, file *File) error
}

type ServeFile func(fsys FS, file *File) error

func (fn ServeFile) ServeFile(fsys FS, file *File) error {
	return fn(fsys, file)
}

type fileServer struct {
	cache Cache
	fn    func(fsys FS, file *File) error
	genfs fs.FS
	path  string
}

var _ generator = (*fileServer)(nil)

func (f *fileServer) Generate(target string) (fs.File, error) {
	if file, err := f.cache.Get(target); err == nil {
		return virtual.Open(file), nil
	}
	// Always return an empty directory if we request the root
	if f.path == target {
		return virtual.Open(&virtual.File{
			Path: f.path,
			Mode: fs.ModeDir,
		}), nil
	}
	scopedFS := &scopedFS{f.cache, f.genfs, target}
	file := &File{nil, f.path, target}
	// g.fsys.log.Fields(log.Fields{
	// 	"target": target,
	// 	"path":   g.node.Path(),
	// }).Debug("budfs: running file server function")
	if err := f.fn(scopedFS, file); err != nil {
		return nil, err
	}
	vfile := &virtual.File{
		Path: target,
		Mode: fs.FileMode(0),
		Data: file.Data,
	}
	if err := f.cache.Set(target, vfile); err != nil {
		return nil, err
	}
	return virtual.Open(vfile), nil
}
