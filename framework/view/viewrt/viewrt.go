package viewrt

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/livebud/bud/framework/view/ssr"
	"github.com/livebud/bud/package/js"
	"github.com/livebud/bud/package/log"
)

type FS = fs.FS

func New(fsys FS, log log.Log, vm js.VM) *Handler {
	return &Handler{http.FS(fsys), fsys, log, vm}
}

type Handler struct {
	hfs  http.FileSystem
	fsys FS
	log  log.Log
	vm   js.VM
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	file, err := h.hfs.Open(r.URL.Path)
	if err != nil {
		h.log.Field("error", err).Error("view: open error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	stat, err := file.Stat()
	if err != nil {
		h.log.Field("error", err).Error("view: stat error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Always add application/javascript since we're now directly targeting routes
	w.Header().Add("Content-Type", "application/javascript")
	http.ServeContent(w, r, r.URL.Path, stat.ModTime(), file)
}

func (h *Handler) Renderer(route string, props interface{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res, err := h.render(route, props)
		if err != nil {
			h.log.Field("error", err).Error("view: render error")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		headers := w.Header()
		for key, value := range res.Headers {
			headers.Set(key, value)
		}
		w.WriteHeader(res.Status)
		w.Write([]byte(res.Body))
	})
}

func (h *Handler) render(path string, props interface{}) (*ssr.Response, error) {
	propBytes, err := json.Marshal(props)
	if err != nil {
		return nil, err
	}
	script, err := fs.ReadFile(h.fsys, "bud/view/_ssr.js")
	if err != nil {
		return nil, err
	}
	// Evaluate the server
	expr := fmt.Sprintf(`%s; bud.render(%q, %s)`, script, path, propBytes)
	result, err := h.vm.Eval("_ssr.js", expr)
	if err != nil {
		return nil, err
	}
	// Unmarshal the response
	res := new(ssr.Response)
	if err := json.Unmarshal([]byte(result), res); err != nil {
		return nil, err
	}
	if res.Status < 100 || res.Status > 999 {
		return nil, fmt.Errorf("view: invalid status code %d", res.Status)
	}
	return res, nil
}
