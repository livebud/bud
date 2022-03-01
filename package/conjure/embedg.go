package conjure

import (
	"context"
	"io/fs"
	"time"
)

type Embed struct {
	Data    []byte
	Mode    fs.FileMode
	ModTime time.Time
	Sys     interface{}
}

func (e *Embed) GenerateFile(ctx context.Context, file *File) error {
	file.Data = e.Data
	file.Mode = e.Mode
	file.modTime = e.ModTime
	file.sys = e.Sys
	return nil
}
