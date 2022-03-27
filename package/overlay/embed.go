package overlay

import (
	"context"

	"gitlab.com/mnm/bud/package/conjure"
)

type Embed conjure.Embed

var _ FileGenerator = (*Embed)(nil)
var _ FileServer = (*Embed)(nil)

func (e *Embed) GenerateFile(_ context.Context, _ F, file *File) error {
	return (*conjure.Embed)(e).GenerateFile(file.File)
}

func (e *Embed) ServeFile(_ context.Context, _ F, file *File) error {
	return (*conjure.Embed)(e).ServeFile(file.File)
}
