package view

import (
	"fmt"

	"github.com/livebud/bud/example/zero/bud/pkg/viewer"
	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/gotemplate"
	"github.com/livebud/bud/package/imports"
	"github.com/livebud/bud/runtime/generator"
)

func New(flag *framework.Flag, module *gomod.Module, viewer viewer.Map) *Generator {
	return &Generator{flag, module, viewer}
}

type Generator struct {
	flag   *framework.Flag
	module *gomod.Module
	viewer viewer.Map
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
	viewer viewer.Map,
) *View {
	postsIndexPage := &PostsIndexPage{
		viewer: viewer[".gohtml"],
		page: &view.Page{
			View: &view.View{
				Key: "posts/index",
				Path: "posts/index.gohtml",
			},
			Frames: []*view.View{
				{
					Key: "posts/frame",
					Path: "posts/frame.gohtml",
				},
			},
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

	postsIntroPage := &PostsIntroPage{
		viewer: viewer[".gohtml"],
		page: &view.Page{
			View: &view.View{
				Key: "posts/intro",
				Path: "posts/intro.md",
			},
			Frames: []*view.View{
				{
					Key: "posts/frame",
					Path: "posts/frame.gohtml",
				},
			},
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
		map[string]renderer{
			"posts/index": postsIndexPage,
			"posts/intro": postsIntroPage,
		},
		&PostsView{
			postsIndexPage,
			postsIntroPage,
		},
	}
}

type renderer interface {
	Render(ctx context.Context, propMap view.PropMap) ([]byte, error)
	RenderError(ctx context.Context, propMap view.PropMap, err error) []byte
}

type View struct {
	pages map[string]renderer
	Posts *PostsView
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
	return page.Render(ctx, propMap)
}

func (v *View) RenderError(ctx context.Context, key string, propMap view.PropMap, err error) []byte {
	page, ok := v.pages[key]
	if !ok {
		return []byte(fmt.Sprintf("no page %q to render error: %s", key, err))
	}
	return page.RenderError(ctx, propMap, err)
}

type PostsView struct {
	Index *PostsIndexPage
	Intro *PostsIntroPage
}

func (p *PostsView) Mount(r *router.Router) error {
	r.Mount(p.Index)
	r.Mount(p.Intro)
	return nil
}

type PostsIndexPage struct {
	page *view.Page
	viewer view.Viewer
}

func (p *PostsIndexPage) Mount(r *router.Router) error {
	r.Get("/posts", p)
	return nil
}

func (p *PostsIndexPage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	propMap := map[string]interface{}{
		// TODO: static props?
		p.page.Key: nil,
	}
	html, err := p.Render(r.Context(), propMap)
	if err != nil {
		html = p.RenderError(r.Context(), propMap, err)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")
	w.Write(html)
}

func (p *PostsIndexPage) Render(ctx context.Context, propMap view.PropMap) ([]byte, error) {
	return p.viewer.Render(ctx, p.page, propMap)
}

func (p *PostsIndexPage) RenderError(ctx context.Context, propMap view.PropMap, err error) []byte {
	return p.viewer.RenderError(ctx, p.page, propMap, err)
}

type PostsIntroPage struct {
	page *view.Page
	viewer view.Viewer
}

func (p *PostsIntroPage) Mount(r *router.Router) error {
	r.Get("/posts/intro", p)
	return nil
}

// TODO: consolidate
func (p *PostsIntroPage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	propMap := map[string]interface{}{
		// TODO: static props?
		p.page.Key: nil,
	}
	html, err := p.Render(r.Context(), propMap)
	if err != nil {
		html = p.RenderError(r.Context(), propMap, err)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")
	w.Write(html)
}

func (p *PostsIntroPage) Render(ctx context.Context, propMap view.PropMap) ([]byte, error) {
	return p.viewer.Render(ctx, p.page, propMap)
}

func (p *PostsIntroPage) RenderError(ctx context.Context, propMap view.PropMap, err error) []byte {
	return p.viewer.RenderError(ctx, p.page, propMap, err)
}
`

var gen = gotemplate.MustParse("view.gotext", template)

type State struct {
	Imports []*imports.Import
}

func (g *Generator) generateFile(fsys generator.FS, file *generator.File) error {
	fmt.Println("TODO bundle", g.flag.Embed, g.viewer)
	imset := imports.New()
	imset.AddStd("context", "fmt", "net/http")
	imset.AddNamed("view", "github.com/livebud/bud/runtime/view")
	imset.AddNamed("router", "github.com/livebud/bud/package/router")
	imset.AddNamed("viewer", g.module.Import("bud/pkg/viewer"))
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
