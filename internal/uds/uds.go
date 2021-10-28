package uds

import (
	"context"
	"net"
	"net/http"
)

func Listen(path string) (net.Listener, error) {
	return net.Listen("unix", path)
}

func Transport(path string) http.RoundTripper {
	dialer := new(net.Dialer)
	return &http.Transport{
		DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
			return dialer.DialContext(ctx, "unix", path)
		},
	}
}
