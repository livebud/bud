package virtual

import (
	"io/fs"
	"path"

	"github.com/livebud/bud/package/log"
)

// Copy files from one filesystem to another at subpath
func Copy(log log.Log, from fs.FS, to FS, subpaths ...string) error {
	log = log.Field("fn", "copy")
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
			log.Debug("copy dir", fpath)
			return to.MkdirAll(fpath, d.Type())
		}
		log.Debug("copy file", fpath)
		data, err := fs.ReadFile(from, fpath)
		if err != nil {
			return err
		}
		return to.WriteFile(fpath, data, d.Type())
	})
}
