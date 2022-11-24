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
	Layout *View   `json:"layout,omitempty"`
	Frames []*View `json:"frames,omitempty"`
	Main   *View   `json:"main,omitempty"`
}

func (p *Page) Client() string {
	return path.Join("bud", p.Main.Path+".js")
}

type Props = map[string]interface{}
type Context = map[string]interface{}

type View struct {
	Path    string  `json:"path,omitempty"`
	Props   Props   `json:"props,omitempty"`
	Context Context `json:"context,omitempty"`
}
