package testdir_test

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"testing"

	"github.com/livebud/bud/internal/testdir"
	"github.com/livebud/bud/internal/version"
	"github.com/matryer/is"
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
	td.Modules["github.com/livebud/bud-test-plugin"] = "v0.0.2"
	td.Files["controller/controller.go"] = `package controller`
	td.BFiles["public/favicon.ico"] = []byte{0x00}
	td.NodeModules["svelte"] = version.Svelte
	td.NodeModules["livebud"] = "*"
	dir := t.TempDir()
	err := td.Write(dir, testdir.WithBackup(false))
	is.NoErr(err)
	err = exists(os.DirFS(dir),
		"controller/controller.go",
		"public/favicon.ico",
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
	td.Modules["github.com/livebud/bud-test-plugin"] = "v0.0.2"
	td.Files["controller/controller.go"] = `package controller`
	td.BFiles["public/favicon.ico"] = []byte{0x00}
	td.NodeModules["svelte"] = version.Svelte
	td.NodeModules["livebud"] = "*"
	dir := t.TempDir()
	err := td.Write(dir, testdir.WithBackup(true))
	is.NoErr(err)
	err = exists(os.DirFS(dir),
		"controller/controller.go",
		"public/favicon.ico",
		"node_modules/svelte/package.json",
		"node_modules/livebud/package.json",
		"package.json",
		"go.mod",
	)
	td.Modules = map[string]string{}
	delete(td.Files, "controller/controller.go")
	delete(td.BFiles, "public/favicon.ico")
	err = td.Write(dir, testdir.WithBackup(true))
	is.NoErr(err)
	is.NoErr(notExists(os.DirFS(dir),
		"controller/controller.go",
		"public/favicon.ico",
	))
	is.NoErr(exists(os.DirFS(dir),
		"node_modules/livebud/package.json",
		"node_modules/svelte/package.json",
		"package.json",
		"go.mod",
	))
}

func TestSkip(t *testing.T) {
	is := is.New(t)
	td := testdir.New()
	td.Modules["github.com/livebud/bud-test-plugin"] = "v0.0.2"
	td.Files["controller/controller.go"] = `package controller`
	td.BFiles["public/favicon.ico"] = []byte{0x00}
	td.NodeModules["svelte"] = version.Svelte
	td.NodeModules["livebud"] = "*"
	dir := t.TempDir()
	err := td.Write(dir, testdir.WithBackup(true))
	is.NoErr(err)
	err = exists(os.DirFS(dir),
		"controller/controller.go",
		"public/favicon.ico",
		"node_modules/svelte/package.json",
		"node_modules/livebud/package.json",
		"package.json",
		"go.mod",
	)
	td.Modules = map[string]string{}
	delete(td.Files, "controller/controller.go")
	delete(td.BFiles, "public/favicon.ico")
	err = td.Write(dir, testdir.WithBackup(true), testdir.WithSkip(func(name string, isDir bool) bool {
		return (name == "controller" && isDir)
	}))
	is.NoErr(err)
	is.NoErr(notExists(os.DirFS(dir),
		"public/favicon.ico",
	))
	is.NoErr(exists(os.DirFS(dir),
		"controller/controller.go",
		"node_modules/livebud/package.json",
		"node_modules/svelte/package.json",
		"package.json",
		"go.mod",
	))
}
