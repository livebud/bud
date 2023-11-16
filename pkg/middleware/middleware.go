package middleware

import "net/http"

// type Middleware = func(next http.Handler) http.Handler

// Middleware handler
// type Middleware interface {
// 	Middleware(next http.Handler) http.Handler
// }

// // Func is a middleware function
// type Func func(next http.Handler) http.Handler

// func (fn Func) Middleware(next http.Handler) http.Handler {
// 	return fn(next)
// }

// Stack of middleware
// type Stack []Middleware

// // Middleware fn
// func (stack Stack) Middleware(next http.Handler) http.Handler {
// 	return Compose(stack...).Middleware(next)
// }

func Compose(stack ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		if len(stack) == 0 {
			return h
		}
		for i := len(stack) - 1; i >= 0; i-- {
			if stack[i] == nil {
				continue
			}
			h = stack[i](h)
		}
		return h
	}
}

// Compose a stack of middleware into a single middleware
// func Compose(stack ...Middleware) Middleware {
// 	return Func(func(h http.Handler) http.Handler {
// 		if len(stack) == 0 {
// 			return h
// 		}
// 		for i := len(stack) - 1; i >= 0; i-- {
// 			if stack[i] == nil {
// 				continue
// 			}
// 			h = stack[i].Middleware(h)
// 		}
// 		return h
// 	})
// }
