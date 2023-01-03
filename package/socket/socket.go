package socket

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/livebud/bud/internal/urlx"
)

// ErrAddrInUse occurs when a port is already in use
var ErrAddrInUse = syscall.EADDRINUSE

type Listener interface {
	net.Listener
	file
}

type listener struct {
	net.Listener
}

type file interface {
	File() (*os.File, error)
}

func (l *listener) File() (*os.File, error) {
	filer, ok := l.Listener.(file)
	if !ok {
		return nil, fmt.Errorf("socket: %s is not a file", l.Listener.Addr().String())
	}
	return filer.File()
}

func (l *listener) Close() error {
	if err := l.Listener.Close(); err != nil {
		// Ignore errors where the listener has been closed already. This can occur
		// when the server is shutdown before the listener is closed.
		if !errors.Is(err, net.ErrClosed) {
			return err
		}
	}
	return nil
}

// Listen on a path or port
func Listen(path string) (Listener, error) {
	url, err := urlx.Parse(path)
	if err != nil {
		return nil, err
	}
	// Empty host means the path is a unix domain socket
	if url.Host == "" {
		// Unix domain socket path can't be more than 103 characters long
		if len(path) > 103 {
			return nil, fmt.Errorf("socket: unix path too long %q", path)
		}
		addr, err := net.ResolveUnixAddr("unix", path)
		if err != nil {
			return nil, err
		}
		unix, err := net.ListenUnix("unix", addr)
		if err != nil {
			return nil, err
		}
		return &listener{unix}, nil
	}
	// Otherwise, we listen on a TCP port
	addr, err := net.ResolveTCPAddr("tcp", url.Host)
	if err != nil {
		return nil, err
	}
	tcp, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &listener{tcp}, nil
}

// ListenUp is similar to listen, but will increment the port number until it
// finds a free one or reaches the maximum number of attempts
func ListenUp(path string, attempts int) (Listener, error) {
	ln, err := Listen(path)
	if err != nil {
		if !errors.Is(err, ErrAddrInUse) {
			return nil, err
		}
		if attempts--; attempts >= 0 {
			newPath, err := incrementPort(path)
			if err != nil {
				return nil, err
			}
			return ListenUp(newPath, attempts-1)
		}
		return nil, err
	}
	return ln, nil
}

// Takes a address and increments the port by 1
func incrementPort(path string) (string, error) {
	url, err := urlx.Parse(path)
	if err != nil {
		return "", err
	}
	port, err := strconv.Atoi(url.Port())
	if err != nil {
		return "", err
	}
	port++
	url.Host = url.Hostname() + ":" + strconv.Itoa(port)
	return url.String(), nil
}

// Dial creates a connection to an address
func Dial(ctx context.Context, address string) (net.Conn, error) {
	url, err := urlx.Parse(address)
	if err != nil {
		return nil, err
	}
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	// Empty host means the path is a unix domain socket
	if url.Host == "" {
		return dialer.DialContext(ctx, "unix", address)
	}
	return dialer.DialContext(ctx, "tcp", url.Host)
}

// Transport creates a RoundTripper for an HTTP Client
func Transport(path string) (http.RoundTripper, error) {
	url, err := urlx.Parse(path)
	if err != nil {
		return nil, err
	}
	// Empty host means the path is a unix domain socket
	if url.Host == "" {
		dialer := new(net.Dialer)
		return &http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				return dialer.DialContext(ctx, "unix", path)
			},
		}, nil
	}
	return httpTransport(url.Host), nil
}

// httpTransport is a modified from http.DefaultTransport
func httpTransport(host string) http.RoundTripper {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.DialContext(ctx, network, host)
		},
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}

// From turns a file into a Listener or fails trying
func From(file *os.File) (Listener, error) {
	ln, err := net.FileListener(file)
	if err != nil {
		return nil, err
	}
	if err := file.Close(); err != nil {
		return nil, err
	}
	return &listener{ln}, nil
}
