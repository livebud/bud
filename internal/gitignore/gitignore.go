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
}

var defaultIgnore = gitignore.CompileIgnoreLines(alwaysIgnore...).MatchesPath

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
