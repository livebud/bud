package gitignore

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	gitignore "github.com/sabhiram/go-gitignore"
)

var alwaysIgnore = []string{
	"node_modules",
	".git",
	".DS_Store",
	// Regardless of if this directory is committed or not, it should be ignored
	// because this will trigger unnecessary rebuilds during development.
	"bud",
}

var defaultIgnores = append([]string{"/bud"}, alwaysIgnore...)

var defaultIgnore = gitignore.CompileIgnoreLines(defaultIgnores...).MatchesPath

func FromFS(fsys fs.FS) (ignore func(path string) bool) {
	code, err := fs.ReadFile(fsys, ".gitignore")
	if err != nil {
		return defaultIgnore
	}
	lines := strings.Split(string(code), "\n")
	lines = append(lines, alwaysIgnore...)
	ignorer := gitignore.CompileIgnoreLines(lines...)
	return ignorer.MatchesPath
}

func From(dir string) (ignore func(path string) bool) {
	code, err := os.ReadFile(filepath.Join(dir, ".gitignore"))
	if err != nil {
		return defaultIgnore
	}
	lines := strings.Split(string(code), "\n")
	lines = append(lines, alwaysIgnore...)
	ignorer := gitignore.CompileIgnoreLines(lines...)
	return ignorer.MatchesPath
}
