package view

import (
	context "context"
	fmt "fmt"
	viewer "github.com/livebud/bud/example/zero/bud/pkg/viewer"
	gomod "github.com/livebud/bud/package/gomod"
	router "github.com/livebud/bud/package/router"
	virtual "github.com/livebud/bud/package/virtual"
	view "github.com/livebud/bud/runtime/view"
	fs "io/fs"
	http "net/http"
)

func New(
	module *gomod.Module,
	viewer viewer.Viewer,
) *View {

	fsys := module

	// TODO: generate with Pages instead
	postsIndexPage := &PostsIndexPage{
		fsys,
		viewer,
		&view.Page{
			View: &view.View{
				Key:  "posts/index",
				Path: "posts/index.gohtml",
				Ext:  ".gohtml",
			},
			Frames: []*view.View{
				{
					Key:  "posts/frame",
					Path: "posts/frame.gohtml",
					Ext:  ".gohtml",
				},
			},
			Layout: &view.View{
				Key:  "layout",
				Path: "layout.gohtml",
				Ext:  ".gohtml",
			},
			Error: &view.View{
				Key:  "error",
				Path: "error.gohtml",
				Ext:  ".gohtml",
			},
		},
	}

	postsIntroPage := &PostsIntroPage{
		fsys,
		viewer,
		&view.Page{
			View: &view.View{
				Key:  "posts/intro",
				Path: "posts/intro.md",
				Ext:  ".md",
			},
			Frames: []*view.View{
				{
					Key:  "posts/frame",
					Path: "posts/frame.gohtml",
					Ext:  ".gohtml",
				},
			},
			Layout: &view.View{
				Key:  "layout",
				Path: "layout.gohtml",
				Ext:  ".gohtml",
			},
			Error: &view.View{
				Key:  "error",
				Path: "error.gohtml",
				Ext:  ".gohtml",
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

var embeddedFS = virtual.List{}

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
	fsys   fs.FS
	viewer view.Viewer
	page   *view.Page
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
	return p.viewer.Render(ctx, p.fsys, p.page, propMap)
}

func (p *PostsIndexPage) RenderError(ctx context.Context, propMap view.PropMap, err error) []byte {
	return p.viewer.RenderError(ctx, p.fsys, p.page, propMap, err)
}

type PostsIntroPage struct {
	fsys   fs.FS
	viewer view.Viewer
	page   *view.Page
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
	return p.viewer.Render(ctx, p.fsys, p.page, propMap)
}

func (p *PostsIntroPage) RenderError(ctx context.Context, propMap view.PropMap, err error) []byte {
	return p.viewer.RenderError(ctx, p.fsys, p.page, propMap, err)
}
