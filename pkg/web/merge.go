package web

import (
	"net/http"

	"github.com/livebud/bud/internal/httpwrap"
)

func Merge(handlers ...http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw, flush := httpwrap.Wrap(w)
		defer flush()
		for i := len(handlers) - 1; i >= 0; i-- {
			handlers[i].ServeHTTP(rw, r.Clone(r.Context()))
		}
	})
}
