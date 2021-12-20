// Package modtest is for quickly building a usable module in memory.
package modtest

import (
	"io/fs"
	"os"
	"strings"
	"testing"

	"github.com/lithammer/dedent"
	"github.com/matryer/is"
	"gitlab.com/mnm/bud/go/mod"
	"gitlab.com/mnm/bud/internal/modcache"
	"gitlab.com/mnm/bud/vfs"
)

type Module struct {
	Modules  map[string]map[string]string
	Files    map[string]string
	CacheDir string
	AppDir   string
	FS       fs.FS
}

func redent(s string) string {
	return strings.TrimSpace(dedent.Dedent(s)) + "\n"
}

func noWrap(fsys fs.FS) fs.FS {
	return fsys
}

func Make(t testing.TB, m Module) *mod.Module {
	t.Helper()
	is := is.New(t)
	modCache := modcache.Default()
	if m.Modules != nil {
		if m.CacheDir == "" {
			m.CacheDir = t.TempDir()
		}
		modCache = modcache.New(m.CacheDir)
		err := modCache.Write(m.Modules)
		is.NoErr(err)
	}
	if m.AppDir == "" {
		m.AppDir = t.TempDir()
	}
	if m.Files != nil {
		for path, file := range m.Files {
			if path == "go.mod" {
				m.Files[path] = replaceBud(t, file)
				continue
			}
			m.Files[path] = redent(file)
		}
		err := vfs.Write(m.AppDir, vfs.Map(m.Files))
		is.NoErr(err)
	}
	if m.FS == nil {
		m.FS = os.DirFS(m.AppDir)
	}
	modFinder := mod.New(mod.WithCache(modCache), mod.WithFS(m.FS))
	module, err := modFinder.Find(".")
	is.NoErr(err)
	return module
}

func replaceBud(t testing.TB, code string) string {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	budModule, err := mod.New().Find(wd)
	is.NoErr(err)
	module, err := mod.New().Parse("go.mod", []byte(code))
	is.NoErr(err)
	err = module.File().Replace("gitlab.com/mnm/bud", budModule.Directory())
	is.NoErr(err)
	return string(module.File().Format())
}
