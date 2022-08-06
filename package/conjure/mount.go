package conjure

import (
	"io/fs"
	"path/filepath"
)

type mountfs struct {
	path string
	fsys fs.FS
}

func (m *mountfs) Generate(target string) (fs.File, error) {
	// TODO: we shouldn't rely on filepath since paths should be agnostic
	// Unfortunately, there doesn't seem to be a path.Rel()
	rel, err := filepath.Rel(m.path, target)
	if err != nil {
		return nil, err
	}
	return m.fsys.Open(rel)
}
