package controller

import (
	"encoding/json"
	"net/http"
	"reflect"
)

func wrapJSON(handler any) (http.Handler, error) {
	switch h := handler.(type) {
	case http.Handler:
		return h, nil
	case func(http.ResponseWriter, *http.Request):
		return http.HandlerFunc(h), nil
	}
	return wrapValue(defaultReader{}, newJsonWriter(), reflect.ValueOf(handler))
}

func newJsonWriter() *jsonWriter {
	return &jsonWriter{}
}

type jsonWriter struct {
}

var _ writer = (*jsonWriter)(nil)

func (*jsonWriter) WriteEmpty(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func (*jsonWriter) WriteOutput(w http.ResponseWriter, r *http.Request, out any) {
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(out); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (*jsonWriter) WriteError(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Add("Content-Type", "application/json")
	// By default treat errors as 500 errors
	status := http.StatusInternalServerError
	// Allow for custom status codes by implementing the statuser interface
	if s, ok := err.(statuser); ok {
		status = s.Status()
	}
	w.WriteHeader(status)
	props := map[string]any{"error": err.Error()}
	if err := json.NewEncoder(w).Encode(props); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
