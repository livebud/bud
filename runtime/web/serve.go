package web

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/livebud/bud/internal/extrafile"
	"github.com/livebud/bud/internal/sig"
	"github.com/livebud/bud/package/log/console"
	"github.com/livebud/bud/package/socket"
)

// listen first tries pulling the connection from a passed in file descriptor.
// If that fails, it will start listening on a path.
func Listen(prefix, path string) (socket.Listener, error) {
	files := extrafile.Load(prefix)
	if len(files) > 0 {
		// Turn the passed in file descriptor into a listener
		return socket.From(files[0])
	}
	if path == "" {
		path = "localhost:3000"
	}
	// Listen on a path
	listener, err := socket.Listen(path)
	if err != nil {
		return nil, err
	}
	// This is done to avoid logging the live-reload server.
	// TODO: clean this up.
	if prefix == "APP" {
		// Log here because this is the first time we've bound to a resource.
		console.Info("Listening on " + Format(listener))
	}
	return listener, nil
}

// Serve the handler at address
func Serve(ctx context.Context, listener net.Listener, handler http.Handler) error {
	// Create the HTTP server
	server := &http.Server{Addr: listener.Addr().String(), Handler: handler}
	// Make the server shutdownable
	shutdown := shutdown(ctx, server)
	// Serve requests
	if err := server.Serve(listener); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}
	// Handle any errors that occurred while shutting down
	if err := <-shutdown; err != nil {
		return err
	}
	return nil
}

// Shutdown the server when the context is canceled
func shutdown(ctx context.Context, server *http.Server) <-chan error {
	shutdown := make(chan error, 1)
	go func() {
		<-ctx.Done()
		// Wait for one more interrupt to force an immediate shutdown
		forceCtx := sig.Trap(ctx, os.Interrupt)
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
