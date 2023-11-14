package controller

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/livebud/bud/pkg/mux"
	"github.com/livebud/bud/pkg/request"
	"github.com/livebud/bud/pkg/slots"
	"github.com/livebud/bud/pkg/ssr"
	"github.com/livebud/bud/pkg/view"
)

func New(ssr *ssr.Viewer) *Controller {
	return &Controller{ssr}
}

type Controller struct {
	ssr *ssr.Viewer
}

func (c *Controller) Routes(r mux.Router) {
	r.Layout("/", http.HandlerFunc(c.Layout))
	r.Get("/", http.HandlerFunc(c.Index))
	r.Get("/faq", http.HandlerFunc(c.Faq))
}

func (c *Controller) Layout(w http.ResponseWriter, r *http.Request) {
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
	if err := c.ssr.Render(w, "view/layout.tsx", &view.Data{
		Props: map[string]interface{}{
			"page":  string(page),
			"heads": json.RawMessage(heads),
		},
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *Controller) Index(w http.ResponseWriter, r *http.Request) {
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
	if err := c.ssr.Render(slot, "view/index.tsx", &view.Data{Props: props, Slots: slot}); err != nil {
		// TODO: support live reload on error pages
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *Controller) Faq(w http.ResponseWriter, r *http.Request) {
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
	if err := c.ssr.Render(slot, "view/faq.tsx", &view.Data{Props: props, Slots: slot}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
