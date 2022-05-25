package testcli

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/livebud/bud/internal/cli"
	"github.com/livebud/bud/internal/errs"
	"github.com/livebud/bud/internal/once"
	"github.com/livebud/bud/package/exe"
	"github.com/livebud/bud/package/socket"
)

func New(cli *cli.CLI) *TestCLI {
	cli.Env["NO_COLOR"] = "1"
	return &TestCLI{cli: cli}
}

type TestCLI struct {
	cli *cli.CLI

	// make the temporary directory once
	mkTemp once.String
}

func (c *TestCLI) stdio() (stdout, stderr *bytes.Buffer) {
	stdout = new(bytes.Buffer)
	stderr = new(bytes.Buffer)
	// Write to stdio as well so debugging doesn't become too confusing
	c.cli.Stdout = io.MultiWriter(stdout, os.Stdout)
	c.cli.Stderr = io.MultiWriter(stderr, os.Stderr)
	return stdout, stderr
}

// makeTemp creates a temporary directory exactly once per run
func (c *TestCLI) makeTemp() (dir string, err error) {
	return c.mkTemp.Do(func() (string, error) {
		tmpDir := filepath.Join(c.cli.Dir(), "bud", "tmp")
		if err := os.MkdirAll(tmpDir, 0755); err != nil {
			return "", err
		}
		return tmpDir, nil
	})
}

// Env sets an environment variable, overriding any existing key.
func (c *TestCLI) Env(key, value string) {
	c.cli.Env[key] = value
}

// Stdin adds a reader as standard input.
func (c *TestCLI) Stdin(stdin io.Reader) {
	c.cli.Stdin = stdin
}

// Run the CLI with the provided args
func (c *TestCLI) Run(ctx context.Context, args ...string) (stdout, stderr *bytes.Buffer, err error) {
	stdout, stderr = c.stdio()
	if err := c.cli.Run(ctx, args...); err != nil {
		return stdout, stderr, err
	}
	return stdout, stderr, nil
}

func listenUnix(socketPath string) (net.Listener, *http.Client, error) {
	transport, err := socket.Transport(socketPath)
	if err != nil {
		return nil, nil, err
	}
	client := &http.Client{
		Timeout:   5 * time.Second,
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

func injectListener(t testing.TB, cli *cli.CLI) (*http.Client, func()) {
	t.Helper()
	// Start listening on a unix domain socket
	socketPath := filepath.Join(t.TempDir(), "unix.sock")
	listener, client, err := listenUnix(socketPath)
	if err != nil {
		t.Fatalf("unable to listen on socket path %q: %s", socketPath, err)
	}
	// Pull the files and environment from the listener
	files, env, err := socket.Files(listener)
	if err != nil {
		t.Fatalf("unable to derive *os.File from net.Listener: %s", err)
	}
	// Inject into CLI
	cli.ExtraFiles = append(cli.ExtraFiles, files...)
	cli.Env[env.Key()] = env.Value()
	// Return the client and a way to shutdown the listener
	return client, func() {
		if err := listener.Close(); err != nil {
			t.Fatalf("unexpected error while closing listener: %s", err)
		}
	}
}

func (c *TestCLI) Start(ctx context.Context, args ...string) (app *App, stdout *bytes.Buffer, stderr *bytes.Buffer, err error) {
	// Create the temporary directory
	tmpDir, err := c.makeTemp()
	if err != nil {
		return nil, nil, nil, err
	}
	// Start listening on a unix domain socket
	socketPath := filepath.Join(tmpDir, "app.sock")
	listener, client, err := listenUnix(socketPath)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("unable to listen on socket path %q: %s", socketPath, err)
	}
	// Pull the files and environment from the listener
	files, env, err := socket.Files(listener)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("unable to derive *os.File from net.Listener: %s", err)
	}
	// Inject into CLI
	c.cli.ExtraFiles = append(c.cli.ExtraFiles, files...)
	c.cli.Env[env.Key()] = env.Value()
	// Attach to stdout and stderr
	stdout, stderr = c.stdio()
	// Start the process
	process, err := c.cli.Start(ctx, args...)
	if err != nil {
		return nil, stdout, stderr, err
	}
	// Create a process
	return &App{
		process: process,
		client:  client,
		// Add resources to be cleaned up in order once
		resources: []*resource{
			&resource{"app process", process.Close},
			&resource{"app socket", listener.Close},
		},
	}, stdout, stderr, nil
}

type resource struct {
	Name  string
	Close func() error
}

type App struct {
	process   *exe.Cmd
	client    *http.Client
	once      once.Error
	resources []*resource
}

// Wait for the app to finish. This isn't typically necessary
func (a *App) Wait() error {
	return a.process.Wait()
}

// Close cleans up resources exactly once
func (a *App) Close() (err error) {
	return a.once.Do(func() (err error) {
		for _, r := range a.resources {
			if e := r.Close(); e != nil {
				err = errs.Join(e)
			}
		}
		return err
	})
}

func getURL(path string) string {
	return "http://host" + path
}

func (a *App) Get(url string) (*Response, error) {
	res, err := a.client.Get(getURL(url))
	if err != nil {
		return nil, err
	}
	body, err := bufferBody(res)
	if err != nil {
		return nil, err
	}
	// Dump the response
	headers, err := bufferHeaders(res)
	if err != nil {
		return nil, err
	}
	// Verify content-length, then remove it to make tests less fragile
	if err := checkContentLength(res, body); err != nil {
		return nil, err
	}
	res.Header.Del("Content-Length")
	// Check date, then remove it to make tests repeatable
	if err := checkDate(res); err != nil {
		return nil, err
	}
	res.Header.Del("Date")
	return &Response{res, body, headers}, nil
}

// bufferBody allows the response body to be read multiple times
// https://gist.github.com/franchb/d38fd9271e225a105a26c6859df1ce9b
func bufferBody(res *http.Response) ([]byte, error) {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	res.Body.Close()
	return body, nil
}

func bufferHeaders(res *http.Response) ([]byte, error) {
	return httputil.DumpResponse(res, false)
}

func checkContentLength(res *http.Response, body []byte) error {
	cl := res.Header.Get("Content-Length")
	if cl == "" {
		return nil
	}
	clen, err := strconv.Atoi(cl)
	if err != nil {
		return err
	}
	if clen != len(body) {
		return fmt.Errorf("Content-Length (%d) doesn't match the body length (%d)", clen, len(body))
	}
	return nil
}

func checkDate(res *http.Response) error {
	date := res.Header.Get("Date")
	if date == "" {
		return nil
	}
	dt, err := time.Parse(time.RFC1123, date)
	if err != nil {
		return err
	}
	// Date should be within 1 minute. In reality, it should be almost instant
	elapsed := time.Now().Sub(dt)
	if elapsed > time.Minute {
		return fmt.Errorf("Date header is too old %s", elapsed)
	}
	return nil
}

type Response struct {
	res     *http.Response
	body    []byte
	headers []byte
}

// Status returns the response status
func (r *Response) Status() int {
	return r.res.StatusCode
}

func (r *Response) Headers() *bytes.Buffer {
	return bytes.NewBuffer(r.headers)
}

func (r *Response) Body() *bytes.Buffer {
	return bytes.NewBuffer(r.body)
}

// Diff the response the expected HTTP response
func (r *Response) Dump() *bytes.Buffer {
	return bytes.NewBuffer(bytes.Join(
		[][]byte{r.headers, r.body},
		[]byte{'\r', '\n'},
	))
}

// Header gets a value from a key
func (r *Response) Header(key string) string {
	return r.res.Header.Get(key)
}
