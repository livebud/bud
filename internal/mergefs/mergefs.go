package mergefs

import (
	"io/fs"

	"github.com/yalue/merged_fs"
)

// New merges filesystems into one. The earlier filesystems override the later
// filesystems when there are conflicts.
func New(fileSystems ...fs.FS) fs.FS {
	return merged_fs.MergeMultiple(fileSystems...)
}
