package web

import (
	"net/http"

	"github.com/livebud/bud/package/middleware"

	"github.com/livebud/bud/example/zero/bud/pkg/web/controller"
	"github.com/livebud/bud/example/zero/bud/pkg/web/view"
	"github.com/livebud/bud/package/router"
)

type Web struct {
	Controller *controller.Controller
	View       *view.View
}

// TODO: use a router interface
func (w *Web) Router(r *router.Router) *router.Router {
	// r.Get("/posts", w.Controller.Index)
	r.Mount(w.Controller)
	r.Mount(w.View.Posts.Intro)
	return r
}

func (w *Web) Middleware(s middleware.Stack) middleware.Stack {
	return s
}

func (w *Web) Handler(h http.Handler) http.Handler {
	return h
}

func (w *Web) Server(s *http.Server) *http.Server {
	return s
}
