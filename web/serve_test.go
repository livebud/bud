package web_test

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/tj/assert"
	"golang.org/x/sync/errgroup"

	"gitlab.com/mnm/bud/web"
)

// get an available port
func freePort(t testing.TB) int {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	assert.NoError(t, err)
	l, err := net.ListenTCP("tcp", addr)
	assert.NoError(t, err)
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

func TestServeTCP(t *testing.T) {
	is := is.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	url := "localhost:" + strconv.Itoa(freePort(t))
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(205)
	})
	eg := new(errgroup.Group)
	eg.Go(func() error { return web.ServeTCP(ctx, url, handler) })
	res, err := http.Get("http://" + url)
	is.NoErr(err)
	is.Equal(res.StatusCode, 205)
	cancel()
	eg.Wait()
	res, err = http.Get("http://" + url)
	is.True(err != nil)
	is.True(res == nil)
	is.True(strings.Contains(err.Error(), `connection refused`)) // should have stopped
}
