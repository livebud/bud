package budserver

import (
	"fmt"
	"net/http"
	"time"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/pubsub"
	"github.com/livebud/bud/package/hot"
	"github.com/livebud/bud/package/js"
	"github.com/livebud/bud/package/overlay"
	"github.com/livebud/bud/package/router"
)

func New(overlay *overlay.FileSystem, hot *hot.Server, vm js.VM, flag *framework.Flag) *Server {
	router := router.New()
	server := &Server{
		Handler:    router,
		vm:         vm,
		fileServer: http.FileServer(http.FS(overlay)),
	}
	router.Get("/bud/hot", hot)
	router.Post("/bud/view/_ssr.js", http.HandlerFunc(server.renderView))
	router.Get("/bud/view/:path*", http.HandlerFunc(server.serveFile))
	return server
}

type Server struct {
	http.Handler
	vm         js.VM
	fileServer http.Handler

	ps  pubsub.Client
	Now func() time.Time // Used for testing
}

func (s *Server) Reload(path string) {
	s.ps.Publish(path, nil)
}

func (s *Server) renderView(w http.ResponseWriter, r *http.Request) {
	fmt.Println("rendering view", r.URL.Path)
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (s *Server) serveFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("serving file", r.URL.Path)
	s.fileServer.ServeHTTP(w, r)
}
