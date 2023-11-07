package mux

import "net/http"

func Batch(handlers ...http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("batch not implemented yet"))
	})
}
