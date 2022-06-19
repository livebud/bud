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
	"github.com/livebud/bud/internal/errs"
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

func (c *CLI) toCLI() *cli.CLI {
	return &cli.CLI{
		Dir:    c.dir,
		Bus:    c.bus,
		Env:    c.Env.List(),
		Stdin:  c.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

// Flags that can be set from the test suite
// These can be overriden by more specific flags
func prependFlags(args []string) []string {
	return append([]string{
		"--log=" + testlog.Pattern(),
	}, args...)
}

func (c *CLI) Run(ctx context.Context, args ...string) (*Result, error) {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	cli := c.toCLI()
	cli.Stdout = io.MultiWriter(os.Stdout, stdout)
	cli.Stderr = io.MultiWriter(os.Stderr, stderr)
	err := cli.Run(ctx, prependFlags(args)...)
	return &Result{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}, err
}

type Result struct {
	Stdout string
	Stderr string
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
		Timeout:   5 * time.Second,
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	return listener, client, nil
}

func (c *CLI) Start(ctx context.Context, args ...string) (*App, error) {
	log := testlog.Log()
	stdoutReader, stdoutWriter, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	stderrReader, stderrWriter, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	// TODO: listen unix and create client
	webListener, webClient, err := listen(":0")
	if err != nil {
		return nil, err
	}
	// TODO: listen unix and create client
	hotListener, hotClient, err := listen(":0")
	if err != nil {
		return nil, err
	}
	// Setup the CLI
	cli := c.toCLI()
	cli.Stdout = io.MultiWriter(os.Stdout, stdoutWriter)
	cli.Stderr = io.MultiWriter(os.Stderr, stderrWriter)
	cli.Web = webListener
	cli.Hot = hotListener
	// Run the CLI
	ctx, cancel := context.WithCancel(ctx)
	eg, ctx := errgroup.WithContext(ctx)
	// App provides helpers and controls for the running CLI
	app := &App{
		eg:           eg,
		log:          log,
		bus:          c.bus,
		stdoutReader: stdoutReader,
		stdoutWriter: stdoutWriter,
		stderrReader: stderrReader,
		stderrWriter: stderrWriter,
		webClient:    webClient,
		hotClient:    hotClient,
		cancel:       cancel,
		appReadySub:  c.bus.Subscribe("app:ready"),
		cmdErrorSub:  c.bus.Subscribe("cmd:error"),
	}
	// Start running the CLI
	eg.Go(func() error {
		return cli.Run(ctx, prependFlags(args)...)
	})
	// Wait for the app to be ready
	if err := app.Ready(ctx); err != nil {
		cancel()
		return nil, err
	}
	return app, nil
}

type App struct {
	eg           *errgroup.Group
	log          log.Interface
	bus          pubsub.Client
	stdoutReader io.Reader
	stdoutWriter io.WriteCloser
	stderrReader io.Reader
	stderrWriter io.WriteCloser
	webClient    *http.Client
	hotClient    *http.Client
	cancel       context.CancelFunc
	appReadySub  pubsub.Subscription
	cmdErrorSub  pubsub.Subscription
}

func (a *App) Stdout() io.Reader {
	return a.stdoutReader
}

func (a *App) Stderr() io.Reader {
	return a.stdoutReader
}

// Subscribe to an event
func (a *App) Subscribe(topics ...string) pubsub.Subscription {
	return a.bus.Subscribe(topics...)
}

// Publish an event
func (a *App) Publish(topic string, payload []byte) {
	a.bus.Publish(topic, payload)
}

func (a *App) Ready(ctx context.Context) error {
	// Wait for the app to be ready or fail trying
	a.log.Debug("testcli: waiting for app to be ready")
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("testcli: app exited quickly after starting. %w", ctx.Err())
		case err := <-a.cmdErrorSub.Wait():
			return fmt.Errorf("testcli: error starting the generated app. %s", string(err))
		case <-a.appReadySub.Wait():
			a.log.Debug("testcli: app is ready")
			return nil
		case <-time.Tick(time.Second):
			a.log.Debug("testcli: app is not ready yet")
		}
	}
}

// Close the app down
func (a *App) Close() (err error) {
	// Cancel the CLI
	a.cancel()
	// Close the subscriptions
	a.appReadySub.Close()
	a.cmdErrorSub.Close()
	// Close the writers
	err = errs.Join(err, a.stdoutWriter.Close())
	err = errs.Join(err, a.stderrWriter.Close())
	// Wait for the CLI to finish
	err = errs.Join(err, a.eg.Wait())
	return err
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

func (a *App) Request(req *http.Request) (*Response, error) {
	res, err := a.webClient.Do(req)
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
