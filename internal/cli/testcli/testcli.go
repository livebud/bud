package testcli

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/livebud/bud/internal/errs"
	"github.com/livebud/bud/internal/pubsub"

	"github.com/livebud/bud/package/socket"

	"golang.org/x/sync/errgroup"

	"github.com/livebud/bud/internal/cli"
	"github.com/livebud/bud/internal/envs"
)

var log = flag.String("log", "info", "choose a log level")

func New(dir string) *CLI {
	return &CLI{
		dir: dir,
		bus: pubsub.New(),
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
		"--log=" + *log,
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
		// This is extra high right now because we don't currently have any signal
		// that we've built the app and `bud run --embed` can take a long time. This
		// is going to slow down legitimately failing requests, so it's a very
		// temporary solution.
		//
		// TODO: support getting a signal that we've built the app, then lower this
		// deadline.
		Timeout:   60 * time.Second,
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	return listener, client, nil
}

func (c *CLI) Start(ctx context.Context, args ...string) (*App, error) {
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
	eg := new(errgroup.Group)
	ctx, cancelCLI := context.WithCancel(ctx)
	eg.Go(func() error { return cli.Run(ctx, prependFlags(args)...) })
	// App provides helpers and controls for the running CLI
	return &App{
		eg:           eg,
		bus:          c.bus,
		stdoutReader: stdoutReader,
		stdoutWriter: stdoutWriter,
		stderrReader: stderrReader,
		stderrWriter: stderrWriter,
		webClient:    webClient,
		hotClient:    hotClient,
		cancelCLI:    cancelCLI,
	}, nil
}

type App struct {
	eg           *errgroup.Group
	bus          pubsub.Client
	stdoutReader io.Reader
	stdoutWriter io.WriteCloser
	stderrReader io.Reader
	stderrWriter io.WriteCloser
	webClient    *http.Client
	hotClient    *http.Client
	cancelCLI    context.CancelFunc
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

// Close the app down
func (a *App) Close() (err error) {
	// Cancel the CLI
	a.cancelCLI()
	// Close the writers
	err = errs.Join(err, a.stdoutWriter.Close())
	err = errs.Join(err, a.stderrWriter.Close())
	// Wait for the CLI to finish
	err = errs.Join(err, a.eg.Wait())
	return err
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
