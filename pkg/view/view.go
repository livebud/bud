package view

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/livebud/bud/pkg/middleware"
	"github.com/livebud/bud/pkg/mux"
	"github.com/livebud/bud/pkg/slot"
)

var ErrNotInContext = errors.New("viewer: not in context")

// Data for a view
type Data struct {
	// TODO: consider splitting into slot.Slot and slot.Slots since this might
	// simplify the slot API (no more ReadString, etc.)
	Slots *slot.Slots
	Attrs map[string]any
	Props any
}

// Viewer renders views
type Viewer interface {
	// TODO: make this an optional interface
	mux.Routes
	Render(w io.Writer, path string, data *Data) error
}

// Middleware to load viewers into the context
func Middleware(viewer Viewer) middleware.Middleware {
	return middleware.Func(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(WithViewer(r.Context(), viewer)))
		})
	})
}

type contextKey string

const viewerKey contextKey = "view:viewer"

// WithViewer sets the viewer on the context
func WithViewer(parent context.Context, viewer Viewer) context.Context {
	return context.WithValue(parent, viewerKey, viewer)
}

// From gets the viewer from the context or returns an error
func From(ctx context.Context) (Viewer, error) {
	viewer, ok := ctx.Value(viewerKey).(Viewer)
	if !ok {
		return nil, ErrNotInContext
	}
	return viewer, nil
}
