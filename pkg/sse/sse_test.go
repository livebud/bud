package sse_test

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/livebud/bud/pkg/logs"
	"github.com/livebud/bud/pkg/sse"
	"github.com/matryer/is"
)

// var now = time.Date(2021, 8, 4, 14, 56, 0, 0, time.UTC)

func TestHandlerClient(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	log := logs.Discard()
	handler := sse.New(log)
	server := httptest.NewServer(handler)
	defer server.Close()
	stream, err := sse.Dial(log, server.URL)
	is.NoErr(err)
	defer stream.Close()
	err = handler.Publish(ctx, &sse.Event{
		Type: "test",
		Data: []byte("hello"),
	})
	is.NoErr(err)
	event, err := stream.Next(ctx)
	is.NoErr(err)
	is.Equal(event.Type, "test")
	is.Equal(string(event.Data), "hello")
	is.NoErr(stream.Close())
}

func TestMultipleEvents(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	log := logs.Discard()
	handler := sse.New(log)
	server := httptest.NewServer(handler)
	defer server.Close()
	stream, err := sse.Dial(log, server.URL)
	is.NoErr(err)
	defer stream.Close()
	err = handler.Publish(ctx, &sse.Event{
		Data: []byte("1"),
	})
	is.NoErr(err)
	err = handler.Publish(ctx, &sse.Event{
		Data: []byte("2"),
	})
	is.NoErr(err)
	err = handler.Publish(ctx, &sse.Event{
		Data: []byte("3"),
	})
	is.NoErr(err)
	event, err := stream.Next(ctx)
	is.NoErr(err)
	is.Equal(string(event.Data), "1")
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(200*time.Millisecond))
	defer cancel()
	// We should expect the deadline to elapse because broadcast doesn't buffer
	// events, streams must be listening already
	event, err = stream.Next(ctx)
	is.True(err != nil)
	is.True(errors.Is(err, context.DeadlineExceeded))
	is.Equal(event, nil)
	is.NoErr(stream.Close())
}

func TestNoLockup(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	log := logs.Discard()
	handler := sse.New(log)
	server := httptest.NewServer(handler)
	defer server.Close()
	stream, err := sse.Dial(log, server.URL)
	is.NoErr(err)
	defer stream.Close()
	err = handler.Publish(ctx, &sse.Event{
		Data: []byte("1"),
	})
	is.NoErr(err)
	err = handler.Publish(ctx, &sse.Event{
		Data: []byte("2"),
	})
	is.NoErr(err)
	err = handler.Publish(ctx, &sse.Event{
		Data: []byte("3"),
	})
	is.NoErr(err)
	event, err := stream.Next(ctx)
	is.NoErr(err)
	is.Equal(string(event.Data), "1")
	is.NoErr(stream.Close())
}

func TestMultipleClients(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	log := logs.Debugger()
	handler := sse.New(log)
	server := httptest.NewServer(handler)
	defer server.Close()
	stream1, err := sse.Dial(log, server.URL)
	is.NoErr(err)
	defer stream1.Close()
	stream2, err := sse.Dial(log, server.URL)
	is.NoErr(err)
	defer stream2.Close()
	err = handler.Publish(ctx, &sse.Event{
		Data: []byte("1"),
	})
	is.NoErr(err)
	event, err := stream1.Next(ctx)
	is.NoErr(err)
	is.Equal(string(event.Data), "1")
	event, err = stream2.Next(ctx)
	is.NoErr(err)
	is.Equal(string(event.Data), "1")

	is.NoErr(stream1.Close())

	err = handler.Publish(ctx, &sse.Event{
		Data: []byte("2"),
	})
	is.NoErr(err)
	event, err = stream1.Next(ctx)
	is.True(err != nil)
	is.True(errors.Is(err, sse.ErrStreamClosed))
	is.Equal(event, nil)
	event, err = stream2.Next(ctx)
	is.NoErr(err)
	is.Equal(string(event.Data), "2")
	is.NoErr(stream2.Close())

	err = handler.Publish(ctx, &sse.Event{
		Data: []byte("3"),
	})
	is.NoErr(err)
	event, err = stream1.Next(ctx)
	is.True(err != nil)
	is.True(errors.Is(err, sse.ErrStreamClosed))
	is.Equal(event, nil)
	event, err = stream2.Next(ctx)
	is.True(err != nil)
	is.True(errors.Is(err, sse.ErrStreamClosed))
	is.Equal(event, nil)
}

// func TestNoPathUpdate(t *testing.T) {
// 	is := is.New(t)
// 	log := testlogs.New()
// 	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
// 	defer cancel()
// 	ps := pubsub.New()
// 	hotServer := hot.New(log, ps)
// 	hotServer.Now = func() time.Time { return now }
// 	testServer := httptest.NewServer(hotServer)
// 	hotClient, err := hot.Dial(log, testServer.URL)
// 	is.NoErr(err)
// 	ps.Publish("frontend:update", nil)
// 	event, err := hotClient.Next(ctx)
// 	is.NoErr(err)
// 	is.Equal(event.ID, "")
// 	is.Equal(event.Type, "")
// 	is.Equal(string(event.Data), `{"reload":true}`)
// 	is.Equal(event.Retry, 0)
// 	is.NoErr(hotClient.Close())
// 	testServer.Close()
// }

// func TestPage(t *testing.T) {
// 	is := is.New(t)
// 	log := testlogs.New()
// 	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
// 	defer cancel()
// 	ps := pubsub.New()
// 	hotServer := hot.New(log, ps)
// 	hotServer.Now = func() time.Time { return now }
// 	testServer := httptest.NewServer(hotServer)
// 	hotClient, err := hot.Dial(log, testServer.URL+"/bud/hot/view/index.svelte")
// 	is.NoErr(err)
// 	ps.Publish("frontend:update:view/index.svelte", nil)
// 	event, err := hotClient.Next(ctx)
// 	is.NoErr(err)
// 	is.Equal(event.ID, "")
// 	is.Equal(event.Type, "")
// 	is.Equal(string(event.Data), `{"scripts":["/bud/view/index.svelte?ts=1628088960000"]}`)
// 	is.Equal(event.Retry, 0)
// 	ps.Publish("frontend:update:view/index.svelte", nil)
// 	event, err = hotClient.Next(ctx)
// 	is.NoErr(err)
// 	is.Equal(event.ID, "")
// 	is.Equal(event.Type, "")
// 	is.Equal(string(event.Data), `{"scripts":["/bud/view/index.svelte?ts=1628088960000"]}`)
// 	is.Equal(event.Retry, 0)
// 	ps.Publish("frontend:update", nil)
// 	event, err = hotClient.Next(ctx)
// 	is.NoErr(err)
// 	is.Equal(event.ID, "")
// 	is.Equal(event.Type, "")
// 	is.Equal(string(event.Data), `{"scripts":["/bud/view/index.svelte?ts=1628088960000"]}`)
// 	is.Equal(event.Retry, 0)
// 	is.NoErr(hotClient.Close())
// 	testServer.Close()
// }

// func TestReload(t *testing.T) {
// 	is := is.New(t)
// 	log := testlogs.New()
// 	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
// 	defer cancel()
// 	ps := pubsub.New()
// 	hotServer := hot.New(log, ps)
// 	hotServer.Now = func() time.Time { return now }
// 	testServer := httptest.NewServer(hotServer)
// 	hotClient, err := hot.Dial(log, testServer.URL+`/bud/hot/view/index.svelte`)
// 	is.NoErr(err)
// 	ps.Publish("backend:update", nil)
// 	event, err := hotClient.Next(ctx)
// 	is.NoErr(err)
// 	is.Equal(event.ID, "")
// 	is.Equal(event.Type, "")
// 	is.Equal(string(event.Data), `{"reload":true}`)
// 	is.Equal(event.Retry, 0)
// 	is.NoErr(hotClient.Close())
// 	testServer.Close()
// }

// // TODO: consolidate function. This is duplicated in multiple places.
// func listen(path string) (socket.Listener, *http.Client, error) {
// 	listener, err := socket.Listen(path)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	transport, err := socket.Transport(listener.Addr().String())
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	client := &http.Client{
// 		Timeout:   10 * time.Second,
// 		Transport: transport,
// 		CheckRedirect: func(req *http.Request, via []*http.Request) error {
// 			return http.ErrUseLastResponse
// 		},
// 	}
// 	return listener, client, nil
// }

// func TestUnixListener(t *testing.T) {
// 	is := is.New(t)
// 	log := testlogs.New()
// 	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
// 	defer cancel()
// 	listener, client, err := listen(filepath.Join(t.TempDir(), "test.sock"))
// 	is.NoErr(err)
// 	ps := pubsub.New()
// 	hotServer := hot.New(log, ps)
// 	hotServer.Now = func() time.Time { return now }
// 	server := &http.Server{
// 		Addr:    listener.Addr().String(),
// 		Handler: hotServer,
// 	}
// 	defer server.Shutdown(ctx)
// 	eg := new(errgroup.Group)
// 	eg.Go(func() error { return server.Serve(listener) })
// 	hotClient, err := hot.DialWith(client, log, "http://host/")
// 	is.NoErr(err)
// 	ps.Publish("frontend:update", nil)
// 	event, err := hotClient.Next(ctx)
// 	is.NoErr(err)
// 	is.Equal(event.ID, "")
// 	is.Equal(event.Type, "")
// 	is.Equal(string(event.Data), `{"reload":true}`)
// 	is.Equal(event.Retry, 0)
// 	is.NoErr(hotClient.Close())
// 	is.NoErr(server.Shutdown(ctx))
// }

// func TestNoWait(t *testing.T) {
// 	is := is.New(t)
// 	log := testlogs.New()
// 	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
// 	defer cancel()
// 	listener, client, err := listen(filepath.Join(t.TempDir(), "test.sock"))
// 	is.NoErr(err)
// 	ps := pubsub.New()
// 	hotServer := hot.New(log, ps)
// 	hotServer.Now = func() time.Time { return now }
// 	server := &http.Server{
// 		Addr:    listener.Addr().String(),
// 		Handler: hotServer,
// 	}
// 	defer server.Shutdown(ctx)
// 	eg := new(errgroup.Group)
// 	eg.Go(func() error { return server.Serve(listener) })
// 	hotClient, err := hot.DialWith(client, log, "http://host/")
// 	is.NoErr(err)
// 	is.NoErr(hotClient.Close())
// }

// func TestDrainBeforeClose(t *testing.T) {
// 	is := is.New(t)
// 	log := testlogs.New()
// 	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
// 	defer cancel()
// 	listener, client, err := listen(filepath.Join(t.TempDir(), "test.sock"))
// 	is.NoErr(err)
// 	ps := pubsub.New()
// 	hotServer := hot.New(log, ps)
// 	hotServer.Now = func() time.Time { return now }
// 	server := &http.Server{
// 		Addr:    listener.Addr().String(),
// 		Handler: hotServer,
// 	}
// 	defer server.Shutdown(ctx)
// 	eg := new(errgroup.Group)
// 	eg.Go(func() error { return server.Serve(listener) })
// 	hotClient, err := hot.DialWith(client, log, "http://host/")
// 	is.NoErr(err)
// 	ps.Publish("frontend:update", nil)
// 	ps.Publish("frontend:update", nil)
// 	ps.Publish("frontend:update", nil)
// 	ps.Publish("frontend:update", nil)
// 	is.NoErr(hotClient.Close())
// }
