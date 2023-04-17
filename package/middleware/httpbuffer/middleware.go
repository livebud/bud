package httpbuffer

import (
	"bytes"
	"net/http"

	"github.com/felixge/httpsnoop"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/middleware"
)

func New(log log.Log) middleware.Middleware {
	rw := &responseWriter{
		code: 0,
		body: new(bytes.Buffer),
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(original http.ResponseWriter, r *http.Request) {
			w := httpsnoop.Wrap(original, httpsnoop.Hooks{
				WriteHeader: func(_ httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
					return rw.WriteHeader
				},
				Write: func(_ httpsnoop.WriteFunc) httpsnoop.WriteFunc {
					return rw.Write
				},
				Flush: func(flush httpsnoop.FlushFunc) httpsnoop.FlushFunc {
					rw.writeTo(original)
					return flush
				},
			})
			next.ServeHTTP(w, r)
			rw.writeTo(original)
		})
	}
}

type responseWriter struct {
	body  *bytes.Buffer
	code  int
	wrote bool
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.code = statusCode
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	return rw.body.Write(b)
}

func (rw *responseWriter) writeTo(w http.ResponseWriter) {
	// Only write status code once to avoid:
	// "http: superfluous response.WriteHeader"
	// Not concurrency safe.
	if !rw.wrote {
		if rw.code == 0 {
			rw.code = http.StatusOK
		}
		w.WriteHeader(rw.code)
		rw.wrote = true
	}
	rw.body.WriteTo(w)
}
