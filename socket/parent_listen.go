package socket

import (
	"net"
)

func Listen(path string) (net.Listener, error) {
	return listen(path)
}
