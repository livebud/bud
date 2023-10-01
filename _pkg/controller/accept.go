package controller

import (
	"net/http"

	"github.com/livebud/bud/pkg/controller/internal/request"
)

type accept struct {
	HTML writer
	JSON writer
}

var _ writer = (*accept)(nil)

func (a *accept) negotiate(r *http.Request) (writer, bool) {
	ctype, err := request.Negotiate(r, "text/html", "application/json")
	if err != nil {
		return nil, false
	}
	switch ctype {
	case "text/html":
		return a.HTML, true
	case "application/json":
		return a.JSON, true
	default:
		return nil, false
	}
}

func (a *accept) WriteEmpty(w http.ResponseWriter, r *http.Request) {
	writer, ok := a.negotiate(r)
	if !ok {
		http.Error(w, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable)
		return
	}
	writer.WriteEmpty(w, r)
}

func (a *accept) WriteOutput(w http.ResponseWriter, r *http.Request, out any) {
	writer, ok := a.negotiate(r)
	if !ok {
		http.Error(w, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable)
		return
	}
	writer.WriteOutput(w, r, out)
}

func (a *accept) WriteError(w http.ResponseWriter, r *http.Request, err error) {
	writer, ok := a.negotiate(r)
	if !ok {
		http.Error(w, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable)
		return
	}
	writer.WriteError(w, r, err)
}
