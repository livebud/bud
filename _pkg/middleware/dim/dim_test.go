package dim_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"

	"github.com/livebud/bud/middleware/dim"
	"github.com/livebud/bud/pkg/di"
	"github.com/matryer/is"
	"golang.org/x/sync/errgroup"
)

func TestMiddleware(t *testing.T) {
	is := is.New(t)
	in := di.New()
	type Log struct{}
	di.Provide[*Log](in, func() (*Log, error) {
		return &Log{}, nil
	})
	middleware := dim.Provide(in)
	handler := middleware.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		in, err := dim.From(r.Context())
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
		in, err := dim.From(r.Context())
		is.True(err != nil)
		is.True(errors.Is(err, dim.ErrNotInContext))
		is.True(in == nil)
		w.WriteHeader(500)
	})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	handler.ServeHTTP(rec, req)
	res := rec.Result()
	is.Equal(res.StatusCode, 500)
}

func TestConcurrency(t *testing.T) {
	is := is.New(t)
	in := di.New()
	middleware := dim.Provide(in)
	count := 100
	mu := sync.Mutex{}
	paths := make(map[string]struct{}, count)
	requests := make([]*http.Request, count)
	for i := 0; i < count; i++ {
		path := "/" + strconv.Itoa(i)
		paths[path] = struct{}{}
		requests[i] = httptest.NewRequest(http.MethodGet, path, nil)
	}
	handler := middleware.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		in, err := dim.From(r.Context())
		is.NoErr(err)
		req, err := di.Load[*http.Request](in)
		is.NoErr(err)
		is.Equal(req.URL.Path, r.URL.Path)
		mu.Lock()
		delete(paths, req.URL.Path)
		mu.Unlock()
	}))
	eg := new(errgroup.Group)
	for _, req := range requests {
		req := req
		eg.Go(func() error {
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			res := rec.Result()
			is.Equal(res.StatusCode, 200)
			return nil
		})
	}
	err := eg.Wait()
	is.NoErr(err)
	is.Equal(len(paths), 0)
}
