package web

import (
	"net/http"

	"github.com/livebud/bud/example/zero/generator/web"
)

type Web struct {
	// Controller dependencies
	// Middleware dependencies
	// View dependencies
}

func (w *Web) Router(r web.Router) web.Router {
	return r
}

func (w *Web) Middleware(m web.Middlewares) web.Middlewares {
	return m
}

func (w *Web) Handler(h http.Handler) http.Handler {
	return h
}

func (w *Web) Server(s *http.Server) *http.Server {
	return s
}
