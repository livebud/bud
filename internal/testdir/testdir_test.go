package testdir_test

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/internal/testdir"
	"gitlab.com/mnm/bud/package/modcache"
)

func exists(fsys fs.FS, paths ...string) error {
	for _, path := range paths {
		if _, err := fs.Stat(fsys, path); err != nil {
			return err
		}
	}
	return nil
}

func notExists(fsys fs.FS, paths ...string) error {
	for _, path := range paths {
		if _, err := fs.Stat(fsys, path); nil == err {
			return fmt.Errorf("%s exists", path)
		} else if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}
	return nil
}

func TestDir(t *testing.T) {
	is := is.New(t)
	td := testdir.New()
	td.Modules = map[string]modcache.Files{
		"gitlab.com/mnm/bud-tailwind@v0.0.1": modcache.Files{
			"public/tailwind/preflight.css": `/* tailwind */`,
		},
	}
	td.Files["action/action.go"] = `package action`
	td.BFiles["public/favicon.ico"] = []byte{0x00}
	td.NodeModules["svelte"] = `3.46.4`
	td.NodeModules["livebud"] = "*"
	dir := t.TempDir()
	err := td.Write(dir, testdir.WithBackup(false))
	is.NoErr(err)
	err = exists(os.DirFS(dir),
		"action/action.go",
		"public/favicon.ico",
		".mod/gitlab.com/mnm/bud-tailwind@v0.0.1/public/tailwind/preflight.css",
		".npm/livebud.tgz",
		"node_modules/svelte/package.json",
		"node_modules/livebud/package.json",
		"package.json",
		"go.mod",
	)
	is.NoErr(err)
}

func TestRefresh(t *testing.T) {
	is := is.New(t)
	td := testdir.New()
	td.Modules = map[string]modcache.Files{
		"gitlab.com/mnm/bud-tailwind@v0.0.1": modcache.Files{
			"public/tailwind/preflight.css": `/* tailwind */`,
		},
	}
	td.Files["action/action.go"] = `package action`
	td.BFiles["public/favicon.ico"] = []byte{0x00}
	td.NodeModules["svelte"] = `3.46.4`
	td.NodeModules["livebud"] = "*"
	dir := t.TempDir()
	err := td.Write(dir, testdir.WithBackup(false))
	is.NoErr(err)
	err = exists(os.DirFS(dir),
		"action/action.go",
		"public/favicon.ico",
		".mod/gitlab.com/mnm/bud-tailwind@v0.0.1/public/tailwind/preflight.css",
		".npm/livebud.tgz",
		"node_modules/svelte/package.json",
		"node_modules/livebud/package.json",
		"package.json",
		"go.mod",
	)
	td.Modules = map[string]modcache.Files{}
	delete(td.Files, "action/action.go")
	delete(td.BFiles, "public/favicon.ico")
	err = td.Write(dir, testdir.WithBackup(false))
	is.NoErr(err)
	is.NoErr(notExists(os.DirFS(dir),
		"action/action.go",
		"public/favicon.ico",
		".mod/gitlab.com/mnm/bud-tailwind@v0.0.1/public/tailwind/preflight.css",
	))
	is.NoErr(exists(os.DirFS(dir),
		".npm/livebud.tgz",
		"node_modules/livebud/package.json",
		"node_modules/svelte/package.json",
		"package.json",
		"go.mod",
	))
}

func TestSkip(t *testing.T) {
	is := is.New(t)
	td := testdir.New()
	td.Modules = map[string]modcache.Files{
		"gitlab.com/mnm/bud-tailwind@v0.0.1": modcache.Files{
			"public/tailwind/preflight.css": `/* tailwind */`,
		},
	}
	td.Files["action/action.go"] = `package action`
	td.BFiles["public/favicon.ico"] = []byte{0x00}
	td.NodeModules["svelte"] = `3.46.4`
	td.NodeModules["livebud"] = "*"
	dir := t.TempDir()
	err := td.Write(dir, testdir.WithBackup(false))
	is.NoErr(err)
	err = exists(os.DirFS(dir),
		"action/action.go",
		"public/favicon.ico",
		".mod/gitlab.com/mnm/bud-tailwind@v0.0.1/public/tailwind/preflight.css",
		".npm/livebud.tgz",
		"node_modules/svelte/package.json",
		"node_modules/livebud/package.json",
		"package.json",
		"go.mod",
	)
	td.Modules = map[string]modcache.Files{}
	delete(td.Files, "action/action.go")
	delete(td.BFiles, "public/favicon.ico")
	err = td.Write(dir, testdir.WithBackup(false), testdir.WithSkip(func(name string, isDir bool) bool {
		return (name == "action" && isDir) || (name == ".mod" && isDir)
	}))
	is.NoErr(err)
	is.NoErr(notExists(os.DirFS(dir),
		"public/favicon.ico",
	))
	is.NoErr(exists(os.DirFS(dir),
		"action/action.go",
		".mod/gitlab.com/mnm/bud-tailwind@v0.0.1/public/tailwind/preflight.css",
		".npm/livebud.tgz",
		"node_modules/livebud/package.json",
		"node_modules/svelte/package.json",
		"package.json",
		"go.mod",
	))
}