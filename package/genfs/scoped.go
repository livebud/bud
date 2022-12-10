package genfs

import (
	"errors"
	"io/fs"
	"path/filepath"

	"github.com/livebud/bud/internal/glob"
	"github.com/livebud/bud/internal/orderedset"
	"github.com/livebud/bud/internal/valid"
)

type scopedFS struct {
	cache Cache
	genfs fs.FS
	from  string // generator path
}

var _ FS = (*scopedFS)(nil)

// Open implements fs.FS
func (f *scopedFS) Open(name string) (fs.File, error) {
	f.cache.Link(f.from, name)
	file, err := f.genfs.Open(name)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// Watch the paths for changes
func (f *scopedFS) Watch(patterns ...string) error {
	for _, path := range patterns {
		// Not a glob
		if glob.Base(path) == path {
			f.cache.Link(f.from, path)
			continue
		}
		// Compile the pattern into a glob matcher
		matcher, err := glob.Compile(path)
		if err != nil {
			return err
		}
		// Check for changes in the matched paths
		f.cache.Check(f.from, func(path string) bool {
			return matcher.Match(path)
		})
	}
	return nil
}

// Glob implements fs.GlobFS
func (f *scopedFS) Glob(pattern string) (matches []string, err error) {
	// Compile the pattern into a glob matcher
	matcher, err := glob.Compile(pattern)
	if err != nil {
		return nil, err
	}
	// Check for changes in the matched paths
	f.cache.Check(f.from, func(path string) bool {
		return matcher.Match(path)
	})
	// Base is a minor optimization to avoid walking the entire tree
	bases, err := glob.Bases(pattern)
	if err != nil {
		return nil, err
	}
	// Compute the matches for each base
	for _, base := range bases {
		results, err := f.glob(matcher, base)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}
			return nil, err
		}
		matches = append(matches, results...)
	}
	// Deduplicate the matches
	return orderedset.Strings(matches...), nil
}

// ReadDir implements fs.ReadDirFS
func (f *scopedFS) ReadDir(name string) ([]fs.DirEntry, error) {
	// Check for directory changes
	f.cache.Check(f.from, func(path string) bool {
		return path == name || filepath.Dir(path) == name
	})
	des, err := fs.ReadDir(f.genfs, name)
	if err != nil {
		return nil, err
	}
	return des, nil
}

func (f *scopedFS) glob(matcher glob.Matcher, base string) (matches []string, err error) {
	// Walk the directory tree, filtering out non-valid paths
	err = fs.WalkDir(f.genfs, base, valid.WalkDirFunc(func(path string, de fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// If the paths match, add it to the list of matches
		if matcher.Match(path) {
			matches = append(matches, path)
		}
		return nil
	}))
	if err != nil {
		return nil, err
	}
	// return the list of matches
	return matches, nil
}
