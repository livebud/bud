package hot_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/hot"
	"github.com/livebud/bud/package/socket"
)

var now = time.Date(2021, 8, 4, 14, 56, 0, 0, time.UTC)

func TestServer(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	hotServer := hot.New()
	hotServer.Now = func() time.Time { return now }
	testServer := httptest.NewServer(hotServer)
	hotClient, err := hot.Dial(testServer.URL)
	is.NoErr(err)
	hotServer.Reload("*")
	event, err := hotClient.Next(ctx)
	is.NoErr(err)
	is.Equal(event.ID, "")
	is.Equal(event.Type, "")
	is.Equal(string(event.Data), `{"scripts":["?ts=1628088960000"]}`)
	is.Equal(event.Retry, 0)
	is.NoErr(hotClient.Close())
	testServer.Close()
}

func TestPage(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	hotServer := hot.New()
	hotServer.Now = func() time.Time { return now }
	testServer := httptest.NewServer(hotServer)
	query := url.Values{}
	query.Add("page", "/bud/view/index.svelte")
	hotClient, err := hot.Dial(testServer.URL + "?" + query.Encode())
	is.NoErr(err)
	hotServer.Reload("bud/view/index.svelte")
	event, err := hotClient.Next(ctx)
	is.NoErr(err)
	is.Equal(event.ID, "")
	is.Equal(event.Type, "")
	is.Equal(string(event.Data), `{"scripts":["/bud/view/index.svelte?ts=1628088960000"]}`)
	is.Equal(event.Retry, 0)
	hotServer.Reload("bud/view/index.svelte")
	event, err = hotClient.Next(ctx)
	is.NoErr(err)
	is.Equal(event.ID, "")
	is.Equal(event.Type, "")
	is.Equal(string(event.Data), `{"scripts":["/bud/view/index.svelte?ts=1628088960000"]}`)
	is.Equal(event.Retry, 0)
	hotServer.Reload("*")
	event, err = hotClient.Next(ctx)
	is.NoErr(err)
	is.Equal(event.ID, "")
	is.Equal(event.Type, "")
	is.Equal(string(event.Data), `{"scripts":["/bud/view/index.svelte?ts=1628088960000"]}`)
	is.Equal(event.Retry, 0)
	is.NoErr(hotClient.Close())
	testServer.Close()
}

func TestReload(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	hotServer := hot.New()
	hotServer.Now = func() time.Time { return now }
	testServer := httptest.NewServer(hotServer)
	query := url.Values{}
	query.Add("page", "/bud/view/index.svelte")
	hotClient, err := hot.Dial(testServer.URL + "?" + query.Encode())
	is.NoErr(err)
	hotServer.Reload("!")
	event, err := hotClient.Next(ctx)
	is.NoErr(err)
	is.Equal(event.ID, "")
	is.Equal(event.Type, "")
	is.Equal(string(event.Data), `{"reload":true}`)
	is.Equal(event.Retry, 0)
	is.NoErr(hotClient.Close())
	testServer.Close()
}

// TODO: consolidate function. This is duplicated in multiple places.
func listen(path string) (socket.Listener, *http.Client, error) {
	listener, err := socket.Listen(path)
	if err != nil {
		return nil, nil, err
	}
	transport, err := socket.Transport(listener.Addr().String())
	if err != nil {
		return nil, nil, err
	}
	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	return listener, client, nil
}

func TestUnixListener(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	listener, client, err := listen(filepath.Join(t.TempDir(), "test.sock"))
	is.NoErr(err)
	hotServer := hot.New()
	hotServer.Now = func() time.Time { return now }
	server := &http.Server{
		Addr:    listener.Addr().String(),
		Handler: hotServer,
	}
	defer server.Shutdown(ctx)
	eg := new(errgroup.Group)
	eg.Go(func() error { return server.Serve(listener) })
	hotClient, err := hot.DialWith(client, "http://host/")
	is.NoErr(err)
	hotServer.Reload("*")
	event, err := hotClient.Next(ctx)
	is.NoErr(err)
	is.Equal(event.ID, "")
	is.Equal(event.Type, "")
	is.Equal(string(event.Data), `{"scripts":["?ts=1628088960000"]}`)
	is.Equal(event.Retry, 0)
	is.NoErr(hotClient.Close())
	is.NoErr(server.Shutdown(ctx))
}

func TestNoWait(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	listener, client, err := listen(filepath.Join(t.TempDir(), "test.sock"))
	is.NoErr(err)
	hotServer := hot.New()
	hotServer.Now = func() time.Time { return now }
	server := &http.Server{
		Addr:    listener.Addr().String(),
		Handler: hotServer,
	}
	defer server.Shutdown(ctx)
	eg := new(errgroup.Group)
	eg.Go(func() error { return server.Serve(listener) })
	hotClient, err := hot.DialWith(client, "http://host/")
	is.NoErr(err)
	is.NoErr(hotClient.Close())
}

func TestDrainBeforeClose(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	listener, client, err := listen(filepath.Join(t.TempDir(), "test.sock"))
	is.NoErr(err)
	hotServer := hot.New()
	hotServer.Now = func() time.Time { return now }
	server := &http.Server{
		Addr:    listener.Addr().String(),
		Handler: hotServer,
	}
	defer server.Shutdown(ctx)
	eg := new(errgroup.Group)
	eg.Go(func() error { return server.Serve(listener) })
	hotClient, err := hot.DialWith(client, "http://host/")
	is.NoErr(err)
	hotServer.Reload("*")
	hotServer.Reload("*")
	hotServer.Reload("*")
	hotServer.Reload("*")
	is.NoErr(hotClient.Close())
}
