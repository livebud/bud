package testcli

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/livebud/bud/internal/cli"
	"github.com/livebud/bud/internal/errs"
	"github.com/livebud/bud/internal/once"
	"github.com/livebud/bud/package/exe"
	"github.com/livebud/bud/package/hot"
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
		absPath, err := filepath.Abs(c.cli.Dir())
		if err != nil {
			return "", err
		}
		tmpDir := filepath.Join(absPath, "bud", "tmp")
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

func listen(path string) (socket.Listener, *http.Client, error) {
	listener, err := socket.Listen(path)
	if err != nil {
		return nil, nil, err
	}
	transport, err := socket.Transport(listener.Addr().String())
	if err != nil {
		return nil, nil, err
	}
	client := &http.Client{
		Timeout:   30 * time.Second,
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	return listener, client, nil
}

func (c *TestCLI) Start(ctx context.Context, args ...string) (app *App, stdout *bytes.Buffer, stderr *bytes.Buffer, err error) {
	// Create the temporary directory
	tmpDir, err := c.makeTemp()
	if err != nil {
		return nil, nil, nil, err
	}
	// Start listening on a unix domain socket for the app
	appSocketPath := filepath.Join(tmpDir, "app.sock")
	appListener, appClient, err := listen(appSocketPath)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("unable to listen on socket path %q: %s", appSocketPath, err)
	}
	if err := c.cli.Inject("APP", appListener); err != nil {
		return nil, nil, nil, err
	}
	// Start listening on a unix domain socket for the hot
	hotSocketPath := filepath.Join(tmpDir, "hot.sock")
	hotListener, hotClient, err := listen(hotSocketPath)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("unable to listen on socket path %q: %s", hotSocketPath, err)
	}
	if err := c.cli.Inject("HOT", hotListener); err != nil {
		return nil, nil, nil, err
	}
	// Attach to stdout and stderr
	stdout, stderr = c.stdio()
	// Start the process
	process, err := c.cli.Start(ctx, args...)
	if err != nil {
		return nil, stdout, stderr, err
	}
	// Create a process
	return &App{
		process:   process,
		appClient: appClient,
		hotClient: hotClient,
		// Add resources to be cleaned up in order once
		resources: []*resource{
			&resource{"app process", process.Close},
			&resource{"app socket", appListener.Close},
			&resource{"hot socket", hotListener.Close},
		},
	}, stdout, stderr, nil
}

type resource struct {
	Name  string
	Close func() error
}

type App struct {
	process   *exe.Cmd
	appClient *http.Client
	hotClient *http.Client
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

func (a *App) Get(path string) (*Response, error) {
	req, err := http.NewRequest(http.MethodGet, getURL(path), nil)
	if err != nil {
		return nil, err
	}
	return a.Request(req)
}

func (a *App) GetJSON(path string) (*Response, error) {
	req, err := http.NewRequest(http.MethodGet, getURL(path), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	return a.Request(req)
}

func (a *App) Post(path string, body io.Reader) (*Response, error) {
	req, err := http.NewRequest(http.MethodPost, getURL(path), body)
	if err != nil {
		return nil, err
	}
	return a.Request(req)
}

func (a *App) PostJSON(path string, body io.Reader) (*Response, error) {
	req, err := http.NewRequest(http.MethodPost, getURL(path), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return a.Request(req)
}

func (a *App) Patch(path string, body io.Reader) (*Response, error) {
	req, err := http.NewRequest(http.MethodPatch, getURL(path), body)
	if err != nil {
		return nil, err
	}
	return a.Request(req)
}

func (a *App) PatchJSON(path string, body io.Reader) (*Response, error) {
	req, err := http.NewRequest(http.MethodPatch, getURL(path), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return a.Request(req)
}

func (a *App) Delete(path string, body io.Reader) (*Response, error) {
	req, err := http.NewRequest(http.MethodDelete, getURL(path), body)
	if err != nil {
		return nil, err
	}
	return a.Request(req)
}

func (a *App) DeleteJSON(path string, body io.Reader) (*Response, error) {
	req, err := http.NewRequest(http.MethodDelete, getURL(path), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return a.Request(req)
}

func (a *App) Request(req *http.Request) (*Response, error) {
	res, err := a.appClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
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
	// Buffer the headers response
	headers, err := bufferHeaders(res, body)
	if err != nil {
		return nil, err
	}
	return &Response{res, headers, body}, nil
}

// Hot connects to the event stream
func (a *App) Hot(path string) (*hot.Stream, error) {
	return hot.DialWith(a.hotClient, getURL(path))
}

func bufferHeaders(res *http.Response, body []byte) ([]byte, error) {
	dump, err := httputil.DumpResponse(res, false)
	if err != nil {
		return nil, err
	}
	// httputil.DumpResponse() always attaches a Content-Length, regardless of
	// whether or not you remove it. This scanner removes the Content-Lenght
	// manually.
	s := bufio.NewScanner(bytes.NewBuffer(dump))
	b := new(bytes.Buffer)
	for s.Scan() {
		if bytes.Contains(s.Bytes(), []byte("Content-Length")) {
			continue
		}
		b.WriteByte('\n')
		b.Write(s.Bytes())
	}
	if s.Err() != nil {
		return nil, s.Err()
	}
	return b.Bytes(), nil
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
	headers []byte
	body    []byte
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
	b := new(bytes.Buffer)
	b.Write(r.headers)
	b.WriteByte('\n')
	b.Write(r.body)
	return b
}

// Header gets a value from a key
func (r *Response) Header(key string) string {
	return r.res.Header.Get(key)
}

// Query a selector on the page using goquery
func (r *Response) Query(selector string) (*goquery.Selection, error) {
	doc, err := goquery.NewDocumentFromReader(r.Body())
	if err != nil {
		return nil, err
	}
	return doc.Find(selector), nil
}
