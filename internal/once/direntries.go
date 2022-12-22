package once

import (
	"io/fs"
	"sync"
)

// DirEntries ensures we only read the directories once
type DirEntries struct {
	o sync.Once
	v []fs.DirEntry
	e error
}

func (d *DirEntries) Do(fn func() ([]fs.DirEntry, error)) ([]fs.DirEntry, error) {
	d.o.Do(func() { d.v, d.e = fn() })
	return d.v, d.e
}
