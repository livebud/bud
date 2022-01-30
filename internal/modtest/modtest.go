// Package modtest is for quickly building a usable module in memory.
package modtest

import (
	"os"
	"strings"
	"testing"

	"github.com/lithammer/dedent"
	"github.com/matryer/is"
	"gitlab.com/mnm/bud/internal/modcache"
	"gitlab.com/mnm/bud/pkg/gomod"
	"gitlab.com/mnm/bud/vfs"
)

type Module struct {
	Modules  map[string]map[string]string
	Files    map[string][]byte
	CacheDir string
	AppDir   string
	// FS       fs.FS
}

func redent(s string) string {
	return strings.TrimSpace(dedent.Dedent(s)) + "\n"
}

// func noWrap(fsys fs.FS) fs.FS {
// 	return fsys
// }

func Make(t testing.TB, m Module) *gomod.Module {
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
				m.Files[path] = []byte(replaceBud(t, string(file)))
				continue
			}
			m.Files[path] = []byte(redent(string(file)))
		}
		err := vfs.Write(m.AppDir, vfs.Map(m.Files))
		is.NoErr(err)
	}
	module, err := gomod.Find(m.AppDir, gomod.WithModCache(modCache))
	is.NoErr(err)
	return module
}

func replaceBud(t testing.TB, code string) string {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	budModule, err := gomod.Find(wd)
	is.NoErr(err)
	module, err := gomod.Parse("go.mod", []byte(code))
	is.NoErr(err)
	err = module.File().Replace("gitlab.com/mnm/bud", budModule.Directory())
	is.NoErr(err)
	return string(module.File().Format())
}
