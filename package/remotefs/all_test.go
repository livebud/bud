package remotefs_test

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"testing"
	"testing/fstest"

	"github.com/livebud/bud/package/vfs"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/remotefs"
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

func TestReadFile(t *testing.T) {
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
	go remotefs.Serve(fsys, c2)
	data, err := fs.ReadFile(client, "a.txt")
	is.NoErr(err)
	is.Equal(data, []byte("a"))
}

func TestReadDir(t *testing.T) {
	is := is.New(t)
	r1, w1 := io.Pipe()
	r2, w2 := io.Pipe()
	c1 := &conn{r1, w2}
	defer c1.Close()
	c2 := &conn{r2, w1}
	defer c2.Close()
	client := remotefs.NewClient(c1)
	fsys := vfs.Map{
		"tailwind/tailwind.css": []byte("/** tailwind **/"),
	}
	go remotefs.Serve(fsys, c2)
	des, err := fs.ReadDir(client, "tailwind")
	is.NoErr(err)
	is.Equal(len(des), 1)
}

func TestFS(t *testing.T) {
	is := is.New(t)
	r1, w1 := io.Pipe()
	r2, w2 := io.Pipe()
	c1 := &conn{r1, w2}
	defer c1.Close()
	c2 := &conn{r2, w1}
	defer c2.Close()
	client := remotefs.NewClient(c1)
	fsys := fstest.MapFS{
		"tailwind/tailwind.css": &fstest.MapFile{Data: []byte("/** tailwind **/")},
		"markdoc/markdoc.js":    &fstest.MapFile{Data: []byte("/** markdoc **/")},
		"main.go":               &fstest.MapFile{Data: []byte("/** main **/")},
	}
	go remotefs.Serve(fsys, c2)
	is.NoErr(fstest.TestFS(client, "tailwind/tailwind.css", "markdoc/markdoc.js"))
}

func TestOS(t *testing.T) {
	is := is.New(t)
	r1, w1 := io.Pipe()
	r2, w2 := io.Pipe()
	c1 := &conn{r1, w2}
	defer c1.Close()
	c2 := &conn{r2, w1}
	defer c2.Close()
	client := remotefs.NewClient(c1)
	go remotefs.Serve(os.DirFS("."), c2)
	stat, err := fs.Stat(client, ".")
	is.NoErr(err)
	is.True(stat.IsDir())
	dir, err := client.Open(".")
	is.NoErr(err)
	stat, err = dir.Stat()
	is.NoErr(err)
	is.Equal(stat.IsDir(), true)
	stat, err = fs.Stat(client, "client.go")
	is.NoErr(err)
	is.Equal(stat.IsDir(), false)
	file, err := client.Open("client.go")
	is.NoErr(err)
	stat, err = file.Stat()
	is.NoErr(err)
	is.Equal(stat.IsDir(), false)
}

func TestNotExist(t *testing.T) {
	is := is.New(t)
	r1, w1 := io.Pipe()
	r2, w2 := io.Pipe()
	c1 := &conn{r1, w2}
	defer c1.Close()
	c2 := &conn{r2, w1}
	defer c2.Close()
	client := remotefs.NewClient(c1)
	fsys := fstest.MapFS{}
	go remotefs.Serve(fsys, c2)
	data, err := fs.ReadFile(client, "client.go")
	is.True(err != nil)
	is.True(errors.Is(err, fs.ErrNotExist))
	is.Equal(data, nil)
}
