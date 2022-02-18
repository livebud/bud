package testapp

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/lithammer/dedent"
	"github.com/matthewmueller/diff"
	"gitlab.com/mnm/bud/pkg/socket"
)

func New(dir string) *App {
	return &App{dir, map[string]string{}}
}

type App struct {
	dir string
	Env map[string]string
}

func (a *App) Run(args ...string) (stdout string, stderr string, err error) {
	cmd := exec.Command("bud/app", args...)
	cmd.Dir = a.dir

	// Setup stdio
	stdo := new(bytes.Buffer)
	cmd.Stdout = stdo
	stde := new(bytes.Buffer)
	cmd.Stderr = stde

	// Load the environment
	env := make([]string, len(a.Env))
	for key, value := range a.Env {
		env = append(env, key+"="+value)
	}
	cmd.Env = env

	// Run the command
	err = cmd.Run()
	if isExitStatus(err) {
		err = fmt.Errorf("%w: %s", err, stde.String())
	}
	return stdo.String(), stde.String(), err
}

func (a *App) Start(args ...string) (*Process, error) {
	cmd := exec.Command("bud/app", args...)
	cmd.Dir = a.dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// If we're starting the web server, initialize a unix domain socket
	// to listen and pass that socket to the application process
	var ln net.Listener
	var transport http.RoundTripper
	if len(args) == 0 {
		// Start the unix domain socket
		name, err := ioutil.TempDir("", "bud-testapp-*")
		if err != nil {
			return nil, err
		}
		socketPath := filepath.Join(name, "tmp.sock")
		// Heads up: If you see the `bind: invalid argument` error, there's a chance
		// the path is too long. 103 characters appears to be the limit on OSX,
		// https://github.com/golang/go/issues/6895.
		if len(socketPath) > 103 {
			return nil, fmt.Errorf("socket name is too long")
		}
		ln, err := socket.Listen(socketPath)
		if err != nil {
			return nil, err
		}
		transport, err = socket.Transport(socketPath)
		if err != nil {
			return nil, err
		}
		// Add socket configuration to the command
		files, env, err := socket.Files(ln)
		if err != nil {
			return nil, err
		}
		cmd.ExtraFiles = append(cmd.ExtraFiles, files...)
		a.Env[env.Key()] = env.Value()
	}

	// Load the environment
	env := make([]string, len(a.Env))
	for key, value := range a.Env {
		env = append(env, key+"="+value)
	}
	cmd.Env = env

	// Start the webserver
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return &Process{
		cmd: cmd,
		ln:  ln,
		client: &http.Client{
			Timeout:   10 * time.Second,
			Transport: transport,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}, nil
}

type Process struct {
	cmd    *exec.Cmd
	ln     net.Listener // could be nil
	client *http.Client
}

func (c *Process) Close() error {
	p := c.cmd.Process
	if p != nil {
		if err := p.Signal(os.Interrupt); err != nil {
			p.Kill()
		}
	}
	if err := c.cmd.Wait(); err != nil {
		return err
	}
	if err := c.ln.Close(); err != nil {
		return err
	}
	return nil
}

func (c *Process) Restart() error {
	p := c.cmd.Process
	if p != nil {
		if err := p.Signal(os.Interrupt); err != nil {
			p.Kill()
		}
	}
	if err := c.cmd.Wait(); err != nil {
		if !isExitStatus(err) {
			return err
		}
	}
	cmd := exec.Command(c.cmd.Path, c.cmd.Args...)
	cmd.Env = c.cmd.Env
	cmd.Stdout = c.cmd.Stdout
	cmd.Stderr = c.cmd.Stderr
	cmd.Stdin = c.cmd.Stdin
	cmd.ExtraFiles = c.cmd.ExtraFiles
	cmd.Dir = c.cmd.Dir
	if err := cmd.Start(); err != nil {
		return err
	}
	return nil
}

func isExitStatus(err error) bool {
	return err != nil && strings.Contains(err.Error(), "exit status ")
}

func (a *Process) Request(req *http.Request) (*Response, error) {
	res, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	return &Response{
		Response: res,
	}, nil
}

func (a *Process) Get(path string) (*Response, error) {
	req, err := http.NewRequest("GET", "http://host"+path, nil)
	if err != nil {
		return nil, err
	}
	return a.Request(req)
}

func (a *Process) GetJSON(path string) (*Response, error) {
	req, err := http.NewRequest("GET", "http://host"+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	return a.Request(req)
}

func (a *Process) Post(path string, body io.Reader) (*Response, error) {
	req, err := http.NewRequest("POST", "http://host"+path, body)
	if err != nil {
		return nil, err
	}
	return a.Request(req)
}

func (a *Process) PostJSON(path string, body io.Reader) (*Response, error) {
	req, err := http.NewRequest("POST", "http://host"+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return a.Request(req)
}

func (a *Process) Patch(path string, body io.Reader) (*Response, error) {
	req, err := http.NewRequest("PATCH", "http://host"+path, body)
	if err != nil {
		return nil, err
	}
	return a.Request(req)
}

func (a *Process) PatchJSON(path string, body io.Reader) (*Response, error) {
	req, err := http.NewRequest("PATCH", "http://host"+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return a.Request(req)
}

func (a *Process) Delete(path string, body io.Reader) (*Response, error) {
	req, err := http.NewRequest("DELETE", "http://host"+path, body)
	if err != nil {
		return nil, err
	}
	return a.Request(req)
}

func (a *Process) DeleteJSON(path string, body io.Reader) (*Response, error) {
	req, err := http.NewRequest("DELETE", "http://host"+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return a.Request(req)
}

type Response struct {
	*http.Response
}

var now = time.Date(2021, 12, 31, 0, 0, 0, 0, time.UTC)

func (r *Response) Expect(expect string) error {
	// Make the date constant
	if v := r.Response.Header.Get("Date"); v != "" {
		r.Response.Header.Set("Date", now.Format(http.TimeFormat))
	}
	dumpBytes, err := httputil.DumpResponse(r.Response, true)
	if err != nil {
		return diffString(expect, err.Error())
	}
	dump := string(dumpBytes)
	// Check the content length
	// TODO: clean this up. Ideally we can check this before dumping the response.
	// We just need to make sure to reset the body for the dump response
	if v := r.Response.Header.Get("Content-Length"); v != "" {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return diffString(expect, err.Error())
		}
		actual := fmt.Sprintf("Content-Length: %d", len(body))
		expect := fmt.Sprintf("Content-Length: %s", v)
		if err := diffString(expect, actual); err != nil {
			return err
		}
		dump = strings.ReplaceAll(dump, "\r\n"+expect, "")
	}
	return diffHTTP(expect, dump)
}

var contentLength = regexp.MustCompile(`\r\nContent-Length: \d+`)

func (r *Response) ExpectHeaders(expect string) error {
	// Make the date constant
	if v := r.Response.Header.Get("Date"); v != "" {
		r.Response.Header.Set("Date", now.Format(http.TimeFormat))
	}
	dumpBytes, err := httputil.DumpResponse(r.Response, false)
	if err != nil {
		return diffString(expect, err.Error())
	}
	dump := string(dumpBytes)
	dump = contentLength.ReplaceAllString(dump, "")
	return diffHTTP(expect, dump)
}

func (r *Response) Query(selector string) (*goquery.Selection, error) {
	doc, err := goquery.NewDocumentFromReader(r.Body)
	if err != nil {
		return nil, err
	}
	return doc.Find(selector), nil
}

func diffString(expect, actual string) error {
	if expect == actual {
		return nil
	}
	return errors.New(diff.String(expect, actual))
}

func diffHTTP(expect, actual string) error {
	expect = strings.TrimSpace(dedent.Dedent(expect))
	actual = strings.ReplaceAll(strings.TrimSpace(dedent.Dedent(actual)), "\r\n", "\n")
	if expect == actual {
		return nil
	}
	return errors.New(diff.String(expect, actual))
}
