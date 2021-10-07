package view

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	"github.com/go-duo/bud/js"
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

// // FileServer serves the client-side files
// func FileServer(fs fs.FS) *Server {
// 	return &Server{http.FS(fs)}
// }

// type Server struct {
// 	hfs http.FileSystem
// }

// func (s *Server) Middleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if !isClient(r.URL.Path) {
// 			next.ServeHTTP(w, r)
// 			return
// 		}
// 		s.ServeHTTP(w, r)
// 	})
// }

// func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	file, err := s.hfs.Open(r.URL.Path)
// 	if err != nil {
// 		fmt.Println(err)
// 		http.Error(w, err.Error(), 500)
// 		return
// 	}
// 	stat, err := file.Stat()
// 	if err != nil {
// 		fmt.Println(err)
// 		http.Error(w, err.Error(), 500)
// 		return
// 	}
// 	w.Header().Add("Content-Type", "text/javascript")
// 	http.ServeContent(w, r, r.URL.Path, stat.ModTime(), file)
// }

// func Runner(rootDir string, df *dfs.DFS, vm js.VM) *View {
// 	dirfs := os.DirFS(rootDir)
// 	svelte := svelte.New(&svelte.Input{
// 		VM:  vm,
// 		Dev: true,
// 	})
// 	// Add to the existing DFS
// 	df.Add(map[string]dfs.Generator{
// 		"duo/view":         dom.Runner(svelte, rootDir),
// 		"duo/node_modules": dom.NodeModules(rootDir),
// 		"duo/view/_ssr.js": ssr.Generator(dirfs, svelte, rootDir),
// 	})
// 	return &View{df, http.FS(df), vm}
// }

// func Builder(ef *dfs.EFS, vm js.VM) *View {
// 	return &View{ef, http.FS(ef), vm}
// }

func NewServer(fs fs.FS, vm js.VM) *Server {
	return &Server{fs, http.FS(fs), vm}
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
		fmt.Println(err)
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
	script, err := fs.ReadFile(s.fs, "duo/view/_ssr.js")
	if err != nil {
		return nil, err
	}
	// Evaluate the server
	expr := fmt.Sprintf(`%s; duo.render(%q, %s)`, script, path, propBytes)
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
	return strings.HasPrefix(path, "/duo/node_modules/") ||
		strings.HasPrefix(path, "/duo/view/")
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

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	file, err := s.hfs.Open(r.URL.Path)
	if err != nil {
		// TODO: swap with logger
		fmt.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}
	stat, err := file.Stat()
	if err != nil {
		// TODO: swap with logger
		fmt.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Add("Content-Type", "text/javascript")
	http.ServeContent(w, r, r.URL.Path, stat.ModTime(), file)
}
