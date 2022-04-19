package budtest

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/lithammer/dedent"
	"github.com/livebud/bud/internal/bud"
	"github.com/livebud/bud/internal/gobin"
	"github.com/livebud/bud/internal/imhash"
	"github.com/livebud/bud/internal/testdir"
	"github.com/livebud/bud/package/exe"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/modcache"
	"github.com/livebud/bud/package/socket"
	runtime_bud "github.com/livebud/bud/runtime/bud"
	"github.com/matthewmueller/diff"
)

func New(dir string) *Compiler {
	return &Compiler{
		dir: dir,
		Flag: runtime_bud.Flag{
			Embed:  false,
			Minify: false,
			Hot:    true,
		},
		Files:       map[string]string{},
		BFiles:      map[string][]byte{},
		Modules:     modcache.Modules{},
		NodeModules: map[string]string{},
		Env: bud.Env{
			"HOME":       os.Getenv("HOME"),
			"PATH":       os.Getenv("PATH"),
			"TMPDIR":     os.TempDir(),
			"GOPATH":     os.Getenv("GOPATH"),
			"GOMODCACHE": modcache.Default().Directory(),
			"NO_COLOR":   "1",
			// TODO: remove once we can write a sum file to the modcache
			"GOPRIVATE": "*",
		},
		CacheDir: filepath.Join(os.TempDir(), "bud", "cache"),
	}
}

type Compiler struct {
	dir         string
	Flag        runtime_bud.Flag
	Files       map[string]string // String files (convenient)
	BFiles      map[string][]byte // Byte files (for images and binaries)
	Modules     modcache.Modules  // name@version[path[data]]
	NodeModules map[string]string // name[version]
	Env         bud.Env
	CacheDir    string
}

func (c *Compiler) buildBud(ctx context.Context) (budPath string, err error) {
	// Find the bud module
	module, err := gomod.Find(".")
	if err != nil {
		return "", err
	}
	hash, err := imhash.Hash(module, ".")
	if err != nil {
		return "", err
	}
	budPath = filepath.Join(c.CacheDir, hash)
	if _, err := os.Stat(budPath); nil == err {
		return budPath, nil
	}
	if err := gobin.Build(ctx, module, module.Directory("main.go"), budPath); err != nil {
		return "", err
	}
	return budPath, nil
}

func (c *Compiler) Compile(ctx context.Context) (p *Project, err error) {
	dir, err := filepath.Abs(c.dir)
	if err != nil {
		return nil, err
	}
	// Setup directory
	td := testdir.New()
	td.Files = c.Files
	td.BFiles = c.BFiles
	td.Modules = c.Modules
	td.NodeModules = c.NodeModules
	if err := td.Write(dir); err != nil {
		return nil, err
	}
	// Build bud to use with the V8 Client
	budPath, err := c.buildBud(ctx)
	if err != nil {
		return nil, err
	}
	c.Env["BUD_PATH"] = budPath
	// Get the modCache
	modCache := testdir.ModCache(dir)
	c.Env["GOMODCACHE"] = modCache.Directory()
	// Find go.mod
	module, err := gomod.Find(dir, gomod.WithModCache(modCache))
	if err != nil {
		return nil, err
	}
	// Compile the project
	compiler, err := bud.Load(module)
	if err != nil {
		return nil, err
	}
	compiler.Env = c.Env
	compiler.ModCacheRW = true
	project, err := compiler.Compile(ctx, &c.Flag)
	if err != nil {
		return nil, err
	}
	return &Project{
		module:      module,
		project:     project,
		Files:       c.Files,
		BFiles:      c.BFiles,
		Modules:     c.Modules,
		NodeModules: c.NodeModules,
	}, nil
}

type Project struct {
	module  *gomod.Module
	project *bud.Project

	// Used to adjust files over time
	Files       map[string]string // String files (convenient)
	BFiles      map[string][]byte // Byte files (for images and binaries)
	Modules     modcache.Modules  // name@version[path[data]]
	NodeModules map[string]string // name[version]
}

// Rewrite files
func (p *Project) Rewrite() error {
	td := testdir.New()
	td.Files = p.Files
	td.BFiles = p.BFiles
	td.Modules = p.Modules
	td.NodeModules = p.NodeModules
	return td.Write(p.module.Directory(), testdir.WithBackup(false), testdir.WithSkip(func(name string, isDir bool) bool {
		return isDir && name == "bud"
	}))
}

func (p *Project) Exists(paths ...string) error {
	for _, path := range paths {
		if _, err := fs.Stat(p.module, path); err != nil {
			return err
		}
	}
	return nil
}

func (p *Project) Directory(paths ...string) string {
	return p.module.Directory(paths...)
}

func (p *Project) Execute(ctx context.Context, args ...string) (stdout Stdio, stderr Stdio, err error) {
	cmd := p.project.Executor(ctx, args...)
	sout := new(bytes.Buffer)
	serr := new(bytes.Buffer)
	cmd.Stdout = io.MultiWriter(os.Stdout, sout)
	cmd.Stderr = io.MultiWriter(os.Stderr, serr)
	if err := cmd.Run(); err != nil {
		return "", "", err
	}
	return Stdio(sout.String()), Stdio(serr.String()), nil
}

func (p *Project) Build(ctx context.Context) (*App, error) {
	app, err := p.project.Build(ctx)
	if err != nil {
		return nil, err
	}
	return &App{
		module: p.module,
		app:    app,
	}, nil
}

func (p *Project) Run(ctx context.Context) (*Server, error) {
	socketPath, err := getSocketPath()
	if err != nil {
		return nil, err
	}
	listener, err := socket.Listen(socketPath)
	if err != nil {
		return nil, err
	}
	client, err := httpClient(socketPath)
	if err != nil {
		return nil, err
	}
	process, err := p.project.Run(ctx, listener)
	if err != nil {
		return nil, err
	}
	return &Server{
		process:  process,
		listener: listener,
		client:   client,
	}, nil
}

type App struct {
	module *gomod.Module
	app    *runtime_bud.App
}

func (a *App) Exists(paths ...string) error {
	for _, path := range paths {
		if _, err := fs.Stat(a.module, path); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) NotExists(paths ...string) error {
	for _, path := range paths {
		if _, err := fs.Stat(a.module, path); nil == err {
			return fmt.Errorf("%q should not exist", path)
		} else if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}
	return nil
}

func (a *App) Execute(ctx context.Context, args ...string) (stdout Stdio, stderr Stdio, err error) {
	cmd := a.app.Executor(ctx, args...)
	sout := new(bytes.Buffer)
	serr := new(bytes.Buffer)
	cmd.Stdout = io.MultiWriter(os.Stdout, sout)
	cmd.Stderr = io.MultiWriter(os.Stderr, serr)
	if err := cmd.Run(); err != nil {
		return "", "", err
	}
	return Stdio(sout.String()), Stdio(serr.String()), nil
}

func getSocketPath() (string, error) {
	tmpdir, err := os.MkdirTemp("", "budtest-*")
	if err != nil {
		return "", err
	}
	return filepath.Join(tmpdir, "unix.sock"), nil
}

func httpClient(socketPath string) (*http.Client, error) {
	transport, err := socket.Transport(socketPath)
	if err != nil {
		return nil, err
	}
	return &http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}, nil
}

func (a *App) Start(ctx context.Context) (*Server, error) {
	socketPath, err := getSocketPath()
	if err != nil {
		return nil, err
	}
	listener, err := socket.Listen(socketPath)
	if err != nil {
		return nil, err
	}
	client, err := httpClient(socketPath)
	if err != nil {
		return nil, err
	}
	process, err := a.app.Start(ctx, listener)
	if err != nil {
		return nil, err
	}
	return &Server{
		process:  process,
		listener: listener,
		client:   client,
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
	process  *exe.Cmd
	listener net.Listener
	client   *http.Client
}

func (s *Server) Close() error {
	if err := s.process.Close(); err != nil {
		return err
	}
	if err := s.listener.Close(); err != nil {
		return err
	}
	return nil
}

func isExitStatus(err error) bool {
	return err != nil && strings.Contains(err.Error(), "exit status ")
}

func (s *Server) Restart(ctx context.Context) error {
	return s.process.Restart(ctx)
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

func (r *Response) ContainsBody(expect string) error {
	body, err := io.ReadAll(r.Response.Body)
	if err != nil {
		return err
	}
	if strings.Contains(string(body), expect) {
		return nil
	}
	return fmt.Errorf("%s does not contain %q", string(body), expect)
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
	s := new(strings.Builder)
	s.WriteString("\n\x1b[4mExpected\x1b[0m:\n")
	s.WriteString(expect)
	s.WriteString("\n\n")
	s.WriteString("\x1b[4mActual\x1b[0m: \n")
	s.WriteString(actual)
	s.WriteString("\n\n")
	s.WriteString("\x1b[4mDifference\x1b[0m: \n")
	s.WriteString(diff.String(expect, actual))
	s.WriteString("\n")
	return errors.New(s.String())
}

func diffHTTP(expect, actual string) error {
	expect = strings.TrimSpace(dedent.Dedent(expect))
	actual = strings.ReplaceAll(strings.TrimSpace(dedent.Dedent(actual)), "\r\n", "\n")
	return diffString(expect, actual)
}
