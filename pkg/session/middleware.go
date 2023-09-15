package session

import (
	"context"
	"net/http"

	"github.com/ajg/form"
	"github.com/livebud/bud/pkg/middleware"
)

const cookieName = "sid"

type contextKey string

var key contextKey = "session"

type Middleware middleware.Middleware

type wrapper struct {
	Raw     string
	Session any
}

func New() Middleware {
	return middleware.Func(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var value string
			if cookie, err := r.Cookie(cookieName); nil == err {
				value = cookie.Value
			}
			container := &wrapper{
				Raw: value,
			}
			r = r.WithContext(context.WithValue(r.Context(), key, container))
			next.ServeHTTP(w, r)
			sessionValue := r.Context().Value(key)
			if sessionValue == nil {
				return
			}
			container, ok := sessionValue.(*wrapper)
			if !ok {
				return
			} else if container.Session == nil {
				return
			}
			newValue, err := form.EncodeToString(container.Session)
			if err != nil {
				// TODO: Handle this error.
				panic(err)
			}
			http.SetCookie(w, &http.Cookie{
				Name:     cookieName,
				Value:    newValue,
				HttpOnly: true,
			})
		})
	})
}
