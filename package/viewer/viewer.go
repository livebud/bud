package viewer

import (
	"context"
	"path"

	"github.com/livebud/bud/package/router"
)

type Transpiler interface {
	Transpile(name string, src []byte) ([]byte, error)
}

type Viewer interface {
	Render(ctx context.Context, page *Page) ([]byte, error)
	RenderError(ctx context.Context, page *Page) []byte
	RegisterPage(r *router.Router, page *Page)
}

type Page struct {
	*View
	Layout *View   `json:"layout,omitempty"`
	Frames []*View `json:"frames,omitempty"`
}

func (p *Page) Client() string {
	return path.Join("bud", p.Path+".js")
}

type Props = map[string]interface{}
type Context = map[string]interface{}

type View struct {
	Key     string  `json:"key,omitempty"`
	Path    string  `json:"path,omitempty"` // Can be empty
	Props   Props   `json:"props,omitempty"`
	Context Context `json:"context,omitempty"`
}
