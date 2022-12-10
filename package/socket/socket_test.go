package socket_test

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/livebud/bud/internal/urlx"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/socket"
)

func TestLoadTCP(t *testing.T) {
	is := is.New(t)
	listener, err := socket.Listen(":0")
	is.NoErr(err)
	server := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(r.URL.Path))
		}),
	}
	go server.Serve(listener)
	transport, err := socket.Transport(listener.Addr().String())
	is.NoErr(err)
	client := &http.Client{
		Transport: transport,
		Timeout:   time.Second,
	}
	res, err := client.Get("http://" + listener.Addr().String() + "/hello")
	is.NoErr(err)
	body, err := io.ReadAll(res.Body)
	is.NoErr(err)
	is.Equal(string(body), "/hello")
	server.Shutdown(context.Background())
}

func TestLoadNumberOnly(t *testing.T) {
	is := is.New(t)
	listener, err := socket.Listen("0")
	is.NoErr(err)
	server := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(r.URL.Path))
		}),
	}
	go server.Serve(listener)
	transport, err := socket.Transport(listener.Addr().String())
	is.NoErr(err)
	client := &http.Client{
		Transport: transport,
		Timeout:   time.Second,
	}
	res, err := client.Get("http://" + listener.Addr().String() + "/hello")
	is.NoErr(err)
	body, err := io.ReadAll(res.Body)
	is.NoErr(err)
	is.Equal(string(body), "/hello")
	server.Shutdown(context.Background())
}

// This test is used to determine what the maximum socket length is.
// It should always fail.
func TestSocketLength(t *testing.T) {
	t.SkipNow()
	tmpDir := t.TempDir()
	for i := 1; i < 1000; i++ {
		socketPath := filepath.Join(tmpDir, strings.Repeat("a", i)+".sock")
		listener, err := socket.Listen(socketPath)
		if err != nil {
			t.Fatalf("failed at %d: %s", len(socketPath), err)
		}
		if err := listener.Close(); err != nil {
			t.Fatalf("unable to close listener: %s", err)
		}
	}
}

func TestDial(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	listener, err := socket.Listen(":0")
	is.NoErr(err)
	defer listener.Close()
	msg := "hello world"
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			incoming := make([]byte, len(msg))
			if _, err := io.ReadFull(conn, incoming); err != nil {
				conn.Close()
				return
			}
			conn.Write([]byte(string(incoming)))
			conn.Write([]byte(string(incoming)))
			conn.Close()
		}
	}()
	conn, err := socket.Dial(ctx, listener.Addr().String())
	is.NoErr(err)
	defer conn.Close()
	conn.Write([]byte(msg))
	outgoing := make([]byte, len(msg)*2)
	_, err = io.ReadFull(conn, outgoing)
	is.NoErr(err)
	is.Equal(string(outgoing), msg+msg)
}

func TestUDSCleanup(t *testing.T) {
	is := is.New(t)
	listener, err := socket.Listen("./test.sock")
	is.NoErr(err)
	defer listener.Close()
	is.NoErr(listener.Close())
	stat, err := os.Stat("test.sock")
	is.True(errors.Is(err, os.ErrNotExist))
	is.Equal(stat, nil)
}

func TestListenUp(t *testing.T) {
	is := is.New(t)
	// Letting the OS decide leads to port conflicts when run with other tests
	ln0, err := socket.Listen(":10000")
	is.NoErr(err)
	defer ln0.Close()
	ln1, err := socket.ListenUp(ln0.Addr().String(), 5)
	is.NoErr(err)
	defer ln1.Close()
	priorURL, err := urlx.Parse(ln0.Addr().String())
	is.NoErr(err)
	priorPort, err := strconv.Atoi(priorURL.Port())
	is.NoErr(err)
	url, err := urlx.Parse(ln1.Addr().String())
	is.NoErr(err)
	port, err := strconv.Atoi(url.Port())
	is.NoErr(err)
	is.Equal(port, priorPort+1)
}

func TestListenMaxAttemptsReached(t *testing.T) {
	is := is.New(t)
	// Letting the OS decide leads to port conflicts when run with other tests
	ln0, err := socket.Listen(":20000")
	is.NoErr(err)
	defer ln0.Close()
	// This one should work
	ln1, err := socket.ListenUp(ln0.Addr().String(), 1)
	is.NoErr(err)
	defer ln1.Close()
	// This one should fail because we're using ln0 as the base
	ln2, err := socket.ListenUp(ln0.Addr().String(), 1)
	is.True(errors.Is(err, socket.ErrAddrInUse))
	is.Equal(ln2, nil)
}

func TestListenPortTooHigh(t *testing.T) {
	is := is.New(t)
	ln0, err := socket.Listen(":65536")
	ae, ok := err.(*net.AddrError)
	is.True(ok)
	is.Equal(ae.Err, "invalid port")
	is.Equal(ln0, nil)
}
