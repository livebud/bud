package pubsub_test

import (
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/pubsub"
	"golang.org/x/sync/errgroup"
)

func TestPubSub(t *testing.T) {
	is := is.New(t)
	ps := pubsub.New()
	ps.Publish("toast", []byte("nothing to publish to yet"))
	sub := ps.Subscribe("toast")
	eg := new(errgroup.Group)
	eg.Go(func() error {
		msg := <-sub.Wait()
		is.Equal(string(msg), "toast is ready")
		return nil
	})
	ps.Publish("toast", []byte("toast is ready"))
	is.NoErr(eg.Wait())
	sub.Close()
}

func TestCloseTwice(t *testing.T) {
	ps := pubsub.New()
	sub := ps.Subscribe("toast")
	sub.Close()
	sub.Close()
}

func TestSubTwice(t *testing.T) {
	ps := pubsub.New()
	sub := ps.Subscribe("toast")
	ps.Publish("toast", nil)
	<-sub.Wait()
	sub.Close()
	sub = ps.Subscribe("toast")
	select {
	case <-sub.Wait():
		t.Fatal("lingering event")
	default:
	}
	sub.Close()
}
