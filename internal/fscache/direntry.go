package fscache

import (
	"io/fs"
	"time"
)

type DirEntry struct {
	Base    string // Base name
	Mode    fs.FileMode
	ModTime time.Time
	Sys     interface{}
}

var _ fs.DirEntry = (*DirEntry)(nil)

func (e *DirEntry) Name() string {
	return e.Base
}

func (e *DirEntry) IsDir() bool {
	return e.Mode&fs.ModeDir != 0
}

func (e *DirEntry) Type() fs.FileMode {
	return e.Mode
}

func (e *DirEntry) Info() (fs.FileInfo, error) {
	return &fileInfo{
		name:    e.Base,
		mode:    e.Mode,
		modTime: e.ModTime,
		sys:     e.Sys,
	}, nil
}
