package viewer

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/livebud/bud/package/virtual"
)

var ErrViewerNotFound = errors.New("viewer not found")
var ErrPageNotFound = errors.New("page not found")

type Router interface {
	Get(path string, handler http.Handler) error
}

type Props = map[string]interface{}

type Key = string

type Page struct {
	*View  // Entry
	Frames []*View
	Layout *View
	Error  *View
}

type View struct {
	Key  Key
	Path string
}

type Viewer interface {
	Register(router Router)
	Render(ctx context.Context, key string, props Props) ([]byte, error)
	Render2(ctx context.Context, key string, props Props) ([]byte, error)
	RenderError(ctx context.Context, key string, err error, props Props) []byte
	Bundle(ctx context.Context, fsys virtual.FS) error
}

// Viewers group multiple viewers into one viewer.
//
// Note: viewers get mapped by key (e.g. posts/index), not by
// extension (e.g. .gohtml)
type Viewers map[Key]Viewer

var _ Viewer = (Viewers)(nil)

func (viewers Viewers) Register(router Router) {
	for _, viewer := range viewers {
		viewer.Register(router)
	}
}

func (viewers Viewers) Render(ctx context.Context, key string, props Props) ([]byte, error) {
	viewer, ok := viewers[key]
	if !ok {
		return nil, fmt.Errorf("%w %q", ErrViewerNotFound, key)
	}
	return viewer.Render(ctx, key, props)
}

func (viewers Viewers) Render2(ctx context.Context, key string, props Props) ([]byte, error) {
	viewer, ok := viewers[key]
	if !ok {
		return nil, fmt.Errorf("%w %q", ErrViewerNotFound, key)
	}
	return viewer.Render2(ctx, key, props)
}

func (viewers Viewers) RenderError(ctx context.Context, key string, err error, props Props) []byte {
	viewer, ok := viewers[key]
	if !ok {
		return []byte(fmt.Sprintf("%q %v. %v\n", key, ErrViewerNotFound, err))
	}
	return viewer.RenderError(ctx, key, err, props)
}

func (viewers Viewers) Bundle(ctx context.Context, fsys virtual.FS) error {
	for _, viewer := range viewers {
		if err := viewer.Bundle(ctx, fsys); err != nil {
			return err
		}
	}
	return nil
}
