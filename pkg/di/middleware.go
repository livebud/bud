package di

import (
	"context"
	"errors"
	"net/http"
)

var ErrNotInContext = errors.New("di: injector is not in the context")

type contextKey string

var key contextKey = "di"

func Middleware(parent Injector) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			in := Clone(parent)
			Provide[http.ResponseWriter](in, func(in Injector) (http.ResponseWriter, error) {
				return w, nil
			})
			Provide[*http.Request](in, func(in Injector) (*http.Request, error) {
				return r, nil
			})
			Provide[context.Context](in, func(in Injector) (context.Context, error) {
				return r.Context(), nil
			})
			r = r.WithContext(context.WithValue(r.Context(), key, in))
			next.ServeHTTP(w, r)
		})
	}
}

// FromContext returns the injector from a context if present
func FromContext(ctx context.Context) (Injector, error) {
	in, ok := ctx.Value(key).(Injector)
	if !ok {
		return nil, ErrNotInContext
	}
	return in, nil
}
