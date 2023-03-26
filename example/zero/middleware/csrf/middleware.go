package csrf

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
		fmt.Println("checking csrf", m.Env.CSRF.Token)
		next.ServeHTTP(w, r)
	})
}
