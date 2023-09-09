package csrf

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/livebud/bud/pkg/middleware"
	"github.com/livebud/bud/pkg/view"
)

func Default() *Middleware {
	return &Middleware{
		CookieName: "_csrf",
		FormKey:    "csrf",
		Generate: func() (string, error) {
			b := make([]byte, 32)
			if _, err := rand.Read(b); err != nil {
				return "", err
			}
			return base64.StdEncoding.EncodeToString(b), nil
		},
	}
}

type Middleware struct {
	CookieName string
	FormKey    string
	Generate   func() (string, error)
}

var _ middleware.Middleware = (*Middleware)(nil)

func (m *Middleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(view.Set(r.Context(), view.Data{
			m.FormKey: "123",
		}))
		next.ServeHTTP(w, r)
	})
}
