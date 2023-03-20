package view

import (
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/gotemplate"
	"github.com/livebud/bud/package/imports"
	"github.com/livebud/bud/runtime/generator"
)

func New(module *gomod.Module) *Generator {
	return &Generator{module}
}

type Generator struct {
	module *gomod.Module
}

func (g *Generator) Extend(gen generator.FileSystem) {
	gen.GenerateFile("bud/pkg/web/view/view.go", g.generateFile)
}

const template = `package view

{{- if $.Imports }}

import (
	{{- range $import := $.Imports }}
	{{ $import.Name }} "{{ $import.Path }}"
	{{- end }}
)
{{- end }}

func New(
	gohtml *gohtml.Viewer,
) *View {
	viewers := map[string]view.Viewer{
		"gohtml": gohtml,
	}
	pages := map[string]*view.Page{
		"posts/index": &view.Page{
			View: &view.View{
				Key: "posts/index",
				Path: "posts/index.gohtml",
			},
			Frames: []*view.View{},
			Layout: &view.View{
				Key: "layout",
				Path: "layout.gohtml",
			},
			Error: &view.View{
				Key: "error",
				Path: "error.gohtml",
			},
		},
	}
	return &View{
		viewers,
		pages,
	}
}

type View struct {
	viewers map[string]view.Viewer
	pages  map[string]*view.Page
}

var _ view.Interface = (*View)(nil)

// TODO: use a router.Router interface
func (v *View) Mount(r *router.Router) error {
	return nil
}

func (v *View) Render(ctx context.Context, key string, propMap view.PropMap) ([]byte, error) {
	page, ok := v.pages[key]
	if !ok {
		return nil, fmt.Errorf("generator/view: no page for key %s", key)
	}
	// TODO: figure out how to decide on the viewer
	viewer, ok := v.viewers["gohtml"]
	if !ok {
		return nil, fmt.Errorf("generator/view: no viewer for key %s", key)
	}
	return viewer.Render(ctx, page, propMap)
}

func (v *View) RenderError(ctx context.Context, key string, propMap view.PropMap, err error) []byte {
	page, ok := v.pages[key]
	if !ok {
		return []byte(fmt.Sprintf("no page %q to render error: %s", key, err))
	}
	// TODO: figure out how to decide on the viewer
	viewer, ok := v.viewers["gohtml"]
	if !ok {
		return []byte(fmt.Sprintf("no viewer for %q to render error: %s", key, err))
	}
	return viewer.RenderError(ctx, page, propMap, err)
}
`

var gen = gotemplate.MustParse("view.gotext", template)

type State struct {
	Imports []*imports.Import
}

func (g *Generator) generateFile(fsys generator.FS, file *generator.File) error {
	imset := imports.New()
	imset.AddStd("context", "fmt")
	imset.AddNamed("view", "github.com/livebud/bud/runtime/view")
	imset.AddNamed("router", "github.com/livebud/bud/package/router")
	imset.AddNamed("gohtml", g.module.Import("viewer/gohtml"))
	// imset.AddNamed("posts", g.module.Import("controller/posts"))
	// imset.AddNamed("users", g.module.Import("controller/users"))
	// imset.AddNamed("sessions", g.module.Import("controller/sessions"))
	code, err := gen.Generate(&State{
		Imports: imset.List(),
	})
	if err != nil {
		return err
	}
	file.Data = code
	return nil
}
