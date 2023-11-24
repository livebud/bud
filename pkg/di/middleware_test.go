package di_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/livebud/bud/pkg/di"
	"github.com/matryer/is"
)

func TestMiddleware(t *testing.T) {
	is := is.New(t)
	in := di.New()
	di.Loader[*Env](in, loadEnv)
	err := di.Loader[*Stack](in, func(in di.Injector) (*Stack, error) {
		return &Stack{}, nil
	})
	is.NoErr(err)
	called := 0
	// Attach the injector
	di.Append[*Stack](in, func(in di.Injector, s *Stack) error {
		s.Append(di.Middleware(in))
		called++
		return nil
	})
	stack, err := di.Load[*Stack](in)
	is.NoErr(err)
	h := stack.Compose(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		env, err := di.LoadFrom[*Env](r.Context())
		is.NoErr(err)
		is.True(env != nil)
		called++
	}))
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	is.Equal(called, 2)
}
