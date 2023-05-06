package web

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/livebud/bud/internal/sig"
)

type Server interface {
	Serve(ctx context.Context, ln net.Listener) error
}

// Serve the server
func Serve(ctx context.Context, server *http.Server, ln net.Listener) error {
	// Make the server shutdownable
	shutdown := shutdown(ctx, server)
	// Serve requests
	if err := server.Serve(ln); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}
	// Handle any errors that occurred while shutting down
	if err := <-shutdown; err != nil {
		if !errors.Is(err, context.Canceled) {
			return err
		}
	}
	return nil
}

// Shutdown the server when the context is canceled
func shutdown(ctx context.Context, server *http.Server) <-chan error {
	shutdown := make(chan error, 1)
	go func() {
		<-ctx.Done()
		// Wait for one more interrupt to force an immediate shutdown
		forceCtx := sig.Trap(context.Background(), os.Interrupt)
		if err := server.Shutdown(forceCtx); err != nil {
			shutdown <- err
		}
		close(shutdown)
	}()
	return shutdown
}

// Format a listener
func Format(l net.Listener) string {
	address := l.Addr().String()
	if l.Addr().Network() == "unix" {
		return address
	}
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		// Give up trying to format.
		// TODO: figure out if this can occur.
		return address
	}
	// https://serverfault.com/a/444557
	if host == "::" {
		host = "0.0.0.0"
	}
	return fmt.Sprintf("http://%s:%s", host, port)
}
