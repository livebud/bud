package test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/lithammer/dedent"
	"github.com/matryer/is"
	"github.com/matthewmueller/diff"
	"gitlab.com/mnm/bud/go/mod"
	"gitlab.com/mnm/bud/internal/generator"
	"gitlab.com/mnm/bud/internal/modcache"
	"gitlab.com/mnm/bud/internal/npm"
	"gitlab.com/mnm/bud/socket"
	"gitlab.com/mnm/bud/vfs"
)

func replaceBud(code string) (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("unable to load current file path")
	}
	budModule, err := mod.New().Find(filepath.Dir(file))
	if err != nil {
		return "", err
	}
	module, err := mod.New().Parse("go.mod", []byte(code))
	if err != nil {
		return "", err
	}
	err = module.File().Replace("gitlab.com/mnm/bud", budModule.Directory())
	if err != nil {
		return "", err
	}
	return string(module.File().Format()), nil
}

func addModules(code string, modules map[string]modcache.Files) (string, error) {
	module, err := mod.New().Parse("go.mod", []byte(code))
	if err != nil {
		return "", err
	}
	for modulePathVersion := range modules {
		moduleParts := strings.SplitN(modulePathVersion, "@", 2)
		if len(moduleParts) != 2 {
			return "", fmt.Errorf("modcache: invalid module key")
		}
		if err := module.File().AddRequire(moduleParts[0], moduleParts[1]); err != nil {
			return "", err
		}
	}
	return string(module.File().Format()), nil
}

func moduleGoMod(pathVersion string, module modcache.Files) (string, error) {
	moduleParts := strings.SplitN(pathVersion, "@", 2)
	if len(moduleParts) != 2 {
		return "", fmt.Errorf("modcache: invalid module key")
	}
	modulePath := moduleParts[0]
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("unable to load current file path")
	}
	budModule, err := mod.New().Find(filepath.Dir(file))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(dedent.Dedent(`
		module %s

		require (
			gitlab.com/mnm/bud v0.0.0
		)

		replace (
			gitlab.com/mnm/bud => %s
		)
	`), modulePath, budModule.Directory()), nil
}

// Cleanup individual files and root if no files left
func cleanup(t testing.TB, root, dir string) func() {
	t.Helper()
	is := is.New(t)
	return func() {
		if t.Failed() {
			return
		}
		is.NoErr(os.RemoveAll(dir))
		fis, err := os.ReadDir(root)
		if err != nil {
			return
		}
		if len(fis) > 0 {
			return
		}
		is.NoErr(os.RemoveAll(root))
	}
}

// TODO: remove gitlab.com/mnm/duo
var goMod = []byte(`
module app.com

require (
	gitlab.com/mnm/bud v0.0.0
	gitlab.com/mnm/duo v0.0.0-20220108212322-310ab0354067
)
`)

func Generator(t testing.TB) *Gen {
	return &Gen{
		t:       t,
		Modules: map[string]modcache.Files{},
		Files: map[string][]byte{
			"go.mod": goMod,
		},
		NodeModules: map[string]string{},
	}
}

type Gen struct {
	t           testing.TB
	Modules     map[string]modcache.Files // [module@version][filepath] = code
	Files       map[string][]byte
	NodeModules map[string]string // [package] = version
}

func (g *Gen) Generate() (*App, error) {
	ctx := context.Background()
	modCache := modcache.Default()
	// Add modules
	if len(g.Modules) > 0 {
		// Unable to cleanup after modcache because they're readonly and
		// the t.TempDir() fails with permission denied.
		cacheDir, err := ioutil.TempDir("", "modcache-*")
		if err != nil {
			return nil, err
		}
		modCache = modcache.New(cacheDir)
		// Generate a custom go.mod for the module
		for pathVersion, module := range g.Modules {
			if _, ok := module["go.mod"]; !ok {
				gomod, err := moduleGoMod(pathVersion, module)
				if err != nil {
					return nil, err
				}
				module["go.mod"] = gomod
			}
		}
		if err := modCache.Write(g.Modules); err != nil {
			return nil, err
		}
	}
	// Replace bud in Go mod if present
	if code, ok := g.Files["go.mod"]; ok {
		code, err := replaceBud(string(code))
		if err != nil {
			return nil, err
		}
		if len(g.Modules) > 0 {
			code, err = addModules(code, g.Modules)
			if err != nil {
				return nil, err
			}
		}
		g.Files["go.mod"] = []byte(code)
	}
	// Setup the application dir
	appDir := filepath.Join("_tmp", g.t.Name())
	if err := os.RemoveAll(appDir); err != nil {
		return nil, err
	}
	g.t.Cleanup(cleanup(g.t, "_tmp", appDir))
	// Add node_modules
	var nodeModules []string
	if len(g.NodeModules) > 0 {
		packageJSON := &npm.Package{
			Name:         filepath.Base(appDir),
			Version:      "0.0.0",
			Dependencies: map[string]string{},
		}
		for name, version := range g.NodeModules {
			nodeModules = append(nodeModules, name+"@"+version)
			packageJSON.Dependencies[name] = version
		}
		pkgJSON, err := json.MarshalIndent(packageJSON, "", "  ")
		if err != nil {
			return nil, err
		}
		g.Files["package.json"] = append(pkgJSON, '\n')
	}
	// Write the files to the application directory
	err := vfs.Write(appDir, vfs.Map(g.Files))
	if err != nil {
		return nil, err
	}
	// Install node_modules
	if len(nodeModules) > 0 {
		if err := npm.Install(appDir, nodeModules...); err != nil {
			return nil, err
		}
	}
	appFS := vfs.OS(appDir)
	gen, err := generator.Load(appFS, generator.WithCache(modCache))
	if err != nil {
		return nil, err
	}
	if err := gen.Generate(ctx); err != nil {
		return nil, err
	}
	return &App{
		t:      g.t,
		dir:    appDir,
		module: gen.Module(),
		env: env{
			"HOME":       os.Getenv("HOME"),
			"PATH":       os.Getenv("PATH"),
			"GOPATH":     os.Getenv("GOPATH"),
			"GOCACHE":    os.Getenv("GOCACHE"),
			"GOMODCACHE": modCache.Directory(),
			// TODO: remove once we can write a sum file to the modcache
			"GOPRIVATE": "*",
			"NO_COLOR":  "1",
		},
	}, nil
}

type env map[string]string

func (env env) List() (list []string) {
	for key, value := range env {
		list = append(list, key+"="+value)
	}
	return list
}

type App struct {
	t      testing.TB
	dir    string
	module *mod.Module
	env    env
	extras []*os.File
}

func (a *App) ExtraFiles(files ...*os.File) *App {
	a.extras = append(a.extras, files...)
	return a
}

func (a *App) Env(key, value string) *App {
	a.env[key] = value
	return a
}

func (a *App) build() (string, error) {
	binPath := filepath.Join(a.dir, "bud", "main")
	mainPath := filepath.Join("bud", "main.go")
	cmd := exec.Command("go", "build", "-o", binPath, "-mod", "mod", mainPath)
	cmd.Env = a.env.List()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Dir = a.dir
	return binPath, cmd.Run()
}

func (a *App) command(binPath string, args ...string) (*exec.Cmd, error) {
	cmd := exec.Command(binPath, args...)
	cmd.Env = a.env.List()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.ExtraFiles = a.extras
	cmd.Dir = a.dir
	return cmd, nil
}

func (a *App) Run(args ...string) string {
	binPath, err := a.build()
	if err != nil {
		return err.Error()
	}
	cmd, err := a.command(binPath, args...)
	if err != nil {
		return err.Error()
	}
	stdout := new(bytes.Buffer)
	cmd.Stdout = stdout
	if err := cmd.Run(); err != nil {
		return err.Error()
	}
	return stdout.String()
}

func (a *App) Start(args ...string) (*Server, error) {
	binPath, err := a.build()
	if err != nil {
		return nil, err
	}
	// Setup the command
	cmd, err := a.command(binPath, args...)
	if err != nil {
		return nil, err
	}
	// Start the unix domain socket
	socketPath := filepath.Join(a.t.TempDir(), "tmp.sock")
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
	files, env, err := socket.Files(ln)
	if err != nil {
		return nil, err
	}
	transport, err := socket.Transport(socketPath)
	if err != nil {
		return nil, err
	}
	// Add socket configuration to the command
	cmd.ExtraFiles = append(cmd.ExtraFiles, files...)
	cmd.Env = append(cmd.Env, string(env))
	// Start the webserver
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return &Server{
		t:   a.t,
		cmd: cmd,
		ln:  ln,
		client: &http.Client{
			// Timeout:   time.Second,
			Transport: transport,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}, nil
}

func (a *App) Exists(path string) bool {
	is := is.New(a.t)
	if _, err := fs.Stat(a.module, path); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false
		}
		is.NoErr(err)
	}
	return true
}

type Server struct {
	t      testing.TB
	cmd    *exec.Cmd
	ln     net.Listener
	client *http.Client
}

func (c *Server) Close() error {
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

func (a *Server) Request(req *http.Request) (*Response, error) {
	res, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	return &Response{
		t:        a.t,
		Response: res,
	}, nil
}

func (a *Server) Get(path string) (*Response, error) {
	req, err := http.NewRequest("GET", "http://host"+path, nil)
	if err != nil {
		return nil, err
	}
	return a.Request(req)
}

func (a *Server) GetJSON(path string) (*Response, error) {
	req, err := http.NewRequest("GET", "http://host"+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	return a.Request(req)
}

func (a *Server) Post(path string, body io.Reader) (*Response, error) {
	req, err := http.NewRequest("POST", "http://host"+path, body)
	if err != nil {
		return nil, err
	}
	return a.Request(req)
}

func (a *Server) PostJSON(path string, body io.Reader) (*Response, error) {
	req, err := http.NewRequest("POST", "http://host"+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return a.Request(req)
}

func (a *Server) Patch(path string, body io.Reader) (*Response, error) {
	req, err := http.NewRequest("PATCH", "http://host"+path, body)
	if err != nil {
		return nil, err
	}
	return a.Request(req)
}

func (a *Server) PatchJSON(path string, body io.Reader) (*Response, error) {
	req, err := http.NewRequest("PATCH", "http://host"+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return a.Request(req)
}

func (a *Server) Delete(path string, body io.Reader) (*Response, error) {
	req, err := http.NewRequest("DELETE", "http://host"+path, body)
	if err != nil {
		return nil, err
	}
	return a.Request(req)
}

func (a *Server) DeleteJSON(path string, body io.Reader) (*Response, error) {
	req, err := http.NewRequest("DELETE", "http://host"+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return a.Request(req)
}

type Response struct {
	t testing.TB
	*http.Response
}

var now = time.Date(2021, 12, 31, 0, 0, 0, 0, time.UTC)

func (r *Response) Expect(expect string) {
	r.t.Helper()
	// Make the date constant
	if v := r.Response.Header.Get("Date"); v != "" {
		r.Response.Header.Set("Date", now.Format(http.TimeFormat))
	}
	dumpBytes, err := httputil.DumpResponse(r.Response, true)
	if err != nil {
		diff.TestString(r.t, expect, err.Error())
		return
	}
	dump := string(dumpBytes)
	// Check the content length
	// TODO: clean this up. Ideally we can check this before dumping the response.
	// We just need to make sure to reset the body for the dump response
	if v := r.Response.Header.Get("Content-Length"); v != "" {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			diff.TestString(r.t, expect, err.Error())
			return
		}
		actual := fmt.Sprintf("Content-Length: %d", len(body))
		expect := fmt.Sprintf("Content-Length: %s", v)
		diff.TestString(r.t, expect, actual)
		dump = strings.ReplaceAll(dump, "\r\n"+expect, "")
	}
	diff.TestHTTP(r.t, expect, dump)
}

var contentLength = regexp.MustCompile(`\r\nContent-Length: \d+`)

func (r *Response) ExpectHeaders(expect string) {
	r.t.Helper()
	// Make the date constant
	if v := r.Response.Header.Get("Date"); v != "" {
		r.Response.Header.Set("Date", now.Format(http.TimeFormat))
	}
	dumpBytes, err := httputil.DumpResponse(r.Response, false)
	if err != nil {
		diff.TestString(r.t, expect, err.Error())
		return
	}
	dump := string(dumpBytes)
	dump = contentLength.ReplaceAllString(dump, "")
	diff.TestHTTP(r.t, expect, dump)
}

func (r *Response) Query(selector string) *goquery.Selection {
	r.t.Helper()
	doc, err := goquery.NewDocumentFromReader(r.Body)
	if err != nil {
		r.t.Fatal(err)
	}
	return doc.Find(selector)
}
