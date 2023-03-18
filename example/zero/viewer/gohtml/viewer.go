package gohtml

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"path"
	"text/template"

	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/router"
	"github.com/livebud/bud/runtime/view"
)

// TODO: may want a view.FS that is a wrapper around a fs.Sub(module, "view")
func New(log log.Log, module *gomod.Module) *Viewer {
	fmt.Println("new gohtml viewer!")
	return &Viewer{log, module}
}

type Viewer struct {
	log    log.Log
	module *gomod.Module
}

var _ view.Viewer = (*Viewer)(nil)

func (v *Viewer) Register(r *router.Router, pages []*view.Page) {

}

func (v *Viewer) Render(ctx context.Context, page *view.Page, propMap view.PropMap) ([]byte, error) {
	v.log.Info("rendering gohtml", page.Path)
	code, err := fs.ReadFile(v.module, path.Join("view", page.Path))
	if err != nil {
		return nil, err
	}
	tpl, err := template.New(page.Path).Parse(string(code))
	if err != nil {
		return nil, err
	}
	out := new(bytes.Buffer)
	if err := tpl.Execute(out, propMap[page.Key]); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func (v *Viewer) RenderError(ctx context.Context, page *view.Page, propMap view.PropMap, err error) []byte {
	fmt.Println("rendering gohtml error!", page, propMap, err)
	return []byte{}
}

func (v *Viewer) Bundle(ctx context.Context, fsys view.Writable, pages []*view.Page) error {
	fmt.Println("bundling gohtml!")
	return nil
}
