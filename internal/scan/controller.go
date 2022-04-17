package scan

import (
	"io/fs"

	"github.com/livebud/bud/internal/valid"
)

func Controllers(fsys fs.FS) Scanner {
	return Dir(fsys, func(de fs.DirEntry) bool {
		if de.IsDir() {
			return valid.Dir(de.Name())
		} else {
			return valid.ControllerFile(de.Name())
		}
	})
}
