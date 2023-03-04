package finder

import (
	"errors"
	"io/fs"

	"github.com/livebud/bud/internal/glob"
	"github.com/livebud/bud/internal/orderedset"
	"github.com/livebud/bud/package/valid"
)

// Find files that match the pattern and are added as entries to the selector
func Find(fsys fs.FS, pattern string, selector func(path string, isDir bool) (entries []string)) (matches []string, err error) {
	matcher, err := glob.Compile(pattern)
	if err != nil {
		return nil, err
	}
	bases, err := glob.Bases(pattern)
	if err != nil {
		return nil, err
	}
	// Compute the matches for each base
	for _, base := range bases {
		// Walk the directory tree, filtering out non-valid paths
		err = fs.WalkDir(fsys, base, valid.WalkDirFunc(func(path string, de fs.DirEntry, err error) error {
			if err != nil {
				if errors.Is(err, fs.ErrNotExist) {
					return nil
				}
				return err
			} else if !matcher.Match(path) {
				return nil
			}
			matched := selector(path, de.IsDir())
			if len(matched) == 0 {
				return nil
			}
			// Ensure all matches paths exist
			if _, err := fs.Stat(fsys, path); err != nil {
				if errors.Is(err, fs.ErrNotExist) {
					return nil
				}
				return err
			}
			matches = append(matches, matched...)
			return nil
		}))
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}
			return nil, err
		}
	}
	return orderedset.Strings(matches...), nil
}
