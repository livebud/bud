package mergefs

import (
	"io/fs"

	"github.com/yalue/merged_fs"
)

// Merge the filesystem together
// TODO: this would probably be faster if the dependency supported merging more
// than one filesystem at once
func Merge(first fs.FS, remaining ...fs.FS) (merged fs.FS) {
	merged = first
	for _, fsys := range remaining {
		merged = merged_fs.NewMergedFS(merged, fsys)
	}
	return merged
}
