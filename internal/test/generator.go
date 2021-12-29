package test

import (
	"bytes"
	"context"
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/gen"
	"gitlab.com/mnm/bud/go/mod"
	"gitlab.com/mnm/bud/internal/generator"
	"gitlab.com/mnm/bud/internal/modcache"
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

func Generator(t testing.TB) *Gen {
	return &Gen{
		t:       t,
		Modules: map[string]modcache.Files{},
		Files:   map[string]string{},
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
	appDir := g.t.TempDir()
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

func (a *App) run(binPath string, args ...string) (*exec.Cmd, error) {
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
	cmd, err := a.run(binPath, args...)
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

func (a *App) Start(args ...string) (*Command, error) {
	binPath, err := a.build()
	if err != nil {
		return nil, err
	}
	cmd, err := a.run(binPath, args...)
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return &Command{
		cmd: cmd,
	}, nil
}

type Command struct {
	cmd *exec.Cmd
}

func (c *Command) Wait() error {
	return c.cmd.Wait()
}

func (c *Command) Close() error {
	p := c.cmd.Process
	if p != nil {
		if err := p.Signal(os.Interrupt); err != nil {
			p.Kill()
		}
	}
	return c.cmd.Wait()
}

func (a *App) Exists(path string) bool {
	is := is.New(a.t)
	if _, err := fs.Stat(a.module, path); err != nil {
		if errors.Is(err, fs.ErrNotExist) || errors.Is(err, gen.ErrSkipped) {
			return false
		}
		is.NoErr(err)
	}
	return true
}
