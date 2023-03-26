package session

import (
	"fmt"
	"net/http"

	"github.com/livebud/bud/example/zero/env"
)

type Middleware struct {
	Env *env.Env
}

func (m *Middleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		// Save the session
		fmt.Println("saving session", m.Env.Session.Key)
	})
}
