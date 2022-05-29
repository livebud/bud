package socket

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/livebud/bud/internal/urlx"
)

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

func (l listener) File() (*os.File, error) {
	filer, ok := l.Listener.(file)
	if !ok {
		return nil, fmt.Errorf("socket: %s is not a file", l.Listener.Addr().String())
	}
	return filer.File()
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
		return listener{unix}, nil
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
	return listener{tcp}, nil
}

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
