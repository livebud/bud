package budfs

import (
	"io/fs"

	"github.com/livebud/bud/internal/virtual"
)

type Embed struct {
	Data []byte
	Mode fs.FileMode
}

var _ Generator = (*Embed)(nil)
var _ FileGenerator = (*Embed)(nil)

func (e *Embed) GenerateFile(fsys FS, file *File) error {
	file.Data = e.Data
	file.Mode = e.Mode
	return nil
}

func (e *Embed) Generate(target string) (fs.File, error) {
	return &virtual.File{
		Name: target,
		Data: e.Data,
		Mode: e.Mode,
	}, nil
}
