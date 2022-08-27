package testcli

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/lithammer/dedent"
	"github.com/livebud/bud/internal/cli/bud"
	"github.com/livebud/bud/internal/once"
	"github.com/livebud/bud/internal/pubsub"
	"github.com/matthewmueller/diff"

	"github.com/livebud/bud/package/hot"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/log/testlog"
	"github.com/livebud/bud/package/socket"

	"golang.org/x/sync/errgroup"

	"github.com/livebud/bud/internal/cli"
	"github.com/livebud/bud/internal/envs"
)

func New(dir string) *CLI {
	ps := pubsub.New()
	return &CLI{
		dir: dir,
		bus: ps,
		Env: envs.Map{
			"NO_COLOR": "1",
			"HOME":     os.Getenv("HOME"),
			"PATH":     os.Getenv("PATH"),
			"GOPATH":   os.Getenv("GOPATH"),
			"TMPDIR":   os.TempDir(),
		},
		Stdin: nil,
	}
}

type CLI struct {
	dir   string
	bus   pubsub.Client
	Env   envs.Map
	Stdin io.Reader
}

// Flags that can be set from the test suite
// These can be overridden by more specific flags
func prependFlags(args []string) []string {
	return append([]string{
		"--log=" + testlog.Pattern(),
	}, args...)
}

func (c *CLI) Run(ctx context.Context, args ...string) (*Result, error) {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	cli := cli.New(&bud.Input{
		Dir:    c.dir,
		Bus:    c.bus,
		Env:    c.Env.List(),
		Stdin:  c.Stdin,
		Stdout: io.MultiWriter(os.Stdout, stdout),
		Stderr: io.MultiWriter(os.Stderr, stderr),
	})
	err := cli.Run(ctx, prependFlags(args)...)
	return &Result{stdout, stderr}, err
}

type Result struct {
	stdout *bytes.Buffer
	stderr *bytes.Buffer
}

func (r *Result) Stdout() string {
	return r.stdout.String()
}

func (r *Result) Stderr() string {
	return r.stderr.String()
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
		Timeout:   60 * time.Second,
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	return listener, client, nil
}

func (c *CLI) Start(ctx context.Context, args ...string) (*Client, error) {
	log := testlog.New()
	// TODO: listen unix and create client
	webln, webc, err := listen(":0")
	if err != nil {
		return nil, err
	}
	// TODO: listen unix and create client
	budln, budc, err := listen(":0")
	if err != nil {
		return nil, err
	}
	// Setup the CLI
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	cli := cli.New(&bud.Input{
		Dir:    c.dir,
		Env:    c.Env.List(),
		Stdin:  c.Stdin,
		Stdout: io.MultiWriter(os.Stdout, stdout),
		Stderr: io.MultiWriter(os.Stderr, stderr),
		Bus:    c.bus,
		WebLn:  webln,
		BudLn:  budln,
	})
	// Run the CLI
	ctx, cancel := context.WithCancel(ctx)
	eg, ctx := errgroup.WithContext(ctx)
	// Start running the CLI
	eg.Go(func() error { return cli.Run(ctx, prependFlags(args)...) })
	// App provides helpers and controls for the running CLI
	client := &Client{
		eg:     eg,
		log:    log,
		bus:    c.bus,
		stdout: stdout,
		stderr: stderr,
		webc:   webc,
		hotc:   budc,
		// Close function
		close: func() error {
			// Cancel the CLI
			cancel()
			// Wait for the CLI to finish
			return eg.Wait()
		},
	}
	// Wait for the client to be ready
	if err := client.Ready(ctx); err != nil {
		return nil, err
	}
	return client, nil
}

// Client for interacting with the running app
type Client struct {
	eg     *errgroup.Group
	log    log.Interface
	bus    pubsub.Client
	stdout *bytes.Buffer
	stderr *bytes.Buffer
	webc   *http.Client
	hotc   *http.Client
	once   once.Error
	close  func() error
}

// Stdout at a point in time
func (c *Client) Stdout() string {
	return c.stdout.String()
}

// Stderr at a point in time
func (c *Client) Stderr() string {
	return c.stderr.String()
}

// Close the app down
func (c *Client) Close() error {
	return c.once.Do(c.close)
}

// Hot connects to the event stream
func (c *Client) Hot(path string) (*hot.Stream, error) {
	return hot.DialWith(c.hotc, c.log, getURL(path))
}

func bufferHeaders(res *http.Response, body []byte) ([]byte, error) {
	// Coerce mime types before buffering the header
	if err := coerceMimes(res); err != nil {
		return nil, err
	}
	dump, err := httputil.DumpResponse(res, false)
	if err != nil {
		return nil, err
	}
	// httputil.DumpResponse() always attaches a Content-Length, regardless of
	// whether or not you remove it. This scanner removes the Content-Length
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
	elapsed := time.Since(dt)
	if elapsed > time.Minute {
		return fmt.Errorf("date header is too old %s", elapsed)
	}
	return nil
}

// Mime types are platform-specific. This function coerces them to a common
// name so they can be tested. This list will grow over time.
func coerceMimes(res *http.Response) error {
	contentType := res.Header.Get("Content-Type")
	if contentType == "" {
		return nil
	}
	// Coerce the JS content types to application/javascript
	switch contentType {
	case `text/javascript; charset=utf-8`:
		res.Header.Set("Content-Type", "application/javascript")
	}
	return nil
}

func (c *Client) Do(req *http.Request) (*Response, error) {
	res, err := c.webc.Do(req)
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

func getURL(path string) string {
	return "http://host" + path
}

// Wait for the application to be ready
func (c *Client) Ready(ctx context.Context) error {
	readySub := c.bus.Subscribe("app:ready")
	defer readySub.Close()
	errorSub := c.bus.Subscribe("app:error")
	defer errorSub.Close()
	for {
		select {
		case <-readySub.Wait():
			return nil
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errorSub.Wait():
			return errors.New(string(err))
		case <-time.After(time.Second * 1):
			c.log.Debug("testcli: waiting for app to be ready")
		}
	}
}

func (c *Client) Get(path string) (*Response, error) {
	c.log.Debug("testcli: get request", "path", path)
	req, err := c.GetRequest(path)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *Client) GetJSON(path string) (*Response, error) {
	c.log.Debug("testcli: get json request", "path", path)
	req, err := c.GetRequest(path)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	return c.Do(req)
}

func (c *Client) GetRequest(path string) (*http.Request, error) {
	return http.NewRequest(http.MethodGet, getURL(path), nil)
}

func (c *Client) Post(path string, body io.Reader) (*Response, error) {
	c.log.Debug("testcli: post request", "path", path)
	req, err := c.PostRequest(path, body)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *Client) PostJSON(path string, body io.Reader) (*Response, error) {
	c.log.Debug("testcli: post json request", "path", path)
	req, err := c.PostRequest(path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return c.Do(req)
}

func (c *Client) PostRequest(path string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(http.MethodPost, getURL(path), body)
}

func (c *Client) Patch(path string, body io.Reader) (*Response, error) {
	c.log.Debug("testcli: patch request", "path", path)
	req, err := c.PatchRequest(path, body)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *Client) PatchJSON(path string, body io.Reader) (*Response, error) {
	c.log.Debug("testcli: patch json request", "path", path)
	req, err := c.PatchRequest(path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return c.Do(req)
}

func (c *Client) PatchRequest(path string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(http.MethodPatch, getURL(path), body)
}

func (c *Client) Delete(path string, body io.Reader) (*Response, error) {
	c.log.Debug("testcli: delete request", "path", path)
	req, err := c.DeleteRequest(path, body)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *Client) DeleteJSON(path string, body io.Reader) (*Response, error) {
	c.log.Debug("testcli: delete json request", "path", path)
	req, err := c.DeleteRequest(path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return c.Do(req)
}

func (c *Client) DeleteRequest(path string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(http.MethodDelete, getURL(path), body)
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

func (r *Response) Diff(expected string) error {
	expected = strings.TrimSpace(dedent.Dedent(expected))
	actual := strings.TrimSpace(r.Dump().String())
	return difference(expected, actual)
}

func (r *Response) DiffHeaders(expected string) error {
	expected = strings.TrimSpace(dedent.Dedent(expected))
	actual := strings.TrimSpace(r.Headers().String())
	return difference(expected, actual)
}

func difference(expected, actual string) error {
	if expected == actual {
		return nil
	}
	var b bytes.Buffer
	b.WriteString("\n\x1b[4mExpected\x1b[0m:\n")
	b.WriteString(expected)
	b.WriteString("\n\n")
	b.WriteString("\x1b[4mActual\x1b[0m: \n")
	b.WriteString(actual)
	b.WriteString("\n\n")
	b.WriteString("\x1b[4mDifference\x1b[0m: \n")
	b.WriteString(diff.String(expected, actual))
	b.WriteString("\n")
	return errors.New(b.String())
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
