package test

import (
	"bytes"
	"context"
	"errors"
	"io/fs"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/go/mod"
	"gitlab.com/mnm/bud/internal/generator"
	"gitlab.com/mnm/bud/internal/modcache"
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

const goMod = `
module app.com

require (
	gitlab.com/mnm/bud v0.0.0
	gitlab.com/mnm/bud-tailwind v0.0.0-20211228175933-3ca601f1a518
)
`

func Generator(t testing.TB) *Gen {
	return &Gen{
		t:       t,
		Modules: map[string]modcache.Files{},
		Files: map[string]string{
			"go.mod": goMod,
		},
	}
}

type Gen struct {
	t       testing.TB
	Modules map[string]modcache.Files
	Files   map[string]string
}

func (g *Gen) Generate() (*App, error) {
	ctx := context.Background()
	modCache := modcache.Default()
	// Add modules
	if len(g.Modules) > 0 {
		cacheDir := g.t.TempDir()
		modCache = modcache.New(cacheDir)
		if err := modCache.Write(g.Modules); err != nil {
			return nil, err
		}
	}
	// Replace bud in Go mod if present
	if code, ok := g.Files["go.mod"]; ok {
		code, err := replaceBud(code)
		if err != nil {
			return nil, err
		}
		g.Files["go.mod"] = code
	}
	appDir := filepath.Join("_tmp", g.t.Name())
	g.t.Cleanup(cleanup(g.t, "_tmp", appDir))
	// Write the files to the application directory
	err := vfs.Write(appDir, vfs.Map(g.Files))
	if err != nil {
		return nil, err
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
			"NO_COLOR":   "1",
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
		cmd: cmd,
		ln:  ln,
		client: &http.Client{
			Timeout:   time.Second,
			Transport: transport,
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

func (a *Server) Request(req *http.Request) (*http.Response, error) {
	return a.client.Do(req)
}

func (a *Server) Get(path string) (*http.Response, error) {
	req, err := http.NewRequest("GET", "http://host"+path, nil)
	if err != nil {
		return nil, err
	}
	return a.Request(req)
}

func (a *Server) Post(path, body string) (*http.Response, error) {
	req, err := http.NewRequest("POST", "http://host"+path, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}
	return a.Request(req)
}

func (a *Server) Put(path, body string) (*http.Response, error) {
	req, err := http.NewRequest("PUT", "http://host"+path, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}
	return a.Request(req)
}

func (a *Server) Delete(path, body string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", "http://host"+path, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}
	return a.Request(req)
}
