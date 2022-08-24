package remotefs

import (
	"context"
	"encoding/gob"
	"io/fs"
	"strings"
	"time"

	"github.com/keegancsmith/rpc"
	"github.com/livebud/bud/internal/virtual"
	"github.com/livebud/bud/package/socket"
)

func init() {
	gob.Register(&virtual.File{})
	gob.Register(&virtual.Dir{})
	gob.Register(&virtual.DirEntry{})
}

// client timeout defaults to 10 seconds
const clientTimeout = 10 * time.Second

func Dial(ctx context.Context, addr string) (*Client, error) {
	conn, err := socket.Dial(ctx, addr)
	if err != nil {
		return nil, err
	}
	return &Client{rpc.NewClient(conn), clientTimeout}, nil
}

type Client struct {
	rpc *rpc.Client
	// Timeout for RPC calls
	Timeout time.Duration
}

var _ fs.FS = (*Client)(nil)
var _ fs.ReadDirFS = (*Client)(nil)

func (c *Client) Open(name string) (fs.File, error) {
	vfile := new(fs.File)
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()
	if err := c.rpc.Call(ctx, "remotefs.Open", name, vfile); err != nil {
		if isNotExist(err) {
			return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
		}
		return nil, err
	}
	return *vfile, nil
}

func (c *Client) ReadDir(name string) (des []fs.DirEntry, err error) {
	vdes := new([]fs.DirEntry)
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()
	err = c.rpc.Call(ctx, "remotefs.ReadDir", name, &vdes)
	if err != nil {
		if isNotExist(err) {
			return nil, &fs.PathError{Op: "readdir", Path: name, Err: fs.ErrNotExist}
		}
		return nil, err
	}
	return *vdes, nil
}

func (c *Client) Close() error {
	return c.rpc.Close()
}

// isNotExist is needed because the error has been serialized and passed between
// processes so errors.Is(err, fs.ErrNotExist) no longer is true.
func isNotExist(err error) bool {
	return strings.HasSuffix(err.Error(), fs.ErrNotExist.Error())
}
