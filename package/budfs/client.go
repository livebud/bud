package budfs

import (
	"context"
	"net/rpc"

	"github.com/livebud/bud/package/socket"
)

func Dial(ctx context.Context, address string) (*Proxy, error) {
	conn, err := socket.Dial(ctx, addr)
	if err != nil {
		return nil, err
	}
	return &Proxy{rpc.NewClient(conn)}, nil
}

type Proxy struct {
	rpc *rpc.Client
}
