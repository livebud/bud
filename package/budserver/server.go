package budserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"strings"

	"github.com/livebud/bud/package/budclient"
	"github.com/livebud/bud/package/hot"
	"github.com/livebud/bud/package/log"

	"github.com/livebud/bud/internal/pubsub"

	"github.com/livebud/bud/package/js"
	"github.com/livebud/bud/package/router"
)

func New(fsys fs.FS, bus pubsub.Client, log log.Interface, vm js.VM) *Server {
	router := router.New()
	server := &Server{
		Handler: router,
		fsys:    fsys,
		hfs:     http.FS(fsys),
		log:     log,
		bus:     bus,
		vm:      vm,
	}
	// Routes that are proxied to from the browser through the app to bud
	router.Post("/bud/view/:route*", http.HandlerFunc(server.render))
	router.Get("/bud/view/:path*", http.HandlerFunc(server.serve))
	router.Get("/bud/node_modules/:path*", http.HandlerFunc(server.serve))
	// Routes that are directly requested by the browser to
	router.Get("/bud/hot/:page*", hot.New(log, bus))
	// Private routes between the app and bud
	router.Post("/bud/events", http.HandlerFunc(server.createEvent))
	return server
}

type Server struct {
	http.Handler
	fsys fs.FS
	hfs  http.FileSystem
	bus  pubsub.Publisher
	log  log.Interface
	vm   js.VM
}

var _ http.Handler = (*Server)(nil)

func (s *Server) render(w http.ResponseWriter, r *http.Request) {
	// Read the body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Load the props
	var props map[string]interface{}
	if err := json.Unmarshal(body, &props); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	script, err := fs.ReadFile(s.fsys, "bud/view/_ssr.js")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	route := "/" + r.URL.Query().Get("route")
	expr := fmt.Sprintf(`%s; bud.render(%q, %s)`, script, route, body)
	result, err := s.vm.Eval("_ssr.js", expr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(result))
}

func (s *Server) serve(w http.ResponseWriter, r *http.Request) {
	s.log.Debug("devserver: serving", "file", r.URL.Path)
	file, err := s.hfs.Open(r.URL.Path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			http.Error(w, err.Error(), 404)
			return
		}
		http.Error(w, err.Error(), 500)
		return
	}
	stat, err := file.Stat()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	// Maintain support to resolve and run "/bud/node_modules/livebud/runtime".
	if strings.HasPrefix(r.URL.Path, "/bud/node_modules/") ||
		strings.HasSuffix(r.URL.Path, ".svelte") {
		w.Header().Set("Content-Type", "application/javascript")
	}
	http.ServeContent(w, r, r.URL.Path, stat.ModTime(), file)
	s.log.Debug("devserver: served", "file", r.URL.Path)
}

func (s *Server) createEvent(w http.ResponseWriter, r *http.Request) {
	// Read the body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Unmarshal the request body into an event
	var event budclient.Event
	if err := json.Unmarshal(body, &event); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Publish the event
	s.bus.Publish(event.Topic, event.Data)
	// Return a No Content response
	w.WriteHeader(http.StatusNoContent)
}
