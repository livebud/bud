package overlay

import (
	"context"

	"github.com/livebud/bud/package/budfs/genfs"
)

type Embed genfs.EmbedFile

var _ FileGenerator = (*Embed)(nil)
var _ FileServer = (*Embed)(nil)

func (e *Embed) GenerateFile(_ context.Context, _ F, file *File) error {
	return (*genfs.EmbedFile)(e).GenerateFile(file.File)
}

func (e *Embed) ServeFile(_ context.Context, _ F, file *File) error {
	return (*genfs.EmbedFile)(e).GenerateFile(file.File)
}
