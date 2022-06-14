package devserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"strings"

	"github.com/livebud/bud/package/hot"

	"github.com/livebud/bud/internal/pubsub"

	"github.com/livebud/bud/package/js"
	"github.com/livebud/bud/package/router"
)

func New(fsys fs.FS, ps pubsub.Subscriber, vm js.VM) *Server {
	router := router.New()
	server := &Server{
		Handler: router,
		fsys:    fsys,
		hfs:     http.FS(fsys),
		ps:      ps,
		vm:      vm,
	}
	router.Post("/bud/view/:route*", http.HandlerFunc(server.render))
	router.Get("/bud/view/:path*", http.HandlerFunc(server.serve))
	router.Get("/bud/node_modules/:path*", http.HandlerFunc(server.serve))
	// TODO: pass the pubsub subscriber to the hot handler
	router.Get("/bud/hot", hot.New())
	return server
}

type Server struct {
	http.Handler
	fsys fs.FS
	hfs  http.FileSystem
	ps   pubsub.Subscriber
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
}
