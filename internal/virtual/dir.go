package virtual

import (
	"encoding/json"
	"io"
	"io/fs"
)

type Dir struct {
	Path    string
	Entries []fs.DirEntry
	offset  int
}

var _ fs.ReadDirFile = (*Dir)(nil)
var _ fs.File = (*Dir)(nil)
var _ Entry = (*Dir)(nil)

func (d *Dir) MarshalJSON() ([]byte, error) {
	type localType Dir
	// Ensure entries is always set on dir because it's used to differentiate
	// between files and directories.
	if d.Entries == nil {
		d.Entries = []fs.DirEntry{}
	}
	return json.Marshal(localType(*d))
}

func (d *Dir) Close() error {
	return nil
}

func (d *Dir) Stat() (fs.FileInfo, error) {
	return &FileInfo{
		Path:    d.Path,
		ModeDir: true,
	}, nil
}

func (d *Dir) Read(p []byte) (int, error) {
	return 0, &fs.PathError{Op: "read", Path: d.Path, Err: fs.ErrInvalid}
}

func (d *Dir) ReadDir(count int) ([]fs.DirEntry, error) {
	n := len(d.Entries) - d.offset
	if count > 0 && n > count {
		n = count
	}
	if n == 0 && count > 0 {
		return nil, io.EOF
	}
	list := make([]fs.DirEntry, n)
	for i := range list {
		list[i] = d.Entries[d.offset+i]
	}
	d.offset += n
	return list, nil
}

func (d *Dir) Open() fs.File {
	return &Dir{
		Path:    d.Path,
		Entries: d.Entries,
		offset:  0, // reset offset
	}
}
