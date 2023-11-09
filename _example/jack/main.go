package main

import (
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"

	"github.com/livebud/bud/pkg/logs"
	"github.com/livebud/bud/pkg/mod"
	"github.com/livebud/bud/pkg/mux"
	"github.com/livebud/bud/pkg/slots"
	"github.com/livebud/bud/pkg/view"
	"github.com/livebud/bud/pkg/view/css"
	"github.com/livebud/bud/pkg/view/preact"
)

func main() {
	log := logs.Default()
	module := mod.MustFind()
	router := mux.New()
	preact := preact.New(module)
	css := css.New(module)
	router.Get("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slot := slots.New()
		if err := preact.Render(slot, "view/index.tsx", &view.Data{Slots: slot}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		slot.Close()
		slot = slot.Next()
		page, err := io.ReadAll(slot)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		head, err := io.ReadAll(slot.Slot("head"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := preact.RenderHTML(w, "view/layout.tsx", &view.Data{
			Props: map[string]interface{}{
				"page": string(page),
				"head": string(head),
			},
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		slot.Close()
	}))
	router.Get("/view/{path*}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		viewPath := strings.TrimPrefix(r.URL.Path, "/")
		switch path.Ext(r.URL.Path) {
		case ".js":
			viewPath = strings.TrimSuffix(viewPath, ".js")
			if err := preact.RenderJS(w, viewPath, &view.Data{}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		case ".css":
			w.Header().Set("Content-Type", "text/css")
			if err := css.Render(w, viewPath, &view.Data{}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		default:
			http.Error(w, fmt.Sprintf("%q not found", r.URL.Path), http.StatusNotFound)
		}
	}))
	router.Get("/{path*}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.FileServer(http.Dir("public")).ServeHTTP(w, r)
	}))
	log.Infof("Listening on http://localhost%s", ":8080")
	http.ListenAndServe(":8080", router)
}
