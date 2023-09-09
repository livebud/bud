package controller

import (
	"net/http"

	"github.com/livebud/bud/pkg/web"
)

func New() *Controller {
	return &Controller{}
}

func Register(router web.Router, controller *Controller) {
	router.Get("/", http.HandlerFunc(controller.Index))
}

type Controller struct {
}

func (c *Controller) Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}
