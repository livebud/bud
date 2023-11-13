package mod

import (
	gopath "path"
	"path/filepath"

	"github.com/livebud/bud/pkg/virt"
	"golang.org/x/mod/modfile"
)

type Module struct {
	dir  string
	file *modfile.File
	virt.FS
}

var _ virt.FS = (*Module)(nil)

func (m *Module) Directory(subpaths ...string) string {
	return filepath.Join(append([]string{m.dir}, subpaths...)...)
}

func (m *Module) Import(subpaths ...string) string {
	modulePath := m.file.Module.Mod.Path
	subPath := gopath.Join(subpaths...)
	if modulePath == "std" {
		return subPath
	}
	return gopath.Join(modulePath, subPath)
}

// Sub returns an FS corresponding to the subtree rooted at fsys's dir.
func (m *Module) Sub(subpaths ...string) *Module {
	dir := m.Directory(subpaths...)
	return &Module{
		dir: dir,
		FS:  virt.OS(dir),
	}
}
