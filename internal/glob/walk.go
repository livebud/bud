package glob

import (
	"io/fs"

	"github.com/gobwas/glob"
)

func Walk(fsys fs.FS, pattern string, fn fs.WalkDirFunc) error {
	matcher, err := glob.Compile(pattern)
	if err != nil {
		return err
	}
	base := Base(pattern)
	return fs.WalkDir(fsys, base, func(path string, de fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !matcher.Match(path) {
			return nil
		}
		return fn(path, de, err)
	})
}
