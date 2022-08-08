package conjure

import (
	"io/fs"
)

type mountfs struct {
	path string // unused currently
	fsys fs.FS
}

// TODO: it's a bit weird that the mounted filesystem isn't relative to the
// mountpoint, but it's because *overlay.Filesystem depends on module's FS which
// always points to the root.
func (m *mountfs) Generate(target string) (fs.File, error) {
	return m.fsys.Open(target)
}
