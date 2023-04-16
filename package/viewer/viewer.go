package viewer

import (
	"context"
	"errors"
	"io/fs"
	"net/http"

	"github.com/livebud/bud/framework/controller/controllerrt/request"

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
	Key    Key
	Path   string
	Ext    string
	Client string // View client
}

type Page struct {
	*View  // Entry
	Frames []*View
	Layout *View
	Error  *View
	Route  string
	Client string // Entry client
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

// StaticPropMap returns a prop map for static views based on the request data.
func StaticPropMap(page *Page, r *http.Request) (PropMap, error) {
	props := map[string]interface{}{}
	if err := request.Unmarshal(r, &props); err != nil {
		return nil, err
	}
	propMap := PropMap{}
	propMap[page.Key] = props
	if page.Layout != nil {
		propMap[page.Layout.Key] = props
	}
	for _, frame := range page.Frames {
		propMap[frame.Key] = props
	}
	if page.Error != nil {
		propMap[page.Error.Key] = props
	}
	return propMap, nil
}
