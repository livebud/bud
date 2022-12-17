package remotefs_test

import (
	"context"
	"errors"
	"io/fs"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/testsub"
	"github.com/livebud/bud/package/remotefs"
	"github.com/livebud/bud/package/socket"
	"github.com/livebud/bud/package/vfs"
)

func listen(t testing.TB) (net.Listener, error) {
	socketPath := filepath.Join(t.TempDir(), t.Name()+".sock")
	return socket.Listen(socketPath)
}

func TestReadFile(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	server, err := listen(t)
	is.NoErr(err)
	defer server.Close()
	client, err := remotefs.Dial(ctx, server.Addr().String())
	is.NoErr(err)
	fsys := vfs.Map{
		"a.txt": []byte("a"),
	}
	go remotefs.Serve(fsys, server)
	data, err := fs.ReadFile(client, "a.txt")
	is.NoErr(err)
	is.Equal(data, []byte("a"))
}

func TestReadDir(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	server, err := listen(t)
	is.NoErr(err)
	defer server.Close()
	client, err := remotefs.Dial(ctx, server.Addr().String())
	is.NoErr(err)
	fsys := vfs.Map{
		"tailwind/tailwind.css": []byte("/** tailwind **/"),
	}
	go remotefs.Serve(fsys, server)
	des, err := fs.ReadDir(client, "tailwind")
	is.NoErr(err)
	is.Equal(len(des), 1)
}

func TestFS(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	server, err := listen(t)
	is.NoErr(err)
	defer server.Close()
	client, err := remotefs.Dial(ctx, server.Addr().String())
	is.NoErr(err)
	fsys := fstest.MapFS{
		"tailwind/tailwind.css": &fstest.MapFile{Data: []byte("/** tailwind **/")},
		"markdoc/markdoc.js":    &fstest.MapFile{Data: []byte("/** markdoc **/")},
		"main.go":               &fstest.MapFile{Data: []byte("/** main **/")},
	}
	go remotefs.Serve(fsys, server)
	is.NoErr(fstest.TestFS(client, "tailwind/tailwind.css", "markdoc/markdoc.js"))
}

func TestOS(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	server, err := listen(t)
	is.NoErr(err)
	defer server.Close()
	client, err := remotefs.Dial(ctx, server.Addr().String())
	is.NoErr(err)
	go remotefs.Serve(os.DirFS("."), server)
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
	server, err := listen(t)
	is.NoErr(err)
	defer server.Close()
	client, err := remotefs.Dial(ctx, server.Addr().String())
	is.NoErr(err)
	fsys := fstest.MapFS{}
	go remotefs.Serve(fsys, server)
	data, err := fs.ReadFile(client, "client.go")
	is.True(err != nil)
	is.True(errors.Is(err, fs.ErrNotExist))
	is.Equal(data, nil)
}

func TestOSNotExist(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	ctx := context.Background()
	server, err := listen(t)
	is.NoErr(err)
	defer server.Close()
	client, err := remotefs.Dial(ctx, server.Addr().String())
	is.NoErr(err)
	fsys := os.DirFS(dir)
	go remotefs.Serve(fsys, server)
	data, err := fs.ReadFile(client, "client.go")
	is.True(err != nil)
	is.True(errors.Is(err, fs.ErrNotExist))
	is.Equal(data, nil)
}

func TestCommand(t *testing.T) {
	is := is.New(t)
	parent := func(t testing.TB, cmd *exec.Cmd) {
		ctx := context.Background()
		command := remotefs.Command{
			Env:    cmd.Env,
			Stderr: os.Stderr,
			Stdout: os.Stdout,
		}
		processfs, err := command.Start(ctx, cmd.Path, cmd.Args[1:]...)
		is.NoErr(err)
		defer processfs.Close()
		code, err := fs.ReadFile(processfs, "a.txt")
		is.NoErr(err)
		is.Equal(string(code), "a")
		is.NoErr(processfs.Close())
	}
	child := func(t testing.TB) {
		ctx := context.Background()
		fsys := fstest.MapFS{
			"a.txt": &fstest.MapFile{Data: []byte("a")},
		}
		err := remotefs.ServeFrom(ctx, fsys, "BUD_REMOTEFS")
		is.NoErr(err)
	}
	testsub.Run(t, parent, child)
}
