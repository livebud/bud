package httpwrap

import (
	"bytes"
	"net/http"

	"github.com/felixge/httpsnoop"
)

func Wrap(w http.ResponseWriter) (http.ResponseWriter, func()) {
	rb := &responseBuffer{
		code: 0,
		body: new(bytes.Buffer),
	}
	rw := httpsnoop.Wrap(w, httpsnoop.Hooks{
		WriteHeader: func(_ httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
			return rb.WriteHeader
		},
		Write: func(_ httpsnoop.WriteFunc) httpsnoop.WriteFunc {
			return rb.Write
		},
		Flush: func(flush httpsnoop.FlushFunc) httpsnoop.FlushFunc {
			rb.writeTo(w)
			return flush
		},
	})
	return rw, func() {
		rb.writeTo(w)
	}
}

type responseBuffer struct {
	body  *bytes.Buffer
	code  int
	wrote bool
}

func (rb *responseBuffer) WriteHeader(statusCode int) {
	rb.code = statusCode
}

func (rb *responseBuffer) Write(b []byte) (int, error) {
	return rb.body.Write(b)
}

func (rb *responseBuffer) writeTo(w http.ResponseWriter) {
	// Only write status code once to avoid: "http: superfluous
	// response.WriteHeader". Not concurrency safe.
	if !rb.wrote {
		if rb.code == 0 {
			rb.code = http.StatusOK
		}
		w.WriteHeader(rb.code)
		rb.wrote = true
	}
	rb.body.WriteTo(w)
}
