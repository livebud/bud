package goldmark

import (
	"bytes"
	"fmt"

	"github.com/livebud/bud/package/log"
	transpiler "github.com/livebud/bud/runtime/transpiler2"
	"github.com/yuin/goldmark"
)

// TODO: finish markdown transpiler

type Transpiler struct {
	Log log.Log
}

func (t *Transpiler) MdToGohtml(file *transpiler.File) error {
	t.Log.Info("transpiling:", file.Path())
	var html bytes.Buffer
	if err := goldmark.Convert(file.Data, &html); err != nil {
		return fmt.Errorf("goldmark markdown error: %w", err)
	}
	file.Data = html.Bytes()
	return nil
}
