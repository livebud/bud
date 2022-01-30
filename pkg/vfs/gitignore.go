package vfs

import (
	"bytes"
	"io/fs"

	"github.com/monochromegane/go-gitignore"
)

type defaultIgnore struct{}

func (defaultIgnore) Match(path string, isDir bool) bool {
	switch {
	case path == "node_modules" && isDir:
		return true
	default:
		return false
	}
}

func GitIgnore(fsys fs.FS) *gitIgnore {
	gi, err := fs.ReadFile(fsys, ".gitignore")
	if err != nil {
		return &gitIgnore{fsys, defaultIgnore{}}
	}
	matcher := gitignore.NewGitIgnoreFromReader(".gitignore", bytes.NewBuffer(gi))
	return &gitIgnore{fsys, matcher}
}

type gitIgnore struct {
	fs.FS
	m gitignore.IgnoreMatcher
}

// ReadDir implements the fs.ReadDirFS to pass the capability down
func (g *gitIgnore) ReadDir(name string) (entries []fs.DirEntry, err error) {
	des, err := fs.ReadDir(g.FS, name)
	if err != nil {
		return nil, err
	}
	for _, de := range des {
		if g.m.Match(de.Name(), de.IsDir()) {
			continue
		}
		entries = append(entries, de)
	}
	return entries, nil
}

func GitIgnoreRW(rw ReadWritable) ReadWritable {
	return &gitIgnoreRW{GitIgnore(rw), rw}
}

type gitIgnoreRW struct {
	*gitIgnore
	rw ReadWritable
}

func (g *gitIgnoreRW) MkdirAll(path string, perm fs.FileMode) error {
	return g.rw.MkdirAll(path, perm)
}
func (g *gitIgnoreRW) WriteFile(name string, data []byte, perm fs.FileMode) error {
	return g.rw.WriteFile(name, data, perm)
}
func (g *gitIgnoreRW) RemoveAll(path string) error {
	return g.rw.RemoveAll(path)
}
