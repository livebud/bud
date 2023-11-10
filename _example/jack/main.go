package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"

	"github.com/livebud/bud/pkg/logs"
	"github.com/livebud/bud/pkg/mod"
	"github.com/livebud/bud/pkg/mux"
	"github.com/livebud/bud/pkg/request"
	"github.com/livebud/bud/pkg/slots"
	"github.com/livebud/bud/pkg/view"
	"github.com/livebud/bud/pkg/view/css"
	"github.com/livebud/bud/pkg/view/preact"
)

type VNode struct {
	Name     string            `json:"name,omitempty"`
	Attrs    map[string]string `json:"attrs,omitempty"`
	Children any               `json:"children,omitempty"`
	Value    string            `json:"value,omitempty"`
}

func main() {
	log := logs.Default()
	module := mod.MustFind()
	router := mux.New()
	preact := preact.New(module)
	css := css.New(module)
	router.Get("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var props struct {
			Success *bool `json:"success"`
		}
		if err := request.Unmarshal(r, &props); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		slot := slots.New()
		if err := preact.Render(slot, "view/index.tsx", &view.Data{Props: props, Slots: slot}); err != nil {
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
		jsonProps, err := json.Marshal(props)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		slot.Slot("head").Write([]byte(fmt.Sprintf(`<script id="bud#props" type="text/template" defer>%s</script>`, string(jsonProps))))
		slot.Slot("head").Write([]byte(fmt.Sprintf(`<script src="/view/%s.js" defer></script>`, "index.jsx")))
		head, err := io.ReadAll(slot.Slot("head"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_ = head
		if err := preact.RenderHTML(w, "view/layout.tsx", &view.Data{
			Props: map[string]interface{}{
				"page": string(page),
				"head": []VNode{
					{
						Name:     "title",
						Children: []string{"Standup Jack"},
					},
					{
						Name: "script",
						Attrs: map[string]string{
							"id":    "bud#props",
							"type":  "text/template",
							"defer": "",
						},
						Children: []string{"{}"},
					},
					{
						Name: "script",
						Attrs: map[string]string{
							"src":   "/view/index.jsx.js",
							"type":  "application/javascript",
							"defer": "",
						},
					},
				},
			},
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		slot.Close()
	}))
	router.Get("/faq", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slot := slots.New()
		if err := preact.Render(slot, "view/faq.tsx", &view.Data{Slots: slot}); err != nil {
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
		_ = head
		if err := preact.RenderHTML(w, "view/layout.tsx", &view.Data{
			Props: map[string]interface{}{
				"page": string(page),
				"head": []VNode{
					{
						Name:     "title",
						Children: []string{"FAQ"},
					},
					{
						Name: "script",
						Attrs: map[string]string{
							"id":    "bud#props",
							"type":  "text/template",
							"defer": "",
						},
						Children: []string{"{}"},
					},
					{
						Name: "script",
						Attrs: map[string]string{
							"src":   "/view/faq.jsx.js",
							"type":  "application/javascript",
							"defer": "",
						},
					},
				},
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
	router.Get("/{path*}", http.FileServer(http.Dir("public")))
	log.Infof("Listening on http://localhost%s", ":8080")
	http.ListenAndServe(":8080", router)
}
