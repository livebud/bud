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

type redirector interface {
	Redirect(r *http.Request) string
}

func (f *formWriter) WriteOutput(w http.ResponseWriter, r *http.Request, out any) {
	redirectPath := r.URL.Path
	if redirector, ok := out.(redirector); ok {
		redirectPath = redirector.Redirect(r)
	}
	w.Header().Add("Location", redirectPath)
	w.WriteHeader(http.StatusSeeOther)
}

func (f *formWriter) WriteError(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Add("Location", backPath(r))
	// TODO: add session flash
	w.WriteHeader(http.StatusSeeOther)
	w.Write([]byte(err.Error()))
}

func backPath(r *http.Request) (location string) {
	if r.Referer() != "" {
		return r.Referer()
	}
	return r.URL.Path
}
