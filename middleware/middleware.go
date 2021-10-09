package middleware

import (
	"net/http"
)

// Interface for implementing middleware
type Middleware interface {
	Middleware(next http.Handler) http.Handler
}

// Function for creating middleware
type Function func(next http.Handler) http.Handler

func (fn Function) Middleware(next http.Handler) http.Handler {
	return fn(next)
}

// Stack of middleware
type Stack []Middleware

// Middleware fn
func (stack Stack) Middleware(next http.Handler) http.Handler {
	return Compose(stack...).Middleware(next)
}

// Compose a stack of middleware into a single middleware
func Compose(stack ...Middleware) Middleware {
	return Function(func(h http.Handler) http.Handler {
		if len(stack) == 0 {
			return h
		}
		for i := len(stack) - 1; i >= 0; i-- {
			h = stack[i].Middleware(h)
		}
		return h
	})
}
