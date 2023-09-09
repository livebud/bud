package middleware

import (
	"net/http"
)

type Func func(http.Handler) http.Handler

func (f Func) Middleware(h http.Handler) http.Handler {
	return f(h)
}

type Middleware interface {
	Middleware(http.Handler) http.Handler
}

// compose a stack of middleware into a single middleware
func compose(middlewares ...Middleware) Middleware {
	return Func(func(h http.Handler) http.Handler {
		if len(middlewares) == 0 {
			return h
		}
		for i := len(middlewares) - 1; i >= 0; i-- {
			if middlewares[i] == nil {
				continue
			}
			h = middlewares[i].Middleware(h)
		}
		return h
	})
}

type Stack []Middleware

func (s Stack) Middleware(h http.Handler) http.Handler {
	return compose(s...).Middleware(h)
}
