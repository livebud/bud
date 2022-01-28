package socket_test

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/socket"
)

func TestUnixPassthrough(t *testing.T) {
	parent := func(t testing.TB) {
		socketPath := filepath.Join(t.TempDir(), "tmp.sock")
		is := is.New(t)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		// Ignore -test.count otherwise this will continue recursively
		var args []string
		for _, arg := range os.Args[1:] {
			if strings.HasPrefix(arg, "-test.count=") {
				continue
			}
			args = append(args, arg)
		}
		cmd := exec.CommandContext(ctx, os.Args[0], append(args, "-test.v=true", "-test.run=^TestUnixPassthrough$")...)
		listener, err := socket.Listen(socketPath)
		is.NoErr(err)
		is.Equal(listener.Addr().Network(), "unix")
		is.True(strings.HasSuffix(listener.Addr().String(), "tmp.sock"))
		extras, env, err := socket.Files(listener)
		is.NoErr(err)
		is.Equal(len(extras), 1)
		is.Equal(string(env), "LISTEN_FDS=1")
		is.Equal(env.Key(), "LISTEN_FDS")
		is.Equal(env.Value(), "1")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.ExtraFiles = append(cmd.ExtraFiles, extras...)
		cmd.Env = append(os.Environ(), "CHILD=1", string(env))
		is.NoErr(cmd.Start())
		transport, err := socket.Transport(socketPath)
		is.NoErr(err)
		client := &http.Client{
			Transport: transport,
			Timeout:   time.Second,
		}
		res, err := client.Get("http://host/hello")
		is.NoErr(err)
		body, err := ioutil.ReadAll(res.Body)
		is.NoErr(err)
		is.Equal(string(body), "/hello")
		_, err = client.Get("http://host/hello")
		is.NoErr(err)
		is.NoErr(cmd.Wait())
	}

	child := func(t testing.TB) {
		is := is.New(t)
		// Note: this should timeout if we actually use this new socket path and not
		// the file that's been passed in
		listener, err := socket.Load(filepath.Join(t.TempDir(), "tmp.sock"))
		is.NoErr(err)
		count := 0
		server := &http.Server{}
		server.Addr = listener.Addr().String()
		server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch count {
			case 0:
				count++
				os.Stderr.Write([]byte("stderr"))
				os.Stdout.Write([]byte("stdout"))
				w.Write([]byte(r.URL.Path))
			case 1:
				go server.Shutdown(context.Background())
				return
			}
		})
		if err := server.Serve(listener); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				is.NoErr(err)
			}
		}
	}

	if value := os.Getenv("CHILD"); value != "" {
		child(t)
	} else {
		parent(t)
	}
}

func TestTCPPassthrough(t *testing.T) {
	parent := func(t testing.TB) {
		is := is.New(t)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		// Ignore -test.count otherwise this will continue recursively
		var args []string
		for _, arg := range os.Args[1:] {
			if strings.HasPrefix(arg, "-test.count=") {
				continue
			}
			args = append(args, arg)
		}
		cmd := exec.CommandContext(ctx, os.Args[0], append(args, "-test.v=true", "-test.run=TestUnixPassthrough")...)
		listener, err := socket.Listen(":0")
		is.NoErr(err)
		is.Equal(listener.Addr().Network(), "tcp")
		is.True(strings.HasPrefix(listener.Addr().String(), "[::]:"))
		extras, env, err := socket.Files(listener)
		is.NoErr(err)
		is.Equal(len(extras), 1)
		is.Equal(string(env), "LISTEN_FDS=1")
		is.Equal(env.Key(), "LISTEN_FDS")
		is.Equal(env.Value(), "1")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.ExtraFiles = append(cmd.ExtraFiles, extras...)
		cmd.Env = append(os.Environ(), "CHILD=1", string(env))
		is.NoErr(cmd.Start())
		transport, err := socket.Transport(":0")
		is.NoErr(err)
		client := &http.Client{
			Transport: transport,
			Timeout:   time.Second,
		}
		res, err := client.Get("http://" + listener.Addr().String() + "/hello")
		is.NoErr(err)
		body, err := ioutil.ReadAll(res.Body)
		is.NoErr(err)
		is.Equal(string(body), "/hello")
		res, err = client.Get("http://" + listener.Addr().String() + "/hello")
		is.NoErr(err)
		is.NoErr(cmd.Wait())
	}

	child := func(t testing.TB) {
		is := is.New(t)
		// Note: this should timeout if we actually use this new socket path and not
		// the file that's been passed in
		listener, err := socket.Load(":0")
		is.NoErr(err)
		count := 0
		server := &http.Server{}
		server.Addr = listener.Addr().String()
		server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch count {
			case 0:
				count++
				w.Write([]byte(r.URL.Path))
			case 1:
				go server.Shutdown(context.Background())
				return
			}
		})
		if err := server.Serve(listener); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				is.NoErr(err)
			}
		}
	}

	if value := os.Getenv("CHILD"); value != "" {
		child(t)
	} else {
		parent(t)
	}
}

func TestLoadTCP(t *testing.T) {
	is := is.New(t)
	listener, err := socket.Load(":0")
	is.NoErr(err)
	server := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(r.URL.Path))
		}),
	}
	go server.Serve(listener)
	transport, err := socket.Transport(":0")
	is.NoErr(err)
	client := &http.Client{
		Transport: transport,
		Timeout:   time.Second,
	}
	res, err := client.Get("http://" + listener.Addr().String() + "/hello")
	is.NoErr(err)
	body, err := ioutil.ReadAll(res.Body)
	is.NoErr(err)
	is.Equal(string(body), "/hello")
	server.Shutdown(context.Background())
}

func TestLoadNumberOnly(t *testing.T) {
	is := is.New(t)
	listener, err := socket.Load("0")
	is.NoErr(err)
	server := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(r.URL.Path))
		}),
	}
	go server.Serve(listener)
	transport, err := socket.Transport("0")
	is.NoErr(err)
	client := &http.Client{
		Transport: transport,
		Timeout:   time.Second,
	}
	res, err := client.Get("http://" + listener.Addr().String() + "/hello")
	is.NoErr(err)
	body, err := ioutil.ReadAll(res.Body)
	is.NoErr(err)
	is.Equal(string(body), "/hello")
	server.Shutdown(context.Background())
}
