package gohtml

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"path"
	"text/template"

	"github.com/livebud/bud/example/zero/bud/pkg/transpiler"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/router"
	"github.com/livebud/bud/runtime/view"
)

// TODO: may want a view.FS that is a wrapper around a fs.Sub(module, "view")
func New(log log.Log, module *gomod.Module, tr transpiler.Transpiler) *Viewer {
	return &Viewer{log, module, tr}
}

type Viewer struct {
	log    log.Log
	module *gomod.Module
	tr     transpiler.Transpiler
}

var _ view.Viewer = (*Viewer)(nil)

func (v *Viewer) Register(r *router.Router, pages []*view.Page) {

}

func (v *Viewer) Mount(r *router.Router) error {
	return nil
}

func (v *Viewer) parseTemplate(templatePath string) (*renderer, error) {
	viewPath := path.Join("view", templatePath)
	code, err := fs.ReadFile(v.module, viewPath)
	if err != nil {
		return nil, err
	}
	code, err = v.tr.Transpile(viewPath, ".gohtml", code)
	if err != nil {
		return nil, fmt.Errorf("gohtml unable to transpile %s: %w", viewPath, err)
	}
	tpl, err := template.New(templatePath).Parse(string(code))
	if err != nil {
		return nil, err
	}
	return &renderer{tpl}, nil
}

func (v *Viewer) render(ctx context.Context, templatePath string, props interface{}) ([]byte, error) {
	entry, err := v.parseTemplate(templatePath)
	if err != nil {
		return nil, err
	}
	return entry.Render(ctx, props)
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
	html, err := v.render(ctx, page.Path, propMap[page.Key])
	if err != nil {
		return nil, err
	}
	for _, frame := range page.Frames {
		// TODO: support other props
		html, err = v.render(ctx, frame.Path, string(html))
		if err != nil {
			return nil, err
		}
	}
	// TODO: support other props
	return v.render(ctx, page.Layout.Path, string(html))
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
