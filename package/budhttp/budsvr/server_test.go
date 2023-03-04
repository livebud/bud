package budsvr_test

import (
	"context"
	"errors"
	"testing"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/pubsub"
	"github.com/livebud/bud/package/budhttp/budsvr"
	v8 "github.com/livebud/bud/package/js/v8"
	"github.com/livebud/bud/package/log/testlog"
	"github.com/livebud/bud/package/socket"
	"github.com/livebud/bud/package/virtual"
	"golang.org/x/sync/errgroup"
)

func TestServerClose(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	fsys := virtual.List{}
	vm, err := v8.Load()
	is.NoErr(err)
	bus := pubsub.New()
	budln, err := socket.Listen(":0")
	is.NoErr(err)
	defer budln.Close()
	flag := new(framework.Flag)
	server := budsvr.New(budln, bus, flag, fsys, log, vm)
	server.Start(context.Background())
	is.NoErr(server.Close())
	is.NoErr(server.Close())
	is.NoErr(budln.Close())
}

func TestServerCancel(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	fsys := virtual.List{}
	vm, err := v8.Load()
	is.NoErr(err)
	bus := pubsub.New()
	budln, err := socket.Listen(":0")
	is.NoErr(err)
	defer budln.Close()
	flag := new(framework.Flag)
	server := budsvr.New(budln, bus, flag, fsys, log, vm)
	ctx, cancel := context.WithCancel(context.Background())
	server.Start(ctx)
	cancel()
	is.NoErr(server.Wait())
	is.NoErr(budln.Close())
}

func TestServerListenCancel(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	fsys := virtual.List{}
	vm, err := v8.Load()
	is.NoErr(err)
	bus := pubsub.New()
	budln, err := socket.Listen(":0")
	is.NoErr(err)
	defer budln.Close()
	flag := new(framework.Flag)
	server := budsvr.New(budln, bus, flag, fsys, log, vm)
	eg := new(errgroup.Group)
	ctx, cancel := context.WithCancel(context.Background())
	called := 0
	eg.Go(func() error {
		called++
		return server.Listen(ctx)
	})
	cancel()
	is.NoErr(eg.Wait())
	is.Equal(called, 1)
	is.NoErr(budln.Close())
}

// Random test to double-check nuances of errgroup.Group
func TestErrGroup(t *testing.T) {
	is := is.New(t)
	eg := new(errgroup.Group)
	is.NoErr(eg.Wait())
	called := 0
	eg.Go(func() error {
		called++
		return nil
	})
	is.NoErr(eg.Wait())
	is.Equal(called, 1)
	eg, ctx := errgroup.WithContext(context.Background())
	testErr := errors.New("some error")
	eg.Go(func() error {
		called++
		return testErr
	})
	is.True(errors.Is(eg.Wait(), testErr))
	is.Equal(called, 2)
	select {
	case <-ctx.Done():
	default:
		t.Fatal("context should be canceled")
	}
	is.Equal(ctx.Err(), context.Canceled)
	is.True(errors.Is(eg.Wait(), testErr))
	eg.Go(func() error {
		called++
		return nil
	})
	is.True(errors.Is(eg.Wait(), testErr))
	is.Equal(called, 3)
}

func TestServerCloseNoServe(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	fsys := virtual.List{}
	vm, err := v8.Load()
	is.NoErr(err)
	bus := pubsub.New()
	budln, err := socket.Listen(":0")
	is.NoErr(err)
	defer budln.Close()
	flag := new(framework.Flag)
	server := budsvr.New(budln, bus, flag, fsys, log, vm)
	is.NoErr(server.Close())
	is.NoErr(server.Close())
}

func TestServerWaitNoServe(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	fsys := virtual.List{}
	vm, err := v8.Load()
	is.NoErr(err)
	bus := pubsub.New()
	budln, err := socket.Listen(":0")
	is.NoErr(err)
	defer budln.Close()
	flag := new(framework.Flag)
	server := budsvr.New(budln, bus, flag, fsys, log, vm)
	is.NoErr(server.Wait())
	is.NoErr(server.Wait())
}
