package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

type Stack []Middleware

func Compose(middlewares ...Middleware) Middleware {
	return func(handler http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			handler = middlewares[i](handler)
		}
		return handler
	}
}
