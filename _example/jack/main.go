package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/livebud/bud/pkg/logs"
	"github.com/livebud/bud/pkg/mod"
	"github.com/livebud/bud/pkg/mux"
	"github.com/livebud/bud/pkg/request"
	"github.com/livebud/bud/pkg/slots"
	"github.com/livebud/bud/pkg/sse"
	"github.com/livebud/bud/pkg/view"
	"github.com/livebud/bud/pkg/view/css"
	"github.com/livebud/bud/pkg/view/preact"
	"github.com/livebud/bud/pkg/watcher"
	"golang.org/x/sync/errgroup"
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
	se := sse.New(log)
	ctx := context.Background()
	router.Get("/live.js", se)
	preact := preact.New(module, preact.WithLive("/live.js"), preact.WithEnv(map[string]any{
		"API_URL":            os.Getenv("API_URL"),
		"SLACK_CLIENT_ID":    os.Getenv("SLACK_CLIENT_ID"),
		"SLACK_REDIRECT_URL": os.Getenv("SLACK_REDIRECT_URL"),
		"SLACK_SCOPE":        os.Getenv("SLACK_SCOPE"),
		"SLACK_USER_SCOPE":   os.Getenv("SLACK_USER_SCOPE"),
		"STRIPE_CLIENT_KEY":  os.Getenv("STRIPE_CLIENT_KEY"),
	}))
	css := css.New(module)
	router.Layout("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slot, err := slots.FromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		page, err := io.ReadAll(slot)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		heads, err := io.ReadAll(slot.Slot("heads"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := preact.RenderHTML(w, "view/layout.tsx", &view.Data{
			Props: map[string]interface{}{
				"page":  string(page),
				"heads": json.RawMessage(heads),
			},
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}))
	router.Get("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var props struct {
			Success *bool `json:"success"`
		}
		if err := request.Unmarshal(r, &props); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		slot, err := slots.FromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := preact.Render(slot, "view/index.tsx", &view.Data{Props: props, Slots: slot}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}))
	router.Get("/faq", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var props struct {
			Success *bool `json:"success"`
		}
		if err := request.Unmarshal(r, &props); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		slot, err := slots.FromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := preact.Render(slot, "view/faq.tsx", &view.Data{Props: props, Slots: slot}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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
	eg := new(errgroup.Group)
	eg.Go(func() error {
		log.Infof("Listening on http://localhost%s", ":8080")
		return http.ListenAndServe(":8080", router)
	})
	eg.Go(func() error {
		return watcher.Watch(ctx, ".", func(events *watcher.Events) error {
			eventData, err := json.Marshal(events)
			if err != nil {
				log.Error(err)
			}
			if err := se.Publish(ctx, &sse.Event{Data: []byte(eventData)}); err != nil {
				log.Error(err)
			}
			return nil
		})
	})
	if err := eg.Wait(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
