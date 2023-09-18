package cookies

import (
	"net/http"
	"time"
)

// Interface is an interface for getting and setting cookies
type Interface interface {
	Read(r *http.Request, name string) (*http.Cookie, error)
	Write(w http.ResponseWriter, cookie *http.Cookie) error
}

func New() *Cookies {
	return &Cookies{
		ExpiresIn: 30 * 24 * time.Hour,
		Path:      "/",
		HttpOnly:  true,
	}
}

type Cookies struct {
	ExpiresIn time.Duration
	Path      string
	HttpOnly  bool
}

func (c *Cookies) Read(r *http.Request, name string) (*http.Cookie, error) {
	return r.Cookie(name)
}

func (c *Cookies) Write(w http.ResponseWriter, cookie *http.Cookie) error {
	http.SetCookie(w, cookie)
	return nil
}
