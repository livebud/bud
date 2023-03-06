package virtual

import (
	"io/fs"
	"path"

	"github.com/livebud/bud/package/log"
)

// Copy files from one filesystem to another at subpath
func Copy(log log.Log, from fs.FS, to FS, subpaths ...string) error {
	target := path.Join(subpaths...)
	if target == "" {
		target = "."
	}
	return fs.WalkDir(from, target, func(fpath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		} else if fpath == "." {
			return nil
		}
		if d.IsDir() {
			mode := d.Type()
			// Many of the virtual filesystems don't set a mode. Copying these to an
			// actual filesystem will cause permission errors, so we'll use common
			// permissions when not explicitly set.
			if mode == 0 || mode == fs.ModeDir {
				mode = 0755 | fs.ModeDir
			}
			log.Debug("virtual: copying dir", fpath, mode)
			return to.MkdirAll(fpath, mode)
		}
		data, err := fs.ReadFile(from, fpath)
		if err != nil {
			return err
		}
		// Many of the virtual filesystems don't set a mode. Copying these to an
		// actual filesystem will cause permission errors, so we'll use common
		// permissions when not explicitly set.
		mode := d.Type()
		if mode == 0 {
			mode = 0644
		}
		log.Debug("virtual: copying file", fpath, mode)
		return to.WriteFile(fpath, data, mode)
	})
}
