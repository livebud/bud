package run_test

import (
	"bytes"
	"context"
	"io"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/livebud/bud/internal/testdir"
	"github.com/livebud/bud/package/bud"
	"github.com/livebud/bud/package/bud/run"
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

func TestWelcome(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New()
	err := td.Write(dir)
	is.NoErr(err)
	listener, client, err := listenUnix(t.TempDir())
	is.NoErr(err)
	defer func() { is.NoErr(listener.Close()) }()
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	cmd := &run.Command{
		Dir:      dir,
		Listener: listener,
		Stdout:   stdout,
		Stderr:   stderr,
		Env: bud.Env{
			"TMPDIR":   t.TempDir(),
			"NO_COLOR": "1",
		},
	}
	process, err := cmd.Start(ctx)
	is.NoErr(err)
	defer func() { is.NoErr(process.Close()) }()
	is.Equal(stdout.String(), "")
	is.Equal(stderr.String(), "")
	res, err := client.Get("http://host/")
	is.NoErr(err)
	defer res.Body.Close()
	is.Equal(res.StatusCode, 200)
	body, err := io.ReadAll(res.Body)
	is.NoErr(err)
	is.True(strings.Contains(string(body), "Hey Bud"))
}

func TestController(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
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
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	cmd := &run.Command{
		Dir:      dir,
		Listener: listener,
		Stdout:   stdout,
		Stderr:   stderr,
		Env: bud.Env{
			"TMPDIR": t.TempDir(),
		},
	}
	process, err := cmd.Start(ctx)
	is.NoErr(err)
	defer func() { is.NoErr(process.Close()) }()
	is.Equal(stdout.String(), "")
	is.Equal(stderr.String(), "")
	res, err := client.Get("http://host/")
	is.NoErr(err)
	is.Equal(res.StatusCode, 200)
	body, err := io.ReadAll(res.Body)
	is.NoErr(err)
	is.True(strings.Contains(string(body), "from index"))
}
