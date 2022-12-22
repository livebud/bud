package genfs

import (
	"io/fs"
	"sort"

	"github.com/livebud/bud/internal/once"
)

func newDirEntrySet() *dirEntrySet {
	return &dirEntrySet{
		seen: map[string]struct{}{},
	}
}

// dirEntrySet is an ordered set of directory entries.
type dirEntrySet struct {
	seen    map[string]struct{}
	entries []fs.DirEntry
}

func (s *dirEntrySet) Add(entry fs.DirEntry) {
	name := entry.Name()
	if _, ok := s.seen[name]; ok {
		return
	}
	s.seen[name] = struct{}{}
	s.entries = append(s.entries, entry)
}

func (s *dirEntrySet) List() []fs.DirEntry {
	sort.Slice(s.entries, func(i, j int) bool {
		return s.entries[i].Name() < s.entries[j].Name()
	})
	return s.entries
}

func newDirEntry(genfs fs.FS, name string, mode fs.FileMode, path string) *dirEntry {
	return &dirEntry{
		genfs:    genfs,
		name:     name,
		mode:     mode,
		path:     path,
		statOnce: new(once.FileInfo),
	}
}

type dirEntry struct {
	genfs    fs.FS
	name     string
	mode     fs.FileMode
	path     string
	statOnce *once.FileInfo
}

var _ fs.DirEntry = (*dirEntry)(nil)

func (d *dirEntry) Name() string {
	return d.name
}

func (d *dirEntry) Type() fs.FileMode {
	return d.mode
}

func (d *dirEntry) IsDir() bool {
	return d.mode.IsDir()
}

func (d *dirEntry) Info() (fs.FileInfo, error) {
	stat, err := d.statOnce.Do(func() (fs.FileInfo, error) { return fs.Stat(d.genfs, d.path) })
	if err != nil {
		return nil, err
	}
	return stat, nil
}
