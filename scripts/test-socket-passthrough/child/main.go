package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/sync/errgroup"

	"github.com/livebud/bud/internal/sig"

	"github.com/livebud/bud/package/socket"

	"github.com/livebud/bud/internal/extrafile"
)

func main() {
	if err := run(sig.Trap(context.Background(), os.Interrupt)); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	files := extrafile.Load("APP")
	if len(files) == 0 {
		return fmt.Errorf("no files passed through")
	}
	listener, err := socket.From(files[0])
	if err != nil {
		return err
	}
	server := &http.Server{
		Addr: ":0",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("hello world"))
		}),
	}
	eg := new(errgroup.Group)
	eg.Go(func() error {
		fmt.Println("listening on", listener.Addr())
		return server.Serve(listener)
	})
	<-ctx.Done()
	if err := server.Shutdown(context.Background()); err != nil {
		return err
	}
	if err := eg.Wait(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
	return nil
}
