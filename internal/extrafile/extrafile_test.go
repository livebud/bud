package extrafile_test

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/livebud/bud/internal/extrafile"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/socket"
)

func TestNoFiles(t *testing.T) {
	is := is.New(t)
	files, env, err := extrafile.Prepare("APP", 0)
	is.NoErr(err)
	is.Equal(len(files), 0)
	is.Equal(len(env), 0)
}

func listen(addr string) (socket.Listener, *http.Client, error) {
	listener, err := socket.Listen(addr)
	if err != nil {
		return nil, nil, err
	}
	transport, err := socket.Transport(listener.Addr().String())
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
	return listener, client, nil
}

func TestUnixPassthrough(t *testing.T) {
	// Parent process
	parent := func(t testing.TB) {
		is := is.New(t)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		dir := t.TempDir()
		appSocket, appClient, err := listen(filepath.Join(dir, "app.sock"))
		is.NoErr(err)
		defer appSocket.Close()
		hotSocket, hotClient, err := listen(filepath.Join(dir, "hot.sock"))
		is.NoErr(err)
		defer hotSocket.Close()
		appFiles, appEnv, err := extrafile.Prepare("APP", 0, appSocket)
		is.NoErr(err)
		is.Equal(len(appFiles), 1)
		is.Equal(len(appEnv), 2)
		hotFiles, hotEnv, err := extrafile.Prepare("HOT", 1, hotSocket)
		is.NoErr(err)
		is.Equal(len(hotFiles), 1)
		is.Equal(len(hotEnv), 2)
		// Ignore -test.count otherwise this will continue recursively
		var args []string
		for _, arg := range os.Args[1:] {
			if strings.HasPrefix(arg, "-test.count=") {
				continue
			}
			args = append(args, arg)
		}
		cmd := exec.CommandContext(ctx, os.Args[0], append(args, "-test.v=true", "-test.run=^"+t.Name()+"$")...)
		listener, err := socket.Listen(":0")
		is.NoErr(err)
		is.Equal(listener.Addr().Network(), "tcp")
		is.True(strings.HasPrefix(listener.Addr().String(), "127.0.0.1:"))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.ExtraFiles = append(cmd.ExtraFiles, appFiles...)
		cmd.ExtraFiles = append(cmd.ExtraFiles, hotFiles...)
		cmd.Env = append(os.Environ(), "CHILD=1")
		cmd.Env = append(cmd.Env, appEnv...)
		cmd.Env = append(cmd.Env, hotEnv...)
		is.NoErr(cmd.Start())

		// Test app socket
		res, err := appClient.Get("http://unix/ping")
		is.NoErr(err)
		is.Equal(res.StatusCode, 200)
		body, err := io.ReadAll(res.Body)
		is.NoErr(err)
		is.Equal(string(body), "app pong")
		res, err = appClient.Get("http://unix/close")
		is.NoErr(err)
		is.Equal(res.StatusCode, 200)

		// Test hot socket
		res, err = hotClient.Get("http://unix/ping")
		is.NoErr(err)
		is.Equal(res.StatusCode, 200)
		body, err = io.ReadAll(res.Body)
		is.NoErr(err)
		is.Equal(string(body), "hot pong")
		res, err = hotClient.Get("http://unix/close")
		is.NoErr(err)
		is.Equal(res.StatusCode, 200)

		is.NoErr(cmd.Wait())
	}

	// Child process
	child := func(t testing.TB) {
		is := is.New(t)
		appFiles := extrafile.Load("APP")
		is.Equal(len(appFiles), 1)
		hotFiles := extrafile.Load("HOT")
		is.Equal(len(hotFiles), 1)

		appListener, err := socket.From(appFiles[0])
		is.NoErr(err)
		hotListener, err := socket.From(hotFiles[0])
		is.NoErr(err)

		// Manually flush response so we can shutdown the server in the handler
		flush := func(w http.ResponseWriter) {
			if flusher, ok := w.(http.Flusher); ok {
				flusher.Flush()
			}
		}

		// Serve app
		appServer := &http.Server{
			Addr: appListener.Addr().String(),
		}
		appServer.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/ping":
				w.Write([]byte("app pong"))
			case "/close":
				flush(w)
				is.NoErr(appServer.Shutdown(context.Background()))
			default:
				w.WriteHeader(404)
			}
		})

		// Serve hot
		hotServer := &http.Server{
			Addr: hotListener.Addr().String(),
		}
		hotServer.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/ping":
				w.Write([]byte("hot pong"))
			case "/close":
				flush(w)
				is.NoErr(hotServer.Shutdown(context.Background()))
			default:
				w.WriteHeader(404)
			}
		})

		serve := func(server *http.Server, listener net.Listener) error {
			if err := server.Serve(listener); err != nil {
				if errors.Is(err, http.ErrServerClosed) {
					return nil
				}
				return err
			}
			return nil
		}

		eg := new(errgroup.Group)
		eg.Go(func() error { return serve(appServer, appListener) })
		eg.Go(func() error { return serve(hotServer, hotListener) })
		is.NoErr(eg.Wait())
	}

	if value := os.Getenv("CHILD"); value != "" {
		child(t)
	} else {
		parent(t)
	}
}

func TestTCPPassthrough(t *testing.T) {
	// Parent process
	parent := func(t testing.TB) {
		is := is.New(t)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		appSocket, appClient, err := listen(":0")
		is.NoErr(err)
		defer appSocket.Close()
		hotSocket, hotClient, err := listen(":0")
		is.NoErr(err)
		defer hotSocket.Close()
		appFiles, appEnv, err := extrafile.Prepare("APP", 0, appSocket)
		is.NoErr(err)
		is.Equal(len(appFiles), 1)
		is.Equal(len(appEnv), 2)
		hotFiles, hotEnv, err := extrafile.Prepare("HOT", 1, hotSocket)
		is.NoErr(err)
		is.Equal(len(hotFiles), 1)
		is.Equal(len(hotEnv), 2)
		// Ignore -test.count otherwise this will continue recursively
		var args []string
		for _, arg := range os.Args[1:] {
			if strings.HasPrefix(arg, "-test.count=") {
				continue
			}
			args = append(args, arg)
		}
		cmd := exec.CommandContext(ctx, os.Args[0], append(args, "-test.v=true", "-test.run=^"+t.Name()+"$")...)
		listener, err := socket.Listen(":0")
		is.NoErr(err)
		is.Equal(listener.Addr().Network(), "tcp")
		is.True(strings.HasPrefix(listener.Addr().String(), "127.0.0.1:"))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.ExtraFiles = append(cmd.ExtraFiles, appFiles...)
		cmd.ExtraFiles = append(cmd.ExtraFiles, hotFiles...)
		cmd.Env = append(os.Environ(), "CHILD=1")
		cmd.Env = append(cmd.Env, appEnv...)
		cmd.Env = append(cmd.Env, hotEnv...)
		is.NoErr(cmd.Start())

		// Test app socket
		res, err := appClient.Get("http://" + appSocket.Addr().String() + "/ping")
		is.NoErr(err)
		is.Equal(res.StatusCode, 200)
		body, err := io.ReadAll(res.Body)
		is.NoErr(err)
		is.Equal(string(body), "app pong")
		res, err = appClient.Get("http://" + appSocket.Addr().String() + "/close")
		is.NoErr(err)
		is.Equal(res.StatusCode, 200)

		// Test hot socket
		res, err = hotClient.Get("http://" + hotSocket.Addr().String() + "/ping")
		is.NoErr(err)
		is.Equal(res.StatusCode, 200)
		body, err = io.ReadAll(res.Body)
		is.NoErr(err)
		is.Equal(string(body), "hot pong")
		res, err = hotClient.Get("http://" + hotSocket.Addr().String() + "/close")
		is.NoErr(err)
		is.Equal(res.StatusCode, 200)

		is.NoErr(cmd.Wait())
	}

	// Child process
	child := func(t testing.TB) {
		is := is.New(t)
		appFiles := extrafile.Load("APP")
		is.Equal(len(appFiles), 1)
		hotFiles := extrafile.Load("HOT")
		is.Equal(len(hotFiles), 1)

		appListener, err := socket.From(appFiles[0])
		is.NoErr(err)
		hotListener, err := socket.From(hotFiles[0])
		is.NoErr(err)

		// Manually flush response so we can shutdown the server in the handler
		flush := func(w http.ResponseWriter) {
			if flusher, ok := w.(http.Flusher); ok {
				flusher.Flush()
			}
		}

		// Serve app
		appServer := &http.Server{
			Addr: appListener.Addr().String(),
		}
		appServer.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/ping":
				w.Write([]byte("app pong"))
			case "/close":
				flush(w)
				is.NoErr(appServer.Shutdown(context.Background()))
			default:
				w.WriteHeader(404)
			}
		})

		// Serve hot
		hotServer := &http.Server{
			Addr: hotListener.Addr().String(),
		}
		hotServer.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/ping":
				w.Write([]byte("hot pong"))
			case "/close":
				flush(w)
				is.NoErr(hotServer.Shutdown(context.Background()))
			default:
				w.WriteHeader(404)
			}
		})

		serve := func(server *http.Server, listener net.Listener) error {
			if err := server.Serve(listener); err != nil {
				if errors.Is(err, http.ErrServerClosed) {
					return nil
				}
				return err
			}
			return nil
		}

		eg := new(errgroup.Group)
		eg.Go(func() error { return serve(appServer, appListener) })
		eg.Go(func() error { return serve(hotServer, hotListener) })
		is.NoErr(eg.Wait())
	}

	if value := os.Getenv("CHILD"); value != "" {
		child(t)
	} else {
		parent(t)
	}
}
