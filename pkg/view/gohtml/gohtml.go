package gohtml

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"strings"

	"github.com/livebud/bud/pkg/view"
)

type HTML = template.HTML
type FuncMap = template.FuncMap

var defaultFuncs = map[string]func(data *view.Data) any{
	"slot": func(data *view.Data) any {
		return func(name ...string) (template.HTML, error) {
			if len(name) == 0 {
				slot, err := io.ReadAll(data.Slots)
				if err != nil {
					return "", err
				}
				return template.HTML(slot), nil
			}
			sb := new(strings.Builder)
			for _, name := range name {
				slot, err := io.ReadAll(data.Slots.Slot(name))
				if err != nil {
					return "", err
				}
				sb.Write(slot)
			}
			return template.HTML(sb.String()), nil
		}
	},
	"attr": func(data *view.Data) any {
		return func(key string) (any, error) {
			return data.Attrs[key], nil
		}
	},
	"head": func(data *view.Data) any {
		return func(elements ...string) (template.HTML, error) {
			if len(elements) == 0 {
				head, err := io.ReadAll(data.Slots.Slot("head"))
				if err != nil {
					return "", err
				}
				return template.HTML(head), nil
			}
			for _, element := range elements {
				pipe := data.Slots.Slot("head")
				pipe.Write([]byte(element))
			}
			return "", nil
		}
	},
}

func New(fsys fs.FS) *Viewer {
	return &Viewer{fsys, defaultFuncs}
}

var _ view.Viewer = (*Viewer)(nil)

type Viewer struct {
	fsys  fs.FS
	funcs map[string]func(data *view.Data) any
}

func (v *Viewer) Func(name string, fn func(data *view.Data) any) {
	v.funcs[name] = fn
}

func (v *Viewer) Render(w io.Writer, path string, data *view.Data) error {
	tpl, err := v.parseTemplate(path, data)
	if err != nil {
		return err
	}
	return tpl.Execute(w, data.Props)
}

// func (v *Viewer) Middleware() http.Handler {
// 	return view.Middleware(v).Middleware(next)
// }

func (v *Viewer) parseTemplate(path string, data *view.Data) (*template.Template, error) {
	code, err := fs.ReadFile(v.fsys, path)
	if err != nil {
		return nil, fmt.Errorf("gohtml: unable to parse template %q. %w", path, err)
	}
	// Load all the template functions
	funcMap := template.FuncMap{}
	for name, loader := range v.funcs {
		funcMap[name] = loader(data)
	}
	// Parse the template
	tpl, err := template.New(path).Funcs(funcMap).Parse(string(code))
	if err != nil {
		return nil, err
	}
	return tpl, nil
}
