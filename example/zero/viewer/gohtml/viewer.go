package gohtml

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io/fs"
	"path"

	"github.com/livebud/bud/example/zero/bud/pkg/transpiler"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/router"
	"github.com/livebud/bud/runtime/view"
)

// TODO: may want a view.FS that is a wrapper around a fs.Sub(module, "view")
func New(log log.Log, tr transpiler.Transpiler) *Viewer {
	return &Viewer{log, tr}
}

type Viewer struct {
	log log.Log
	tr  transpiler.Transpiler
}

var _ view.Viewer = (*Viewer)(nil)

func (v *Viewer) Register(r *router.Router, pages []*view.Page) {

}

func (v *Viewer) Mount(r *router.Router) error {
	return nil
}

func (v *Viewer) parseTemplate(fsys fs.FS, templatePath string) (*template.Template, error) {
	// TODO: decide if we want to scope to the view path or module path
	viewPath := path.Join("view", templatePath)
	code, err := fs.ReadFile(fsys, viewPath)
	if err != nil {
		return nil, fmt.Errorf("unable to parse template %q. %w", templatePath, err)
	}
	// TODO: don't transpile when embedded
	code, err = v.tr.Transpile(viewPath, ".gohtml", code)
	if err != nil {
		return nil, fmt.Errorf("gohtml unable to transpile %s: %w", viewPath, err)
	}
	tpl, err := template.New(templatePath).Parse(string(code))
	if err != nil {
		return nil, err
	}
	return tpl, nil
}

func (v *Viewer) render(ctx context.Context, fsys fs.FS, templatePath string, props interface{}) ([]byte, error) {
	tpl, err := v.parseTemplate(fsys, templatePath)
	if err != nil {
		return nil, err
	}
	return render(ctx, tpl, props)
}

func render(ctx context.Context, tpl *template.Template, props interface{}) ([]byte, error) {
	out := new(bytes.Buffer)
	// TODO: pass context through
	if err := tpl.Execute(out, props); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func (v *Viewer) Render(ctx context.Context, fsys fs.FS, page *view.Page, propMap view.PropMap) ([]byte, error) {
	v.log.Info("rendering gohtml", page.Path)
	html, err := v.render(ctx, fsys, page.Path, propMap[page.Key])
	if err != nil {
		return nil, err
	}
	for _, frame := range page.Frames {
		// TODO: support other props
		html, err = v.render(ctx, fsys, frame.Path, template.HTML(html))
		if err != nil {
			return nil, err
		}
	}
	// TODO: support other props
	return v.render(ctx, fsys, page.Layout.Path, template.HTML(html))
}

func (v *Viewer) RenderError(ctx context.Context, fsys fs.FS, page *view.Page, propMap view.PropMap, originalError error) []byte {
	v.log.Info("rendering gohtml", page.Error.Path)
	errorEntry, err := v.parseTemplate(fsys, page.Error.Path)
	if err != nil {
		return []byte(fmt.Sprintf("unable to read error template %q to render error %s. %s", page.Error.Path, err, originalError))
	}
	layout, err := v.parseTemplate(fsys, page.Layout.Path)
	if err != nil {
		return []byte(fmt.Sprintf("unable to parse layout template %q to render error %s. %s", page.Error.Path, err, originalError))
	}
	state := errorState{
		Message: originalError.Error(),
	}
	html, err := render(ctx, errorEntry, state)
	if err != nil {
		return []byte(fmt.Sprintf("unable to render error template %q to render error %s. %s", page.Error.Path, err, originalError))
	}
	html, err = render(ctx, layout, template.HTML(html))
	if err != nil {
		return []byte(fmt.Sprintf("unable to render layout template %q to render error %s. %s", page.Error.Path, err, originalError))
	}
	return html
}

type errorState struct {
	Message string
}

func (v *Viewer) Bundle(ctx context.Context, fsys fs.FS, pages view.Pages, embeds view.Embeds) (err error) {
	for _, page := range pages {
		// TODO: decide if we want to scope to the view path or module path
		pagePath := path.Join("view", page.Path)
		// Embed the page
		pageEmbed, err := v.embedView(fsys, pagePath)
		if err != nil {
			return err
		}
		embeds[pagePath] = pageEmbed

		// Embed the layout
		if page.Layout != nil {
			layoutPath := path.Join("view", page.Layout.Path)
			if _, ok := embeds[layoutPath]; ok {
				continue
			}
			layoutEmbed, err := v.embedView(fsys, layoutPath)
			if err != nil {
				return err
			}
			embeds[layoutPath] = layoutEmbed
		}

		// Embed the frames
		for _, frame := range page.Frames {
			framePath := path.Join("view", frame.Path)
			if _, ok := embeds[framePath]; ok {
				continue
			}
			frameEmbed, err := v.embedView(fsys, framePath)
			if err != nil {
				return err
			}
			embeds[framePath] = frameEmbed
		}

		// Embed the error
		if page.Error != nil {
			errorPath := path.Join("view", page.Error.Path)
			if _, ok := embeds[errorPath]; ok {
				continue
			}
			errorEmbed, err := v.embedView(fsys, errorPath)
			if err != nil {
				return err
			}
			embeds[errorPath] = errorEmbed
		}
	}
	return nil
}

func (v *Viewer) embedView(fsys fs.FS, viewPath string) (*view.Embed, error) {
	code, err := fs.ReadFile(fsys, viewPath)
	if err != nil {
		return nil, err
	}
	// TODO: bring back pre-transpilation
	// Sanity check that the transpiled code is valid
	if _, err := template.New(viewPath).Parse(string(code)); err != nil {
		return nil, fmt.Errorf("gohtml: unable to parse transpiled template %q. %w", viewPath, err)
	}
	return &view.Embed{
		Path: viewPath,
		Data: code,
	}, nil
}
