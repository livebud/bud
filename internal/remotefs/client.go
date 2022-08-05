package remotefs

import (
	"io"
	"io/fs"
	"net/rpc"

	"github.com/livebud/bud/internal/virtual"
)

func NewClient(conn io.ReadWriteCloser) *Client {
	return &Client{rpc.NewClient(conn)}
}

type Client struct {
	rpc *rpc.Client
}

var _ fs.FS = (*Client)(nil)
var _ fs.ReadDirFS = (*Client)(nil)

func (c *Client) Open(name string) (fs.File, error) {
	vfile := new(virtual.File)
	err := c.rpc.Call("remotefs.Open", name, vfile)
	return vfile, err
}

func (c *Client) ReadDir(name string) (des []fs.DirEntry, err error) {
	vdes := []*virtual.DirEntry{}
	err = c.rpc.Call("remotefs.ReadDir", name, &vdes)
	if err != nil {
		return nil, err
	}
	for _, vde := range vdes {
		des = append(des, vde)
	}
	return des, nil
}

func (c *Client) Close() error {
	return c.rpc.Close()
}
