package di_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/livebud/bud/pkg/di"
	"github.com/matryer/is"
)

func TestContext(t *testing.T) {
	is := is.New(t)
	in := di.New()
	di.Loader[*Env](in, loadEnv)
	env, err := di.Load[*Env](in)
	is.NoErr(err)
	is.True(env != nil)
	env, err = di.Load[*Env](in)
	is.NoErr(err)
	is.True(env != nil)
	ctx := di.WithInjector(context.Background(), in)
	env, err = di.LoadFrom[*Env](ctx)
	is.NoErr(err)
	is.True(env != nil)
}

type Stack struct {
	fns []func(http.Handler) http.Handler
}

func (s *Stack) Append(fn func(http.Handler) http.Handler) {
	s.fns = append(s.fns, fn)
}

func (s *Stack) Compose(bottom http.Handler) http.Handler {
	stack := func(next http.Handler) http.Handler {
		for _, m := range s.fns {
			next = m(next)
		}
		return next
	}
	return stack(bottom)
}

func (s *Stack) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Compose(http.NotFoundHandler()).ServeHTTP(w, r)
}
