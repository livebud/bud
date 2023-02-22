package view

import (
	"context"
	"errors"
	"io/fs"
	"net/http"
	"path"

	"github.com/livebud/bud/package/router"
)

var ErrViewerNotFound = errors.New("viewer not found")
var ErrPageNotFound = errors.New("page not found")

type Key = string
type Ext = string
type PropMap = map[Key]interface{}
type ReadFS = fs.FS

type WriteFS interface {
	MkdirAll(path string, perm fs.FileMode) error
	WriteFile(name string, data []byte, perm fs.FileMode) error
}

// Interface for bud/internal/web/view
type Interface interface {
	Register(r *router.Router)
	Render(ctx context.Context, key string, propMap PropMap) ([]byte, error)
	RenderError(ctx context.Context, key string, propMap PropMap, err error) []byte
}

type View struct {
	Key  Key
	Path string
}

// Client is the standard route for specific views. This is typically used for
// hot reloading individual views.
func (v *View) Client() string {
	return "/bud/" + path.Clean(v.Path) + ".js"
}

type Page struct {
	*View  // Entry
	Frames []*View
	Layout *View
	Error  *View
	Viewer Ext // Chosen viewer for page
}

// Client is the standard entry route for pages that need a client
func (p *Page) Client() string {
	return "/bud/" + path.Clean(p.View.Path) + ".entry.js"
}

type Viewer interface {
	Register(r *router.Router, pages []*Page)
	Render(ctx context.Context, page *Page, propMap PropMap) ([]byte, error)
	RenderError(ctx context.Context, page *Page, propMap PropMap, err error) []byte
	Bundle(ctx context.Context, fsys WriteFS, pages []*Page) error
}

// Renderer is a temporary function we're using until we refactor the generated
// controllers
func Renderer(view Interface, key Key, propMap PropMap) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		html, err := view.Render(r.Context(), key, propMap)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(view.RenderError(r.Context(), key, propMap, err))
			return
		}
		w.Write(html)
	})
}

// func NewPages() Pages {
// 	return &pageSet{
// 		pages: []*Page{},
// 		index: map[Key]int{},
// 	}
// }

// type Pages interface {
// 	Add(pages ...*Page)
// 	List() []*Page
// 	Get(key Key) (*Page, bool)
// }

// type pageSet struct {
// 	pages []*Page
// 	index map[Key]int
// }

// func (p *pageSet) Add(pages ...*Page) {
// 	for _, page := range pages {
// 		if _, ok := p.index[page.Key]; ok {
// 			continue
// 		}
// 		p.index[page.Key] = len(p.pages)
// 		p.pages = append(p.pages, page)
// 	}
// }

// func (p *pageSet) List() []*Page {
// 	return p.pages
// }

// func (p *pageSet) Get(key Key) (*Page, bool) {
// 	if i, ok := p.index[key]; ok {
// 		return p.pages[i], true
// 	}
// 	return nil, false
// }

type Viewers map[Ext]Viewer

// func (v Viewers) Register(r *router.Router, pages []*Page) {
// 	// for _, viewer := range v {
// 	// 	viewer.Register(r, pages)
// 	// }
// }

type Pages map[Key]*Page

type ViewerPages map[Ext][]*Page
type PageViewer map[Key]Ext
