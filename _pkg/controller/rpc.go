package controller

import (
	"net/http"
	"reflect"
)

type reader interface {
	ReadContext(r *http.Request, t reflect.Type) (reflect.Value, error)
	ReadInput(r *http.Request, t reflect.Type) (reflect.Value, error)
}

type writer interface {
	WriteEmpty(w http.ResponseWriter, r *http.Request)
	WriteOutput(w http.ResponseWriter, r *http.Request, out any)
	WriteError(w http.ResponseWriter, r *http.Request, err error)
}

func getStatus(defaultStatus int, err error) (status int) {
	// Allow for custom status codes by implementing the statuser interface
	if s, ok := err.(statuser); ok {
		return s.Status()
	}
	return defaultStatus
}
