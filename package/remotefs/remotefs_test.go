package remotefs_test

import (
	"context"
	"errors"
	"io/fs"
	"net"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/livebud/bud/package/socket"
	"github.com/livebud/bud/package/vfs"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/remotefs"
)

func listener(t *testing.T) (net.Listener, error) {
	socketPath := filepath.Join(t.TempDir(), t.Name()+".sock")
	return socket.Listen(socketPath)
}

func TestReadFile(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	ln, err := listener(t)
	is.NoErr(err)
	defer ln.Close()
	client, err := remotefs.Dial(ctx, ln.Addr().String())
	is.NoErr(err)
	fsys := vfs.Map{
		"a.txt": []byte("a"),
	}
	go remotefs.Serve(fsys, ln)
	data, err := fs.ReadFile(client, "a.txt")
	is.NoErr(err)
	is.Equal(data, []byte("a"))
}

func TestReadDir(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	ln, err := listener(t)
	is.NoErr(err)
	defer ln.Close()
	client, err := remotefs.Dial(ctx, ln.Addr().String())
	is.NoErr(err)
	fsys := vfs.Map{
		"tailwind/tailwind.css": []byte("/** tailwind **/"),
	}
	go remotefs.Serve(fsys, ln)
	des, err := fs.ReadDir(client, "tailwind")
	is.NoErr(err)
	is.Equal(len(des), 1)
}

func TestFS(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	ln, err := listener(t)
	is.NoErr(err)
	defer ln.Close()
	client, err := remotefs.Dial(ctx, ln.Addr().String())
	is.NoErr(err)
	fsys := fstest.MapFS{
		"tailwind/tailwind.css": &fstest.MapFile{Data: []byte("/** tailwind **/")},
		"markdoc/markdoc.js":    &fstest.MapFile{Data: []byte("/** markdoc **/")},
		"main.go":               &fstest.MapFile{Data: []byte("/** main **/")},
	}
	go remotefs.Serve(fsys, ln)
	is.NoErr(fstest.TestFS(client, "tailwind/tailwind.css", "markdoc/markdoc.js"))
}

func TestOS(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	ln, err := listener(t)
	is.NoErr(err)
	defer ln.Close()
	client, err := remotefs.Dial(ctx, ln.Addr().String())
	is.NoErr(err)
	go remotefs.Serve(os.DirFS("."), ln)
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
	ctx := context.Background()
	ln, err := listener(t)
	is.NoErr(err)
	defer ln.Close()
	client, err := remotefs.Dial(ctx, ln.Addr().String())
	is.NoErr(err)
	fsys := fstest.MapFS{}
	go remotefs.Serve(fsys, ln)
	data, err := fs.ReadFile(client, "client.go")
	is.True(err != nil)
	is.True(errors.Is(err, fs.ErrNotExist))
	is.Equal(data, nil)
}
