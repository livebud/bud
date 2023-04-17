package mw

import (
	"net/http"

	"github.com/livebud/bud/example/zero/env"
)

type Middleware struct {
	*env.Env
	// *sessions.Store
}

// func (m *Middleware) Session(next http.Handler) http.Handler {
// 	fmt.Println(m.Store)
// 	// store := sessions.NewCookieStore([]byte(m.Env.Session.Key))
// 	// _ = store
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 	})
// }

func (m *Middleware) CSRF(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}

func (m *Middleware) WrapRW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}
