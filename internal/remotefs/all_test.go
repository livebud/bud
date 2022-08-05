package remotefs_test

import (
	"io"
	"io/fs"
	"testing"

	"github.com/livebud/bud/package/vfs"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/remotefs"
)

type conn struct {
	rc io.ReadCloser
	wc io.WriteCloser
}

func (c *conn) Read(p []byte) (int, error) {
	return c.rc.Read(p)
}

func (c *conn) Write(p []byte) (int, error) {
	return c.wc.Write(p)
}

func (c *conn) Close() error {
	c.rc.Close()
	c.wc.Close()
	return nil
}

func TestClientServer(t *testing.T) {
	is := is.New(t)
	r1, w1 := io.Pipe()
	r2, w2 := io.Pipe()
	c1 := &conn{r1, w2}
	defer c1.Close()
	c2 := &conn{r2, w1}
	defer c2.Close()
	client := remotefs.NewClient(c1)
	fsys := vfs.Map{
		"a.txt": []byte("a"),
	}
	go remotefs.NewServer(fsys, c2)
	data, err := fs.ReadFile(client, "a.txt")
	is.NoErr(err)
	is.Equal(data, []byte("a"))
}
