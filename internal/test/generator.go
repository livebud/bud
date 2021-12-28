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
	"gitlab.com/mnm/bud/fsync"
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
	appFS := vfs.Map(g.Files)
	gen, err := generator.Load(appFS, generator.WithCache(modCache))
	if err != nil {
		return nil, err
	}
	if err := gen.Generate(ctx); err != nil {
		return nil, err
	}
	return &App{
		t:     g.t,
		cache: modCache,
		fs:    gen.Module(),
	}, nil
}

type App struct {
	t      testing.TB
	fs     fs.FS
	cache  *modcache.Cache
	runDir string // Initially empty
}

func (a *App) Run(args ...string) string {
	if a.runDir == "" {
		a.runDir = a.t.TempDir()
		if err := fsync.Dir(a.fs, ".", vfs.OS(a.runDir), "."); err != nil {
			return err.Error()
		}
	}
	mainPath := filepath.Join("bud", "main.go")
	args = append([]string{"run", "-mod=mod", mainPath}, args...)
	cmd := exec.Command("go", args...)
	cmd.Env = append(os.Environ(), "GOMODCACHE="+a.cache.Directory(), "GOPRIVATE=*", "NO_COLOR=1")
	stdout := new(bytes.Buffer)
	cmd.Stdout = stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Dir = a.runDir
	if err := cmd.Run(); err != nil {
		return err.Error()
	}
	return stdout.String()
}

func (a *App) Exists(path string) bool {
	is := is.New(a.t)
	if _, err := fs.Stat(a.fs, path); err != nil {
		if errors.Is(err, fs.ErrNotExist) || errors.Is(err, gen.ErrSkipped) {
			return false
		}
		is.NoErr(err)
	}
	return true
}
