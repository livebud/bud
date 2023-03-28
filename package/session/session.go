package session

import (
	"errors"
	"net/http"
	"time"

	"github.com/livebud/bud/package/cookies"
)

var ErrNotFound = errors.New("session not found")

type Store interface {
	Get(id string) ([]byte, error)
	Set(id string, payload []byte, expires time.Time) error
	Delete(id string) error
}

// New cookie store
func New(cs cookies.Store, w http.ResponseWriter, r *http.Request) *CookieStore {
	return &CookieStore{cs, r, w}
}

// CookieStore is a cookie store
type CookieStore struct {
	cs cookies.Store
	r  *http.Request
	w  http.ResponseWriter
}

var _ Store = (*CookieStore)(nil)

func (c *CookieStore) Get(id string) ([]byte, error) {
	cookie, err := c.cs.Get(c.r, id)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return []byte(cookie.Value), nil
}

func (c *CookieStore) Set(id string, payload []byte, expires time.Time) error {
	cookie := &http.Cookie{
		Name:    id,
		Value:   string(payload),
		Expires: time.Now().Add(24 * time.Hour),
	}
	return c.cs.Set(c.w, cookie)
}

func (c *CookieStore) Delete(id string) error {
	cookie := &http.Cookie{
		Name:   id,
		MaxAge: -1,
	}
	return c.cs.Set(c.w, cookie)
}
