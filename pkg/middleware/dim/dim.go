package dim

import (
	"context"
	"errors"
	"net/http"

	"github.com/livebud/bud/pkg/di"
	"github.com/livebud/bud/pkg/middleware"
)

var ErrNotInContext = errors.New("dim: injector is not in the context")

type Injector di.Injector

type Middleware middleware.Middleware

type contextKey string

const key contextKey = "dim"

func From(ctx context.Context) (Injector, error) {
	in, ok := ctx.Value(key).(Injector)
	if !ok {
		return nil, ErrNotInContext
	}
	return in, nil
}

func Provide(parent di.Injector) Middleware {
	return middleware.Func(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			parent := di.Clone(parent)
			di.Provide[http.ResponseWriter](parent, func() (http.ResponseWriter, error) {
				return w, nil
			})
			di.Provide[*http.Request](parent, func() (*http.Request, error) {
				return r, nil
			})
			di.Provide[context.Context](parent, func() (context.Context, error) {
				return r.Context(), nil
			})
			di.Provide[Injector](parent, func() Injector { return parent })
			in, err := di.Load[Injector](parent)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			r = r.WithContext(context.WithValue(r.Context(), key, in))
			next.ServeHTTP(w, r)
		})
	})
}
