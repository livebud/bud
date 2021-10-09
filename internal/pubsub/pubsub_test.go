package pubsub_test

import (
	"testing"

	"github.com/go-duo/bud/internal/pubsub"
	"github.com/matryer/is"
	"golang.org/x/sync/errgroup"
)

func Test(t *testing.T) {
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
