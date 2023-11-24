package di

import "net/http"

// Middleware creates middleware to inject the injector into the request context
func Middleware(parent Injector) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			in := Clone(parent)
			r = r.WithContext(WithInjector(r.Context(), in))
			next.ServeHTTP(w, r)
		})
	}
}
