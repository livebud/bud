package hot_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/pubsub"
	"github.com/livebud/bud/package/hot"
	"github.com/livebud/bud/package/log/testlog"
	"github.com/livebud/bud/package/socket"
)

var now = time.Date(2021, 8, 4, 14, 56, 0, 0, time.UTC)

func TestNoPathUpdate(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ps := pubsub.New()
	hotServer := hot.New(log, ps)
	hotServer.Now = func() time.Time { return now }
	testServer := httptest.NewServer(hotServer)
	hotClient, err := hot.Dial(log, testServer.URL)
	is.NoErr(err)
	ps.Publish("frontend:update", nil)
	event, err := hotClient.Next(ctx)
	is.NoErr(err)
	is.Equal(event.ID, "")
	is.Equal(event.Type, "")
	is.Equal(string(event.Data), `{"reload":true}`)
	is.Equal(event.Retry, 0)
	is.NoErr(hotClient.Close())
	testServer.Close()
}

func TestPage(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ps := pubsub.New()
	hotServer := hot.New(log, ps)
	hotServer.Now = func() time.Time { return now }
	testServer := httptest.NewServer(hotServer)
	hotClient, err := hot.Dial(log, testServer.URL+"/bud/hot/view/index.svelte")
	is.NoErr(err)
	ps.Publish("frontend:update:view/index.svelte", nil)
	event, err := hotClient.Next(ctx)
	is.NoErr(err)
	is.Equal(event.ID, "")
	is.Equal(event.Type, "")
	is.Equal(string(event.Data), `{"scripts":["/bud/view/index.svelte?ts=1628088960000"]}`)
	is.Equal(event.Retry, 0)
	ps.Publish("frontend:update:view/index.svelte", nil)
	event, err = hotClient.Next(ctx)
	is.NoErr(err)
	is.Equal(event.ID, "")
	is.Equal(event.Type, "")
	is.Equal(string(event.Data), `{"scripts":["/bud/view/index.svelte?ts=1628088960000"]}`)
	is.Equal(event.Retry, 0)
	ps.Publish("frontend:update", nil)
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
	log := testlog.New()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ps := pubsub.New()
	hotServer := hot.New(log, ps)
	hotServer.Now = func() time.Time { return now }
	testServer := httptest.NewServer(hotServer)
	hotClient, err := hot.Dial(log, testServer.URL+`/bud/hot/view/index.svelte`)
	is.NoErr(err)
	ps.Publish("backend:update", nil)
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
	log := testlog.New()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	listener, client, err := listen(filepath.Join(t.TempDir(), "test.sock"))
	is.NoErr(err)
	ps := pubsub.New()
	hotServer := hot.New(log, ps)
	hotServer.Now = func() time.Time { return now }
	server := &http.Server{
		Addr:    listener.Addr().String(),
		Handler: hotServer,
	}
	defer server.Shutdown(ctx)
	eg := new(errgroup.Group)
	eg.Go(func() error { return server.Serve(listener) })
	hotClient, err := hot.DialWith(client, log, "http://host/")
	is.NoErr(err)
	ps.Publish("frontend:update", nil)
	event, err := hotClient.Next(ctx)
	is.NoErr(err)
	is.Equal(event.ID, "")
	is.Equal(event.Type, "")
	is.Equal(string(event.Data), `{"reload":true}`)
	is.Equal(event.Retry, 0)
	is.NoErr(hotClient.Close())
	is.NoErr(server.Shutdown(ctx))
}

func TestNoWait(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	listener, client, err := listen(filepath.Join(t.TempDir(), "test.sock"))
	is.NoErr(err)
	ps := pubsub.New()
	hotServer := hot.New(log, ps)
	hotServer.Now = func() time.Time { return now }
	server := &http.Server{
		Addr:    listener.Addr().String(),
		Handler: hotServer,
	}
	defer server.Shutdown(ctx)
	eg := new(errgroup.Group)
	eg.Go(func() error { return server.Serve(listener) })
	hotClient, err := hot.DialWith(client, log, "http://host/")
	is.NoErr(err)
	is.NoErr(hotClient.Close())
}

func TestDrainBeforeClose(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	listener, client, err := listen(filepath.Join(t.TempDir(), "test.sock"))
	is.NoErr(err)
	ps := pubsub.New()
	hotServer := hot.New(log, ps)
	hotServer.Now = func() time.Time { return now }
	server := &http.Server{
		Addr:    listener.Addr().String(),
		Handler: hotServer,
	}
	defer server.Shutdown(ctx)
	eg := new(errgroup.Group)
	eg.Go(func() error { return server.Serve(listener) })
	hotClient, err := hot.DialWith(client, log, "http://host/")
	is.NoErr(err)
	ps.Publish("frontend:update", nil)
	ps.Publish("frontend:update", nil)
	ps.Publish("frontend:update", nil)
	ps.Publish("frontend:update", nil)
	is.NoErr(hotClient.Close())
}
