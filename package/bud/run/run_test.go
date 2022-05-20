package run_test

import (
	"context"
	"net"
	"net/http"
	"path/filepath"
	"testing"
	"time"

	"github.com/livebud/bud/internal/testdir"
	"github.com/livebud/bud/package/bud/run"
	"github.com/livebud/bud/package/log/console"
	"github.com/livebud/bud/package/socket"
	"github.com/matryer/is"
)

func listenUnix(dir string) (net.Listener, *http.Client, error) {
	socketPath := filepath.Join(dir, "unix.sock")
	transport, err := socket.Transport(socketPath)
	if err != nil {
		return nil, nil, err
	}
	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	listener, err := socket.Listen(socketPath)
	if err != nil {
		return nil, nil, err
	}
	return listener, client, nil
}

func TestEmpty(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := "_tmp"
	td := testdir.New()
	td.Files["controller/controller.go"] = `
		package controller
		type Controller struct {}
		func (c *Controller) Index() string { return "from index" }
	`
	err := td.Write(dir)
	is.NoErr(err)
	listener, client, err := listenUnix(t.TempDir())
	is.NoErr(err)
	defer func() { is.NoErr(listener.Close()) }()
	cmd := &run.Command{
		Dir:      dir,
		Listener: listener,
		Log:      console.Stderr,
	}
	process, err := cmd.Start(ctx)
	is.NoErr(err)
	defer func() { is.NoErr(process.Close()) }()
	res, err := client.Get("http://host/")
	is.NoErr(err)
	is.Equal(res.StatusCode, 200)
}
