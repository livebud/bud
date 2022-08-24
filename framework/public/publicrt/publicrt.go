package publicrt

import (
	"net/http"

	"github.com/livebud/bud/package/budclient"
	"github.com/livebud/bud/package/middleware"
)

func Proxy(client budclient.Client) *liveMiddleware {
	return &liveMiddleware{client}
}

type liveMiddleware struct {
	client budclient.Client
}

var _ middleware.Middleware = (*liveMiddleware)(nil)

func (l *liveMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			next.ServeHTTP(w, r)
		}
		// l.client.Proxy
		// l.client.Live(r)
		next.ServeHTTP(w, r)
	})
}
