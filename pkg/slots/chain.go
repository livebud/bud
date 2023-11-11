package slots

import (
	"io"
	"net/http"

	"github.com/felixge/httpsnoop"
	"github.com/livebud/bud/pkg/request"
)

func Chain(handlers ...http.Handler) http.Handler {
	if len(handlers) == 1 {
		return handlers[0]
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slot := New()
		status := 0
		headers := w.Header()
		for i := 0; i < len(handlers); i++ {
			handler := handlers[i]
			r := request.Clone(r.WithContext(ToContext(r.Context(), slot)))
			innerHeaders := http.Header{}
			innerStatus := 0
			w := httpsnoop.Wrap(w, httpsnoop.Hooks{
				WriteHeader: func(next httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
					return func(code int) {
						innerStatus = code
					}
				},
				Header: func(next httpsnoop.HeaderFunc) httpsnoop.HeaderFunc {
					return func() http.Header {
						return innerHeaders
					}
				},
				Write: func(next httpsnoop.WriteFunc) httpsnoop.WriteFunc {
					return slot.Write
				},
			})
			handler.ServeHTTP(w, r)
			// Prioritize 400s and 500s over 200s
			if innerStatus > status {
				status = innerStatus
			}
			for key := range innerHeaders {
				headers.Set(key, innerHeaders.Get(key))
			}
			slot.Close()
			slot = slot.Next()
		}
		if status > 0 {
			w.WriteHeader(status)
		}
		body, err := io.ReadAll(slot)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(body)
	})
}

func newResponseWriter(w http.ResponseWriter) http.ResponseWriter {
	return httpsnoop.Wrap(w, httpsnoop.Hooks{})
}

type responseWriter struct {
	http.ResponseWriter
}

// func
