package once

import (
	"io/fs"
	"sync"
)

// FileInfo ensures we only read fileinfo once
type FileInfo struct {
	o sync.Once
	v fs.FileInfo
	e error
}

func (d *FileInfo) Do(fn func() (fs.FileInfo, error)) (fs.FileInfo, error) {
	d.o.Do(func() { d.v, d.e = fn() })
	return d.v, d.e
}
