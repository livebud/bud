package request

import (
	"net/http"

	"github.com/timewasted/go-accept-headers"
)

// New request context
func New(r *http.Request) *Context {
	return &Context{r}
}

// Context struct
type Context struct {
	r *http.Request
}

// Unmarshal the request body or parameters
func (c *Context) Unmarshal(r *http.Request, in interface{}) error {
	return Unmarshal(r, in)
}

// Accepts a type
func Accepts(r *http.Request) Acceptable {
	return Acceptable(accept.Parse(r.Header.Get("Accept")))
}

// Acceptable types
type Acceptable accept.AcceptSlice

// Accepts checks if the content type is acceptable
func (as Acceptable) Accepts(ctype string) bool {
	return accept.AcceptSlice(as).Accepts(ctype)
}
