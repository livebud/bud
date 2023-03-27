package web

import (
	"net/http"

	"github.com/livebud/bud/example/zero/bud/pkg/web/controller"
	"github.com/livebud/bud/example/zero/bud/pkg/web/middleware"
	"github.com/livebud/bud/example/zero/bud/pkg/web/view"
	"github.com/livebud/bud/example/zero/mw"
	"github.com/livebud/bud/package/router"
)

type Web struct {
	Controller *controller.Controller
	Middleware *middleware.Middleware
	MW         *mw.Middleware
	View       *view.View
}

// type Stack []func(next http.Handler) http.Handler

// func (w *Web) Stack(s middleware.Stack) Stack {
// 	return Stack{
// 		w.MW.WrapRW,
// 		w.MW.Session,
// 		w.MW.CSRF,
// 	}
// }

func (w *Web) Stack(s middleware.Stack) middleware.Interface {
	return append(s, w.Middleware)
}

// TODO: use a router interface
func (w *Web) Router(r *router.Router) *router.Router {
	// r.Get("/posts", w.Controller.Index)
	r.Mount(w.Controller)
	r.Mount(w.View.Posts.Intro)
	return r
}

func (w *Web) Handler(h http.Handler) http.Handler {
	return h
}

func (w *Web) Server(s *http.Server) *http.Server {
	return s
}
