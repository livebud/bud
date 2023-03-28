package session

import (
	"fmt"
	"net/http"

	"github.com/livebud/bud/example/zero/bud/pkg/sessions"
	"github.com/livebud/bud/example/zero/env"
)

type Middleware struct {
	Env      *env.Env
	Sessions *sessions.Store
	// Sessions *session.Store
	// Store session.Storage
}

func (m *Middleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := m.Sessions.Load(r, "some_id")
		if err != nil {
			fmt.Println("error loading session", err)
		}
		fmt.Println("loaded session", session)
		next.ServeHTTP(w, r)
		// Save the session
		fmt.Println("saving session", m.Env.Session.Key)
	})
}
