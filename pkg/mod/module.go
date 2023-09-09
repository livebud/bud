package mod

import (
	gopath "path"
	"path/filepath"

	"golang.org/x/mod/modfile"
)

type Module struct {
	dir  string
	file *modfile.File
}

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
