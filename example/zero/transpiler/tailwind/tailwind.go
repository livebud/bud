package tailwind

import (
	"bytes"

	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/runtime/transpiler"
)

// TODO: switch to markdown. tailwind is more of a custom web generator

type Transpiler struct {
	Log log.Log
}

func (t *Transpiler) GohtmlToGohtml(file *transpiler.File) error {
	t.Log.Info("transpiling", file.Path())
	file.Data = bytes.ReplaceAll(file.Data, []byte("h1"), []byte("h2"))
	return nil
}
