package svelte

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	"github.com/livebud/bud/framework/view/ssr"
	"github.com/livebud/bud/package/js"
	"github.com/livebud/bud/package/router"
)

func New(fsys fs.FS, vm js.VM) *Viewer {
	return &Viewer{fsys, http.FS(fsys), vm}
}

type Viewer struct {
	fsys fs.FS
	hfs  http.FileSystem
	vm   js.VM
}

func (v *Viewer) Handler(route string, props interface{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v.Render(w, route, props)
	})
}

func (v *Viewer) Render(w http.ResponseWriter, viewPath string, props interface{}) {
	script, err := fs.ReadFile(v.fsys, "bud/view/_ssr.js")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	propBytes, err := json.Marshal(props)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	result, err := v.vm.Eval("_ssr.js", fmt.Sprintf(`%s; bud.render(%q, %s)`, script, viewPath, propBytes))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Unmarshal the response
	res := new(ssr.Response)
	if err := json.Unmarshal([]byte(result), res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if res.Status < 100 || res.Status > 999 {
		err := fmt.Errorf("view: invalid status code %d", res.Status)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	headers := w.Header()
	for key, value := range res.Headers {
		headers.Set(key, value)
	}
	w.WriteHeader(res.Status)
	w.Write([]byte(res.Body))
}

func (v *Viewer) Serve(router *router.Router) {
	router.Get("/bud/view/:path*", http.HandlerFunc(v.serveClient))
	router.Get("/bud/node_modules/:path*", http.HandlerFunc(v.serveClient))
}

func (v *Viewer) serveClient(w http.ResponseWriter, r *http.Request) {
	file, err := v.hfs.Open(r.URL.Path)
	if err != nil {
		// TODO: swap with logger
		fmt.Println("view: open error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	stat, err := file.Stat()
	if err != nil {
		// TODO: swap with logger
		fmt.Println("view: stat error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Maintain support to resolve and run "/bud/node_modules/livebud/runtime".
	if strings.HasPrefix(r.URL.Path, "/bud/node_modules/") {
		w.Header().Add("Content-Type", "text/javascript")
	}
	http.ServeContent(w, r, r.URL.Path, stat.ModTime(), file)
}
