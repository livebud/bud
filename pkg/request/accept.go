package request

import (
	"net/http"

	"github.com/timewasted/go-accept-headers"
)

// Accept returns the negotiated content type from the request. Unlike
// negotiate, if there's an error parsing the header, an empty string is
// returned.
func Accept(r *http.Request, types ...string) string {
	accept, err := Negotiate(r, types...)
	if err != nil {
		return ""
	}
	return accept
}

// Negotiate returns the negotiated content type from the request.
func Negotiate(r *http.Request, types ...string) (string, error) {
	header := r.Header.Get("Accept")
	if header == "" {
		header = r.Header.Get("Content-Type")
	}
	return accept.Negotiate(header, types...)
}
