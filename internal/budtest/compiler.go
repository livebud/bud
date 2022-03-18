package budtest

import (
	"bytes"
	"context"
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
	"gitlab.com/mnm/bud/internal/bud"
	"gitlab.com/mnm/bud/internal/testdir"
	"gitlab.com/mnm/bud/pkg/gomod"
	"gitlab.com/mnm/bud/pkg/socket"
)

func Find(dir string) (*Compiler, error) {
	module, err := gomod.Find(dir)
	if err != nil {
		return nil, err
	}
	compiler := bud.New(module)
	if err != nil {
		return nil, err
	}
	flag := bud.Flag{
		Embed:  false,
		Minify: false,
		Hot:    true,
	}
	return &Compiler{compiler, flag, module}, nil
}

type Compiler struct {
	compiler *bud.Compiler
	Flag     bud.Flag
	module   *gomod.Module
}

func (c *Compiler) Compile(ctx context.Context) (p *Project, err error) {
	project, err := c.compiler.Compile(ctx, c.Flag)
	if err != nil {
		return nil, err
	}
	project.Env["GOCACHE"] = os.Getenv("GOCACHE")
	project.Env["GOMODCACHE"] = testdir.ModCache(c.module.Directory()).Directory()
	project.Env["NO_COLOR"] = "1"
	// TODO: remove once we can write a sum file to the modcache
	project.Env["GOPRIVATE"] = "*"
	return &Project{project}, nil
}

type Project struct {
	project *bud.Project
}

func (p *Project) Execute(ctx context.Context, args ...string) (stdout Stdio, stderr Stdio, err error) {
	cmd := p.project.Executor(ctx, args...)
	sout := new(bytes.Buffer)
	serr := new(bytes.Buffer)
	cmd.Stdout = sout
	cmd.Stderr = serr
	if err := cmd.Run(); err != nil {
		return "", "", err
	}
	return Stdio(sout.String()), Stdio(serr.String()), nil
}

func (p *Project) Run(ctx context.Context) (*Server, error) {
	// Since we're starting the web server, initialize a unix domain socket
	// to listen and pass that socket to the application process
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
	listener, err := socket.Listen(socketPath)
	if err != nil {
		return nil, err
	}
	transport, err := socket.Transport(socketPath)
	if err != nil {
		return nil, err
	}
	cmd, err := p.project.Runner(ctx, listener)
	if err != nil {
		return nil, err
	}
	// Start the webserver
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return &Server{
		cmd: cmd,
		ln:  listener,
		client: &http.Client{
			Timeout:   10 * time.Second,
			Transport: transport,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}, nil
}

type Stdio string

func (s Stdio) Contains(c string) error {
	if strings.Contains(string(s), c) {
		return nil
	}
	return fmt.Errorf("%s does not contain %q", s, c)
}

func (s Stdio) Expect(a string) error {
	if string(s) == a {
		return nil
	}
	return fmt.Errorf("%q does not equal %q", s, a)
}

func (s Stdio) String() string {
	return string(s)
}

type Server struct {
	cmd    *exec.Cmd
	ln     net.Listener
	client *http.Client
}

func (s *Server) Close() error {
	p := s.cmd.Process
	if p != nil {
		if err := p.Signal(os.Interrupt); err != nil {
			p.Kill()
		}
	}
	if err := s.cmd.Wait(); err != nil {
		return err
	}
	if err := s.ln.Close(); err != nil {
		return err
	}
	return nil
}

func isExitStatus(err error) bool {
	return err != nil && strings.Contains(err.Error(), "exit status ")
}

func (s *Server) Restart() error {
	p := s.cmd.Process
	if p != nil {
		if err := p.Signal(os.Interrupt); err != nil {
			p.Kill()
		}
	}
	if err := s.cmd.Wait(); err != nil {
		if !isExitStatus(err) {
			return err
		}
	}
	cmd := exec.Command(s.cmd.Path, s.cmd.Args...)
	cmd.Env = s.cmd.Env
	cmd.Stdout = s.cmd.Stdout
	cmd.Stderr = s.cmd.Stderr
	cmd.Stdin = s.cmd.Stdin
	cmd.ExtraFiles = s.cmd.ExtraFiles
	cmd.Dir = s.cmd.Dir
	if err := cmd.Start(); err != nil {
		return err
	}
	return nil
}

func (s *Server) Request(req *http.Request) (*Response, error) {
	res, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	return &Response{
		Response: res,
	}, nil
}

func (s *Server) Get(path string) (*Response, error) {
	req, err := http.NewRequest("GET", "http://host"+path, nil)
	if err != nil {
		return nil, err
	}
	return s.Request(req)
}

func (s *Server) GetJSON(path string) (*Response, error) {
	req, err := http.NewRequest("GET", "http://host"+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	return s.Request(req)
}

func (s *Server) Post(path string, body io.Reader) (*Response, error) {
	req, err := http.NewRequest("POST", "http://host"+path, body)
	if err != nil {
		return nil, err
	}
	return s.Request(req)
}

func (s *Server) PostJSON(path string, body io.Reader) (*Response, error) {
	req, err := http.NewRequest("POST", "http://host"+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return s.Request(req)
}

func (s *Server) Patch(path string, body io.Reader) (*Response, error) {
	req, err := http.NewRequest("PATCH", "http://host"+path, body)
	if err != nil {
		return nil, err
	}
	return s.Request(req)
}

func (s *Server) PatchJSON(path string, body io.Reader) (*Response, error) {
	req, err := http.NewRequest("PATCH", "http://host"+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return s.Request(req)
}

func (s *Server) Delete(path string, body io.Reader) (*Response, error) {
	req, err := http.NewRequest("DELETE", "http://host"+path, body)
	if err != nil {
		return nil, err
	}
	return s.Request(req)
}

func (s *Server) DeleteJSON(path string, body io.Reader) (*Response, error) {
	req, err := http.NewRequest("DELETE", "http://host"+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return s.Request(req)
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
