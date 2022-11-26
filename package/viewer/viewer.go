package viewer

import (
	"context"

	"github.com/livebud/bud/package/router"
)

type Viewer interface {
	Render(ctx context.Context, page *Page) ([]byte, error)
	RenderError(ctx context.Context, page *Page) []byte
	RegisterClient(r *router.Router, page *Page)
}

type Page struct {
	*View
	Layout *View   `json:"layout,omitempty"`
	Frames []*View `json:"frames,omitempty"`
	Client string  `json:"client,omitempty"`
}

type Props = map[string]interface{}
type Context = map[string]interface{}

type View struct {
	Key     string  `json:"key,omitempty"`     // Can be blank
	Path    string  `json:"path,omitempty"`    // Can be blank
	Props   Props   `json:"props,omitempty"`   // Can be nil
	Context Context `json:"context,omitempty"` // Can be nil
}
