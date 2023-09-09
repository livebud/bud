package controller

import (
	"net/http"
	"reflect"
)

type formWriter struct {
	handler reflect.Value
}

var _ writer = (*formWriter)(nil)

func (f *formWriter) WriteEmpty(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Location", backPath(r))
	w.WriteHeader(http.StatusSeeOther)
}

func (f *formWriter) WriteOutput(w http.ResponseWriter, r *http.Request, out any) {
	// TODO: look at `out` during a POST/PATCH for the redirect key
	w.Header().Add("Location", backPath(r))
	w.WriteHeader(http.StatusSeeOther)
}

func (f *formWriter) WriteError(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Add("Location", backPath(r))
	// TODO: add session flash
	w.WriteHeader(http.StatusSeeOther)
}

func backPath(r *http.Request) (location string) {
	if r.Referer() != "" {
		return r.Referer()
	}
	return r.URL.Path
}
