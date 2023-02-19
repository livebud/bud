package gohtml

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io/fs"

	"github.com/livebud/bud/package/viewer"
	"github.com/livebud/bud/runtime/transpiler"
)

func New(fsys fs.FS, transpiler transpiler.Interface, pages viewer.Pages) *Viewer {
	return &Viewer{fsys, pages, transpiler}
}

type Viewer struct {
	fsys       fs.FS
	pages      viewer.Pages
	transpiler transpiler.Interface
}

var _ viewer.Viewer = (*Viewer)(nil)

func (v *Viewer) Register(router viewer.Router) {
	fmt.Println("register called")
}

func (v *Viewer) Render(ctx context.Context, key string, props viewer.Props) ([]byte, error) {
	page, ok := v.pages[key]
	if !ok {
		return nil, fmt.Errorf("gohtml: %q. %w", key, viewer.ErrPageNotFound)
	}
	entryCode, err := fs.ReadFile(v.fsys, page.Path)
	if err != nil {
		return nil, fmt.Errorf("gohtml: error reading %q. %w", page.Path, err)
	}
	entryCode, err = v.transpiler.Transpile(page.Path, ".gohtml", entryCode)
	if err != nil {
		return nil, fmt.Errorf("gohtml: error transpiling %q. %w", page.Path, err)
	}
	entryTemplate, err := template.New(page.Path).Parse(string(entryCode))
	if err != nil {
		return nil, fmt.Errorf("gohtml: error parsing %q. %w", page.Path, err)
	}
	entryHTML := new(bytes.Buffer)
	if err := entryTemplate.Execute(entryHTML, props[page.Key]); err != nil {
		return nil, fmt.Errorf("gohtml: error executing %q. %w", page.Path, err)
	}
	return entryHTML.Bytes(), nil
}

func (v *Viewer) RenderError(ctx context.Context, key string, err error, props viewer.Props) []byte {
	return []byte("RenderError not implemented")
}

func (v *Viewer) Bundle(ctx context.Context, out viewer.FS) error {
	return nil
}
