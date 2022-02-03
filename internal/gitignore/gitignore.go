package gitignore

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/monochromegane/go-gitignore"
)

func defaultIgnore(path string, isDir bool) bool {
	if !isDir {
		return false
	}
	base := filepath.Base(path)
	return base == "node_modules" || base == ".git"
}

func FromFS(fsys fs.FS) (ignore func(path string, isDir bool) bool) {
	gi, err := fs.ReadFile(fsys, ".gitignore")
	if err != nil {
		return defaultIgnore
	}
	matcher := gitignore.NewGitIgnoreFromReader(".gitignore", bytes.NewBuffer(gi))
	return matcher.Match
}

func From(dir string) (ignore func(path string, isDir bool) bool) {
	code, err := os.ReadFile(filepath.Join(dir, ".gitignore"))
	if err != nil {
		return defaultIgnore
	}
	matcher := gitignore.NewGitIgnoreFromReader(".gitignore", bytes.NewBuffer(code))
	return func(path string, isDir bool) bool {
		rel, err := filepath.Rel(dir, path)
		// Ignore non-relative files
		if err != nil {
			return true
		}
		return matcher.Match(rel, isDir)
	}
}
