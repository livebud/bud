package view

import (
	"context"
	"errors"
	"io/fs"
	"path"

	"github.com/livebud/bud/package/router"
	"github.com/livebud/bud/package/virtual"
)

var ErrViewerNotFound = errors.New("viewer not found")
var ErrPageNotFound = errors.New("page not found")

type Key = string
type Ext = string
type PropMap = map[Key]interface{}
type ReadFS = fs.FS

// Interface for bud/internal/web/view
type Interface interface {
	Mount(r *router.Router) error
	Render(ctx context.Context, key string, propMap PropMap) ([]byte, error)
	RenderError(ctx context.Context, key string, propMap PropMap, err error) []byte
}

type View struct {
	Key  Key
	Path string
	Ext  string
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
}

// Client is the standard entry route for pages that need a client
func (p *Page) Client() string {
	return "/bud/" + path.Clean(p.View.Path) + ".entry.js"
}

type Embed = virtual.File
type Embeds = map[string]*Embed

type Viewer interface {
	Register(r *router.Router, pages []*Page)
	Render(ctx context.Context, fsys fs.FS, page *Page, propMap PropMap) ([]byte, error)
	RenderError(ctx context.Context, fsys fs.FS, page *Page, propMap PropMap, err error) []byte
	Bundle(ctx context.Context, fsys fs.FS, pages Pages, embed Embeds) error
}

type Viewers map[Ext]Viewer
type Pages map[Key]*Page
type ViewerPages map[Ext][]*Page
type PageViewer map[Key]Ext
