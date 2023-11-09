package markdown

import (
	"io"
	"io/fs"

	"github.com/livebud/bud/pkg/mux"
	"github.com/livebud/bud/pkg/view"
	"github.com/yuin/goldmark"
)

func New(fsys fs.FS) view.Viewer {
	return &viewer{fsys}
}

type viewer struct {
	fsys fs.FS
}

var _ view.Viewer = (*viewer)(nil)

func (v *viewer) Routes(router mux.Router) {
}

func (r *viewer) Render(w io.Writer, path string, data *view.Data) error {
	code, err := fs.ReadFile(r.fsys, path)
	if err != nil {
		return err
	}
	if err := goldmark.Convert(code, w); err != nil {
		return err
	}
	return nil
}
