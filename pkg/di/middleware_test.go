package di_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/livebud/bud/pkg/di"
	"github.com/matryer/is"
)

func TestMiddleware(t *testing.T) {
	is := is.New(t)
	in := di.New()
	type Log struct{}
	di.Provide[*Log](in, func() (*Log, error) {
		return &Log{}, nil
	})
	middleware := di.Middleware(in)
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		in, err := di.FromContext(r.Context())
		is.NoErr(err)
		log, err := di.Load[*Log](in)
		is.NoErr(err)
		is.True(log != nil)
		ctx, err := di.Load[context.Context](in)
		is.NoErr(err)
		is.True(ctx != nil)
		rw, err := di.Load[http.ResponseWriter](in)
		is.NoErr(err)
		is.True(rw == w)
		req, err := di.Load[*http.Request](in)
		is.NoErr(err)
		is.True(req == r)
		w.WriteHeader(204)
	}))
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	handler.ServeHTTP(rec, req)
	res := rec.Result()
	is.Equal(res.StatusCode, 204)
}

func TestNoMiddleware(t *testing.T) {
	is := is.New(t)
	in := di.New()
	type Log struct{}
	di.Provide[*Log](in, func() (*Log, error) {
		return &Log{}, nil
	})
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		in, err := di.FromContext(r.Context())
		is.True(err != nil)
		is.True(errors.Is(err, di.ErrNotInContext))
		is.True(in == nil)
		w.WriteHeader(500)
	})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	handler.ServeHTTP(rec, req)
	res := rec.Result()
	is.Equal(res.StatusCode, 500)
}
