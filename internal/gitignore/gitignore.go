package gitignore

import (
	"bytes"
	"io/fs"
	"path/filepath"

	"github.com/monochromegane/go-gitignore"
)

func defaultIgnore(path string, isDir bool) bool {
	return isDir && filepath.Base(path) == "node_modules"
}

func New(fsys fs.FS) (ignored func(path string, isDir bool) bool) {
	gi, err := fs.ReadFile(fsys, ".gitignore")
	if err != nil {
		return defaultIgnore
	}
	matcher := gitignore.NewGitIgnoreFromReader(".gitignore", bytes.NewBuffer(gi))
	return matcher.Match
}
