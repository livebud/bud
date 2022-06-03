package extrafile_test

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/livebud/bud/package/exe"
	"github.com/livebud/bud/package/js/v8client"
	"github.com/livebud/bud/package/js/v8server"
	"golang.org/x/sync/errgroup"

	"github.com/livebud/bud/internal/extrafile"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/socket"
)

func TestNoFiles(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	cmd := exe.Command(ctx, "echo")
	extrafile.Inject(&cmd.ExtraFiles, &cmd.Env, "APP")
	is.Equal(len(cmd.Env), 0)
	is.Equal(len(cmd.ExtraFiles), 0)
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
		// Ignore -test.count otherwise this will continue recursively
		var args []string
		for _, arg := range os.Args[1:] {
			if strings.HasPrefix(arg, "-test.count=") {
				continue
			}
			args = append(args, arg)
		}
		cmd := exe.Command(ctx, os.Args[0], append(args, "-test.v=true", "-test.run=^"+t.Name()+"$")...)
		listener, err := socket.Listen(":0")
		is.NoErr(err)
		is.Equal(listener.Addr().Network(), "tcp")
		is.True(strings.HasPrefix(listener.Addr().String(), "127.0.0.1:"))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = append(os.Environ(), "CHILD=1")
		appFile, err := appSocket.File()
		is.NoErr(err)
		hotFile, err := hotSocket.File()
		is.NoErr(err)
		extrafile.Inject(&cmd.ExtraFiles, &cmd.Env, "APP", appFile)
		extrafile.Inject(&cmd.ExtraFiles, &cmd.Env, "HOT", hotFile)
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
		// Ignore -test.count otherwise this will continue recursively
		var args []string
		for _, arg := range os.Args[1:] {
			if strings.HasPrefix(arg, "-test.count=") {
				continue
			}
			args = append(args, arg)
		}
		cmd := exe.Command(ctx, os.Args[0], append(args, "-test.v=true", "-test.run=^"+t.Name()+"$")...)
		listener, err := socket.Listen(":0")
		is.NoErr(err)
		is.Equal(listener.Addr().Network(), "tcp")
		is.True(strings.HasPrefix(listener.Addr().String(), "127.0.0.1:"))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = append(os.Environ(), "CHILD=1")
		appFile, err := appSocket.File()
		is.NoErr(err)
		hotFile, err := hotSocket.File()
		is.NoErr(err)
		extrafile.Inject(&cmd.ExtraFiles, &cmd.Env, "APP", appFile)
		extrafile.Inject(&cmd.ExtraFiles, &cmd.Env, "HOT", hotFile)
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

func TestV8Passthrough(t *testing.T) {
	// Parent process
	parent := func(t testing.TB) {
		is := is.New(t)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		dir := t.TempDir()
		v8server, err := v8server.Pipe()
		is.NoErr(err)
		defer v8server.Close()
		// Ignore -test.count otherwise this will continue recursively
		var args []string
		for _, arg := range os.Args[1:] {
			if strings.HasPrefix(arg, "-test.count=") {
				continue
			}
			args = append(args, arg)
		}
		cmd := exe.Command(ctx, os.Args[0], append(args, "-test.v=true", "-test.run=^"+t.Name()+"$")...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = dir
		cmd.Env = append(os.Environ(), "CHILD=1")
		extrafile.Inject(&cmd.ExtraFiles, &cmd.Env, "V8", v8server.Files()...)
		is.NoErr(cmd.Start())
		// Create the V8 server
		eg := new(errgroup.Group)
		eg.Go(v8server.Serve)
		// Wait for the command to finish
		is.NoErr(cmd.Wait())
		// Restart and ensure the V8 server is ready to serve clients
		is.NoErr(cmd.Restart(ctx))
		is.NoErr(cmd.Wait())
		// Close the V8 server
		is.NoErr(v8server.Close())
		is.NoErr(eg.Wait())
	}

	// Child process
	child := func(t testing.TB) {
		is := is.New(t)
		client, err := v8client.From("V8")
		is.NoErr(err)
		result, err := client.Eval("eval.js", "2+2")
		is.NoErr(err)
		is.Equal(result, "4")
		// Test that console.log doesn't mess things up
		result, err = client.Eval("eval.js", "console.log('hi')")
		is.NoErr(err)
		is.Equal(result, "undefined")
	}

	if value := os.Getenv("CHILD"); value != "" {
		child(t)
	} else {
		parent(t)
	}
}
