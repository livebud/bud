package budsvr

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"

	"github.com/livebud/bud/package/virtual"

	"github.com/livebud/bud/package/budhttp"
	"github.com/livebud/bud/package/hot"
	"github.com/livebud/bud/package/log"

	"github.com/livebud/bud/internal/pubsub"

	"github.com/livebud/bud/package/js"
	"github.com/livebud/bud/package/router"
)

func New(fsys fs.FS, bus pubsub.Client, log log.Log, vm js.VM) *Server {
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
	router.Get("/open/:path*", http.HandlerFunc(server.open))
	// Routes that are directly requested by the browser to
	router.Get("/bud/hot/:page*", hot.New(log, bus))
	// Private routes between the app and bud
	router.Post("/bud/events", http.HandlerFunc(server.publish))
	// Support eval
	router.Post("/js/script", http.HandlerFunc(server.script))
	router.Post("/js/eval", http.HandlerFunc(server.eval))
	return server
}

type Server struct {
	http.Handler
	fsys fs.FS
	hfs  http.FileSystem
	bus  pubsub.Publisher
	log  log.Log
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
		if errors.Is(err, fs.ErrNotExist) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
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

func (s *Server) open(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	s.log.Debug("devserver: opening", "file", path)
	file, err := s.fsys.Open(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			http.Error(w, err.Error(), 404)
			return
		}
		http.Error(w, err.Error(), 500)
		return
	}
	body, err := virtual.MarshalJSON(file)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
	s.log.Debug("devserver: opened", "file", path)
}

func (s *Server) publish(w http.ResponseWriter, r *http.Request) {
	// Read the body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Unmarshal the request body into an event
	var event budhttp.Event
	if err := json.Unmarshal(body, &event); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Publish the event
	s.bus.Publish(event.Topic, event.Data)
	// Return a No Content response
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) script(w http.ResponseWriter, r *http.Request) {
	// Read the body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Unmarshal the request body into an event
	var script budhttp.Script
	if err := json.Unmarshal(body, &script); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := s.vm.Script(script.Path, script.Script); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Return a No Content response
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) eval(w http.ResponseWriter, r *http.Request) {
	// Read the body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Unmarshal the request body into an event
	var eval budhttp.Eval
	if err := json.Unmarshal(body, &eval); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result, err := s.vm.Eval(eval.Path, eval.Expr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Return the result
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
}
