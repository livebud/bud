package router

import (
	"io/fs"

	"github.com/livebud/bud/_example/jack/controller"
	"github.com/livebud/bud/pkg/mux"
	"github.com/livebud/bud/pkg/sse"
	"github.com/livebud/bud/pkg/ssr"
	"github.com/livebud/bud/pkg/u"
)

func New(fsys fs.FS, sse *sse.Handler) *mux.Router {
	router := mux.New()
	router.Use(sse.Middleware())
	ssr := ssr.New(fsys, "/live.js")
	router.Add(controller.New(ssr))
	router.ServeFS("/view/{path*}", fsys)
	publicFS := u.Must(fs.Sub(fsys, "public"))
	router.ServeFS("/{path*}", publicFS)
	return router
}
