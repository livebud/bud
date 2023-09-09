package request

import (
	"net/http"

	"github.com/timewasted/go-accept-headers"
)

func Negotiate(r *http.Request, ctypes ...string) (string, error) {
	header := r.Header.Get("Accept")
	if header == "" {
		header = r.Header.Get("Content-Type")
	}
	return accept.Negotiate(header, ctypes...)
}
