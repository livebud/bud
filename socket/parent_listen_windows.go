package socket

import (
	"errors"
	"net"
)

// Listen creates a noop Listener because windows doesn't support
// net.FileListener which we need to pass an existing connection through to a
// new process.
func Listen(path string) (net.Listener, error) {
	return &noopListener{}, nil
}

type noopListener struct {
}

var _ net.Listener = (*noopListener)(nil)

func (nl *noopListener) Accept() (Conn, error) {
	return nil, errors.New("socket: accept not implemented")
}

func (nl *noopListener) Close() error {
	return nil, errors.New("socket: close not implemented")
}

func (nl *noopListener) Addr() net.Addr {
	return &noopAddr{}
}

type noopAddr struct{}

func (na *noopAddr) Network() string {
	return "noop"
}

func (na *noopAddr) String() string {
	return ""
}
