package view

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"strings"
	"testing"

	"gitlab.com/mnm/bud/pkg/gen"
	"gitlab.com/mnm/bud/pkg/gomod"
	"gitlab.com/mnm/bud/pkg/js"
	"gitlab.com/mnm/bud/runtime/transform"
	"gitlab.com/mnm/bud/runtime/view/dom"
	"gitlab.com/mnm/bud/runtime/view/ssr"
)

type Response struct {
	Status  int               `json:"status,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    string            `json:"body,omitempty"`
}

func (res *Response) Write(w http.ResponseWriter) {
	// Write the response out
	for key, value := range res.Headers {
		w.Header().Set(key, value)
	}
	w.WriteHeader(res.Status)
	w.Write([]byte(res.Body))
}

// Renderer interface
type Renderer interface {
	Render(path string, props interface{}) (*Response, error)
}

func New(bf gen.FS, vm js.VM) *Server {
	return &Server{fs: bf, hfs: http.FS(bf), vm: vm}
}

func Test(t testing.TB) *Server {
	panic("view: Test not implemented yet")
}

// Live server serves view files on the fly. Used during development.
func Live(module *gomod.Module, genfs gen.FS, vm js.VM, transformer *transform.Transformer) *Server {
	dir := module.Directory()
	genfs.Add(map[string]gen.Generator{
		"bud/view":         dom.Runner(genfs, dir, transformer),
		"bud/node_modules": dom.NodeModules(dir),
		"bud/view/_ssr.js": ssr.Generator(genfs, dir, transformer),
	})
	return &Server{fs: genfs, hfs: http.FS(genfs), vm: vm}
}

// Static server serves the same files every time. Used during production.
func Static(genfs gen.FS, vm js.VM) *Server {
	return &Server{fs: genfs, hfs: http.FS(genfs), vm: vm}
}

type Server struct {
	fs  fs.FS
	hfs http.FileSystem
	vm  js.VM
}

// Map is a convenience function for the common case of passing a map of props
// into a view
type Map map[string]interface{}

// Respond is a convenience function for render
func (s *Server) Respond(w http.ResponseWriter, path string, props interface{}) {
	res, err := s.Render(path, props)
	if err != nil {
		// TODO: swap with logger
		fmt.Println("view: render error", err)
		http.Error(w, err.Error(), 500)
		return
	}
	headers := w.Header()
	for key, value := range res.Headers {
		headers.Set(key, value)
	}
	w.WriteHeader(res.Status)
	w.Write([]byte(res.Body))
}

func (s *Server) Render(path string, props interface{}) (*Response, error) {
	propBytes, err := json.Marshal(props)
	if err != nil {
		return nil, err
	}
	script, err := fs.ReadFile(s.fs, "bud/view/_ssr.js")
	if err != nil {
		return nil, err
	}
	// Evaluate the server
	expr := fmt.Sprintf(`%s; bud.render(%q, %s)`, script, path, propBytes)
	result, err := s.vm.Eval("_ssr.js", expr)
	if err != nil {
		return nil, err
	}
	// Unmarshal the response
	res := new(Response)
	if err := json.Unmarshal([]byte(result), res); err != nil {
		return nil, err
	}
	return res, nil
}

func isClient(path string) bool {
	return strings.HasPrefix(path, "/bud/node_modules/") ||
		strings.HasPrefix(path, "/bud/view/")
}

func (s *Server) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isClient(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}
		s.ServeHTTP(w, r)
	})
}

// Handler returns a handler for a specific server-side route
func (s *Server) Handler(route string, props interface{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.Respond(w, route, props)
	})
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	file, err := s.hfs.Open(r.URL.Path)
	if err != nil {
		// TODO: swap with logger
		fmt.Println("view: open error", err)
		http.Error(w, err.Error(), 500)
		return
	}
	stat, err := file.Stat()
	if err != nil {
		// TODO: swap with logger
		fmt.Println("view: stat error", err)
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Add("Content-Type", "text/javascript")
	http.ServeContent(w, r, r.URL.Path, stat.ModTime(), file)
}
