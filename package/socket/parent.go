package socket

import (
	"context"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/livebud/bud/internal/urlx"
)

type file interface {
	File() (*os.File, error)
}

type Env string

func (e Env) Key() string {
	i := strings.Index(string(e), "=")
	return string(e)[0:i]
}

func (e Env) Value() string {
	i := strings.Index(string(e), "=")
	return string(e)[i+1:]
}

func Files(l net.Listener) (files []*os.File, env Env, err error) {
	filer, ok := l.(file)
	if !ok {
		return []*os.File{}, "", nil
	}
	file, err := filer.File()
	if err != nil {
		return nil, "", err
	}
	return []*os.File{file}, "LISTEN_FDS=1", nil
}

func listen(path string) (net.Listener, error) {
	url, err := urlx.Parse(path)
	if err != nil {
		return nil, err
	}
	// Empty host means the path is a unix domain socket
	if url.Host == "" {
		addr, err := net.ResolveUnixAddr("unix", path)
		if err != nil {
			return nil, err
		}
		return net.ListenUnix("unix", addr)
	}
	// Otherwise, we listen on a TCP port
	addr, err := net.ResolveTCPAddr("tcp", url.Host)
	if err != nil {
		return nil, err
	}
	return net.ListenTCP("tcp", addr)
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
