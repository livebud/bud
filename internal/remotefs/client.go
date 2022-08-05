package remotefs

import (
	"io"
	"io/fs"
	"net/rpc"

	"github.com/livebud/bud/internal/virtual"
)

const name = "remotefs"
const openMethod = name + ".Open"

func NewClient(conn io.ReadWriteCloser) *Client {
	return &Client{rpc.NewClient(conn)}
}

type Client struct {
	rpc *rpc.Client
}

var _ fs.FS = (*Client)(nil)

func (c *Client) Open(name string) (fs.File, error) {
	vfile := new(virtual.File)
	err := c.rpc.Call(openMethod, name, vfile)
	return vfile, err
}

func (c *Client) Close() error {
	return c.rpc.Close()
}
