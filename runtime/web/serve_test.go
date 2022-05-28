package web_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/livebud/bud/internal/is"
	"golang.org/x/sync/errgroup"

	"github.com/livebud/bud/runtime/web"
)

func TestServe(t *testing.T) {
	is := is.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	listener, err := web.Listen("APP", ":0")
	is.NoErr(err)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(205)
	})
	eg := new(errgroup.Group)
	eg.Go(func() error { return web.Serve(ctx, listener, handler) })
	res, err := http.Get("http://" + listener.Addr().String())
	is.NoErr(err)
	is.Equal(res.StatusCode, 205)
	cancel()
	eg.Wait()
	res, err = http.Get("http://" + listener.Addr().String())
	is.True(err != nil)
	is.True(res == nil)
	is.True(strings.Contains(err.Error(), `connection refused`)) // should have stopped
}
