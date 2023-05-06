package https

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/livebud/bud/internal/extrafile"
	"github.com/livebud/bud/internal/sig"
	"github.com/livebud/bud/package/socket"
)

type Service interface {
	Serve(ctx context.Context, ln net.Listener) error
}

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
	ln, err := socket.Listen(path)
	if err != nil {
		return nil, err
	}
	return ln, nil
}

type Server http.Server

var _ http.Handler = (*Server)(nil)

// Serve the server
func (s *Server) Serve(ctx context.Context, ln net.Listener) error {
	// Reference the original server
	server := (*http.Server)(s)
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

// ServeHTTP implements http.Handler
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	(*http.Server)(s).Handler.ServeHTTP(w, r)
}

// Serve the handler at address
func Serve(ctx context.Context, ln net.Listener, handler http.Handler) error {
	// Create the HTTP server
	server := &Server{
		Addr:    ln.Addr().String(),
		Handler: handler,
	}
	return server.Serve(ctx, ln)
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
