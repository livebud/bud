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
	"github.com/livebud/bud/internal/virtual"

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
	// Open a file
	router.Get("/open/:path*", http.HandlerFunc(server.openFile))
	// Eval some JS
	router.Post("/eval", http.HandlerFunc(server.evalJS))
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

func (s *Server) openFile(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("path")
	s.log.Debug("devserver: opening", "file", filePath)
	file, err := s.hfs.Open(filePath)
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
	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	res, err := json.Marshal(virtual.File{
		Name:    filePath,
		Data:    data,
		Mode:    stat.Mode(),
		ModTime: stat.ModTime(),
		Sys:     stat.Sys(),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/javascript")
	w.Write(res)
	s.log.Debug("devserver: opened", "file", filePath)
}

func (s *Server) evalJS(w http.ResponseWriter, r *http.Request) {
	s.log.Debug("devserver: evaling")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var input budclient.JS
	if err := json.Unmarshal(body, &input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result, err := s.vm.Eval(input.Path, input.Expr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte(result))
	s.log.Debug("devserver: evaled")
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
