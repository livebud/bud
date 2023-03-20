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
	return &Viewer{log, module}
}

type Viewer struct {
	log    log.Log
	module *gomod.Module
}

var _ view.Viewer = (*Viewer)(nil)

func (v *Viewer) Register(r *router.Router, pages []*view.Page) {

}

func (v *Viewer) parseTemplate(templatePath string) (*renderer, error) {
	code, err := fs.ReadFile(v.module, path.Join("view", templatePath))
	if err != nil {
		return nil, err
	}
	tpl, err := template.New(templatePath).Parse(string(code))
	if err != nil {
		return nil, err
	}
	return &renderer{tpl}, nil
}

type renderer struct {
	tpl *template.Template
}

func (r *renderer) Render(ctx context.Context, props interface{}) ([]byte, error) {
	out := new(bytes.Buffer)
	// TODO: pass context through
	if err := r.tpl.Execute(out, props); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func (v *Viewer) Render(ctx context.Context, page *view.Page, propMap view.PropMap) ([]byte, error) {
	v.log.Info("rendering gohtml", page.Path)
	entry, err := v.parseTemplate(page.Path)
	if err != nil {
		return nil, err
	}
	layout, err := v.parseTemplate(page.Layout.Path)
	if err != nil {
		return nil, err
	}
	html, err := entry.Render(ctx, propMap[page.Key])
	if err != nil {
		return nil, err
	}
	// TODO: support frames
	// TODO: may want to introduce a slot function
	return layout.Render(ctx, string(html))
}

func (v *Viewer) RenderError(ctx context.Context, page *view.Page, propMap view.PropMap, originalError error) []byte {
	v.log.Info("rendering gohtml", page.Error.Path)
	errorEntry, err := v.parseTemplate(page.Error.Path)
	if err != nil {
		return []byte(fmt.Sprintf("unable to read error template %s to render error: %s. %s", page.Error.Path, err, originalError))
	}
	layout, err := v.parseTemplate(page.Layout.Path)
	if err != nil {
		return []byte(fmt.Sprintf("unable to parse layout template %s to render error: %s. %s", page.Error.Path, err, originalError))
	}
	state := errorState{
		Message: originalError.Error(),
	}
	html, err := errorEntry.Render(ctx, state)
	if err != nil {
		return []byte(fmt.Sprintf("unable to render error template %s to render error: %s. %s", page.Error.Path, err, originalError))
	}
	html, err = layout.Render(ctx, string(html))
	if err != nil {
		return []byte(fmt.Sprintf("unable to render layout template %s to render error: %s. %s", page.Error.Path, err, originalError))
	}
	return html
}

type errorState struct {
	Message string
}

func (v *Viewer) Bundle(ctx context.Context, fsys view.Writable, pages []*view.Page) error {
	fmt.Println("bundling gohtml!")
	return nil
}
