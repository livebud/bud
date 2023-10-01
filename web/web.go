package web

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/livebud/bud/internal/signals"
	"github.com/livebud/bud/internal/socket"
)

type Handler http.Handler

type Router interface {
	http.Handler
	Get(route string, handler http.Handler) error
	Post(route string, handler http.Handler) error
	Put(route string, handler http.Handler) error
	Patch(route string, handler http.Handler) error
	Delete(route string, handler http.Handler) error
	Set(method string, route string, handler http.Handler) error
}

type Server http.Server

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Handler.ServeHTTP(w, r)
}

// Service is a server that can be served
type Service interface {
	Serve(ctx context.Context, ln net.Listener) error
}

func (s *Server) Serve(ctx context.Context, ln net.Listener) error {
	return Serve(ctx, (*http.Server)(s), ln)
}

func (s *Server) ListenAndServe(ctx context.Context) error {
	ln, err := Listen(s.Addr)
	if err != nil {
		return err
	}
	return s.Serve(ctx, ln)
}

func Listen(address string) (net.Listener, error) {
	return socket.Listen(address)
}

// Serve the handler at address
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
		forceCtx := signals.Trap(context.Background(), os.Interrupt)
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
