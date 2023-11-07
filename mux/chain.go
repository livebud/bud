package mux

import (
	"net/http"

	"github.com/livebud/bud/internal/httpwrap"
)

func Chain(handlers ...http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw, flush := httpwrap.Wrap(w)
		defer flush()
		for i := 0; i < len(handlers); i++ {
			handlers[i].ServeHTTP(rw, r.Clone(r.Context()))
		}
	})
}
