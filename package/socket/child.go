package socket

import (
	"net"
	"os"
)

// From turns a file into a Listener or fails trying
func From(file *os.File) (Listener, error) {
	ln, err := net.FileListener(file)
	if err != nil {
		return nil, err
	}
	if err := file.Close(); err != nil {
		return nil, err
	}
	return listener{ln}, nil
}
