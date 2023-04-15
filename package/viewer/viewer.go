package viewer

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

// FS can be used to either use the real filesystem or an embedded one,
// depending on how Bud was built.
type FS = fs.FS

// Interface for viewers
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
	return "/view/" + path.Clean(v.Path) + ".js"
}

type Page struct {
	*View  // Entry
	Frames []*View
	Layout *View
	Error  *View
}

// Client is the standard entry route for pages that need a client
func (p *Page) Client() string {
	return "/view/" + path.Clean(p.View.Path) + ".entry.js"
}

type Embed = virtual.File
type Embeds = map[string]*Embed
type Pages map[Key]*Page

type Viewer interface {
	Mount(r *router.Router) error
	Render(ctx context.Context, key string, propMap PropMap) ([]byte, error)
	RenderError(ctx context.Context, key string, propMap PropMap, err error) []byte
	Bundle(ctx context.Context, embed virtual.Tree) error
}
