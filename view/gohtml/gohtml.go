package gohtml

import (
	"context"
	"fmt"
	"html/template"
	"io"

	"github.com/livebud/bud/view"
)

func New() view.Renderer {
	return &renderer{}
}

type renderer struct {
}

func (r *renderer) Render(ctx context.Context, s view.Slot, file view.File, props any) error {
	tpl, err := r.parseTemplate(ctx, s, file)
	if err != nil {
		return err
	}
	return tpl.Execute(s, props)
}

func (r *renderer) parseTemplate(ctx context.Context, slot view.Slot, file view.File) (*template.Template, error) {
	code, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("gohtml: unable to parse template %q. %w", file.Path(), err)
	}
	contextData := view.GetContext(ctx)
	fns := template.FuncMap{
		"slot": func() (template.HTML, error) {
			html, err := io.ReadAll(slot)
			if err != nil {
				return "", err
			}
			return template.HTML(html), nil
		},
		"context": func() view.Context {
			return contextData
		},
		"csrf": func() (string, error) {
			csrf, ok := contextData["csrf"].(string)
			if !ok {
				return "", fmt.Errorf("gohtml: csrf not found in context")
			}
			return csrf, nil
		},
	}
	tpl, err := template.New(file.Path()).Funcs(fns).Parse(string(code))
	if err != nil {
		return nil, err
	}
	return tpl, nil
}
