package mod

import (
	"io/fs"
	"os"
)

// osfs is an fs.FS that just defers to the operating system. This makes this
// fs.FS non-compliant as it accepts absolute paths (e.g. /a/b/c). I hope this
// issue may clean this up: github.com/golang/go/issues/44279
type osfs struct{}

func (osfs) Open(name string) (fs.File, error) {
	return os.Open(name)
}

func (osfs) Stat(name string) (fs.FileInfo, error) {
	return os.Stat(name)
}

func (osfs) ReadDir(name string) ([]fs.DirEntry, error) {
	return os.ReadDir(name)
}

func (o osfs) Sub(name string) (fs.FS, error) {
	return os.DirFS(name), nil
}

// This subFS doesn't conform to fs.Sub(). We allow dir to be an absolute path
// that starts with /.
func subFS(fsys fs.FS, dir string) (fs.FS, error) {
	if dir == "." {
		return fsys, nil
	}
	if fsys, ok := fsys.(fs.SubFS); ok {
		return fsys.Sub(dir)
	}
	return fs.Sub(fsys, dir)
}
