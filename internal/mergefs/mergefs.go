package mergefs

import (
	"io/fs"

	"github.com/yalue/merged_fs"
)

func New(fileSystems ...fs.FS) fs.FS {
	return merged_fs.MergeMultiple(fileSystems...)
}
