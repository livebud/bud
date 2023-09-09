package gohtml

import (
	"context"
	"fmt"
	"html/template"
	"io"

	"github.com/livebud/bud/pkg/view"
)

func New() view.Renderer {
	return &renderer{}
}

type renderer struct {
}

func (r *renderer) Render(ctx context.Context, s view.Slot, file view.File, data view.Data, props any) error {
	tpl, err := r.parseTemplate(ctx, file)
	if err != nil {
		return err
	}
	return tpl.Execute(s, &templateData{data, props, s})
}

type templateData struct {
	data  view.Data
	Props any
	slot  view.Slot
}

// Slot returns slot data, if there is any. Otherwise returns an empty string
func (d *templateData) Slot() (template.HTML, error) {
	html, err := io.ReadAll(d.slot)
	if err != nil {
		return "", err
	}
	return template.HTML(html), nil
}

func (d *templateData) CSRF() string {
	return d.data["csrf"].(string)
}

func (r *renderer) parseTemplate(ctx context.Context, file view.File) (*template.Template, error) {
	code, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("gohtml: unable to parse template %q. %w", file.Path(), err)
	}
	tpl, err := template.New(file.Path()).Parse(withProps(string(code)))
	if err != nil {
		return nil, err
	}
	return tpl, nil
}

func withProps(template string) string {
	return `{{ with $.Props }}` + template + `{{ else }}` + template + `{{ end }}`
}
