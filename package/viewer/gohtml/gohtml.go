package gohtml

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io/fs"

	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/router"
	"github.com/livebud/bud/package/viewer"
	"github.com/livebud/bud/package/virtual"
	"github.com/livebud/bud/runtime/transpiler"
)

func New(fsys viewer.FS, log log.Log, pages viewer.Pages, tr transpiler.Interface) *Viewer {
	return &Viewer{fsys, log, pages, tr}
}

type Viewer struct {
	fsys  viewer.FS
	log   log.Log
	pages viewer.Pages
	tr    transpiler.Interface
}

var _ viewer.Viewer = (*Viewer)(nil)

func (v *Viewer) Mount(r *router.Router) error {
	return nil
}

func (v *Viewer) parseTemplate(templatePath string) (*template.Template, error) {
	// TODO: decide if we want to scope to the view path or module path
	code, err := fs.ReadFile(v.fsys, templatePath)
	if err != nil {
		return nil, fmt.Errorf("gohtml: unable to parse template %q. %w", templatePath, err)
	}
	// TODO: don't transpile when embedded
	code, err = v.tr.Transpile(templatePath, ".gohtml", code)
	if err != nil {
		return nil, fmt.Errorf("gohtml: unable to transpile %s: %w", templatePath, err)
	}
	tpl, err := template.New(templatePath).Parse(string(code))
	if err != nil {
		return nil, err
	}
	return tpl, nil
}

func (v *Viewer) render(ctx context.Context, templatePath string, props interface{}) ([]byte, error) {
	tpl, err := v.parseTemplate(templatePath)
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

func (v *Viewer) Render(ctx context.Context, key string, propMap viewer.PropMap) ([]byte, error) {
	page, ok := v.pages[key]
	if !ok {
		return nil, fmt.Errorf("gohtml: unable to find page from key %q", key)
	}
	v.log.Info("gohtml: rendering", page.Path)
	html, err := v.render(ctx, page.Path, propMap[page.Key])
	if err != nil {
		return nil, err
	}
	for _, frame := range page.Frames {
		// TODO: support other props
		html, err = v.render(ctx, frame.Path, template.HTML(html))
		if err != nil {
			return nil, err
		}
	}
	if page.Layout != nil {
		// TODO: support other props
		html, err = v.render(ctx, page.Layout.Path, template.HTML(html))
		if err != nil {
			return nil, err
		}
	}
	return html, nil
}

func (v *Viewer) RenderError(ctx context.Context, key string, propMap viewer.PropMap, originalError error) []byte {
	page, ok := v.pages[key]
	if !ok {
		return []byte(fmt.Sprintf("gohtml: unable to find page from key %q to render error %s", key, originalError))
	}
	v.log.Info("gohtml: rendering error", page.Error.Path)
	errorEntry, err := v.parseTemplate(page.Error.Path)
	if err != nil {
		return []byte(fmt.Sprintf("gohtml: unable to read error template %q to render error %s. %s", page.Error.Path, err, originalError))
	}
	layout, err := v.parseTemplate(page.Layout.Path)
	if err != nil {
		return []byte(fmt.Sprintf("gohtml: unable to parse layout template %q to render error %s. %s", page.Error.Path, err, originalError))
	}
	state := errorState{
		Message: originalError.Error(),
	}
	html, err := render(ctx, errorEntry, state)
	if err != nil {
		return []byte(fmt.Sprintf("gohtml: unable to render error template %q to render error %s. %s", page.Error.Path, err, originalError))
	}
	html, err = render(ctx, layout, template.HTML(html))
	if err != nil {
		return []byte(fmt.Sprintf("gohtml: unable to render layout template %q to render error %s. %s", page.Error.Path, err, originalError))
	}
	return html
}

type errorState struct {
	Message string
}

func (v *Viewer) Bundle(ctx context.Context, fs virtual.Tree) (err error) {
	for _, page := range v.pages {
		// Embed the page
		pageEmbed, err := v.embedView(page.Path)
		if err != nil {
			return err
		}
		fs[page.Path] = pageEmbed

		// Embed the layout
		if page.Layout != nil {
			if _, ok := fs[page.Layout.Path]; ok {
				continue
			}
			layoutEmbed, err := v.embedView(page.Layout.Path)
			if err != nil {
				return err
			}
			fs[page.Layout.Path] = layoutEmbed
		}

		// Embed the frames
		for _, frame := range page.Frames {
			if _, ok := fs[frame.Path]; ok {
				continue
			}
			frameEmbed, err := v.embedView(frame.Path)
			if err != nil {
				return err
			}
			fs[frame.Path] = frameEmbed
		}

		// Embed the error
		if page.Error != nil {
			if _, ok := fs[page.Error.Path]; ok {
				continue
			}
			errorEmbed, err := v.embedView(page.Error.Path)
			if err != nil {
				return err
			}
			fs[page.Error.Path] = errorEmbed
		}
	}
	return nil
}

func (v *Viewer) embedView(viewPath string) (*viewer.Embed, error) {
	code, err := fs.ReadFile(v.fsys, viewPath)
	if err != nil {
		return nil, err
	}
	// TODO: bring back pre-transpilation
	// Sanity check that the transpiled code is valid
	if _, err := template.New(viewPath).Parse(string(code)); err != nil {
		return nil, fmt.Errorf("gohtml: unable to parse transpiled template %q. %w", viewPath, err)
	}
	return &viewer.Embed{
		Path: viewPath,
		Data: code,
	}, nil
}
