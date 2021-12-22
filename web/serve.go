package web

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"

	"gitlab.com/mnm/bud/internal/sig"
)

// ServeUnix serves the handler over TCP.
func ServeTCP(ctx context.Context, address string, handler http.Handler) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	return serve(ctx, address, handler, listener)
}

// ServeUnix serves the handler over Unix Domain Sockets.
func ServeUnix(ctx context.Context, path string, handler http.Handler) error {
	listener, err := net.Listen("unix", path)
	if err != nil {
		return err
	}
	return serve(ctx, path, handler, listener)
}

// Serve the handler on the listener
func serve(ctx context.Context, a string, h http.Handler, l net.Listener) error {
	// Create the HTTP server
	server := &http.Server{Addr: a, Handler: h}
	// Make the server shutdownable
	shutdown := shutdown(ctx, server)
	// Serve requests
	if err := server.Serve(l); err != nil {
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
		forceCtx, cancel := sig.Trap(context.Background(), os.Interrupt)
		defer cancel()
		if err := server.Shutdown(forceCtx); err != nil {
			shutdown <- err
		}
		close(shutdown)
	}()
	return shutdown
}
