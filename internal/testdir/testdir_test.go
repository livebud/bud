package testdir_test

import (
	"context"
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/testdir"
	"github.com/livebud/bud/internal/version"
)

func TestDir(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	td := testdir.New(t.TempDir())
	td.Backup = false
	td.Modules["github.com/livebud/bud-test-plugin"] = "v0.0.2"
	td.Files["controller/controller.go"] = `package controller`
	td.BFiles["public/favicon.ico"] = []byte{0x00}
	td.NodeModules["svelte"] = version.Svelte
	td.NodeModules["livebud"] = "*"
	is.NoErr(td.Write(ctx))
	is.NoErr(td.Exists(
		"controller/controller.go",
		"public/favicon.ico",
		"node_modules/svelte/package.json",
		"node_modules/livebud/package.json",
		"package.json",
		"go.mod",
	))
}

func TestRefresh(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	td := testdir.New(t.TempDir())
	td.Backup = true
	td.Modules["github.com/livebud/bud-test-plugin"] = "v0.0.2"
	td.Files["controller/controller.go"] = `package controller`
	td.BFiles["public/favicon.ico"] = []byte{0x00}
	td.NodeModules["svelte"] = version.Svelte
	td.NodeModules["livebud"] = "*"
	is.NoErr(td.Write(ctx))
	is.NoErr(td.Exists(
		"controller/controller.go",
		"public/favicon.ico",
		"node_modules/svelte/package.json",
		"node_modules/livebud/package.json",
		"package.json",
		"go.mod",
	))
	td.Modules = map[string]string{}
	delete(td.Files, "controller/controller.go")
	delete(td.BFiles, "public/favicon.ico")
	is.NoErr(td.Write(ctx))
	is.NoErr(td.NotExists(
		"controller/controller.go",
		"public/favicon.ico",
	))
	is.NoErr(td.Exists(
		"node_modules/livebud/package.json",
		"node_modules/svelte/package.json",
		"package.json",
		"go.mod",
	))
}

func TestSkip(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	td := testdir.New(t.TempDir())
	td.Backup = true
	td.Modules["github.com/livebud/bud-test-plugin"] = "v0.0.2"
	td.Files["controller/controller.go"] = `package controller`
	td.BFiles["public/favicon.ico"] = []byte{0x00}
	td.NodeModules["svelte"] = version.Svelte
	td.NodeModules["livebud"] = "*"
	is.NoErr(td.Write(ctx))
	is.NoErr(td.Exists(
		"controller/controller.go",
		"public/favicon.ico",
		"node_modules/svelte/package.json",
		"node_modules/livebud/package.json",
		"package.json",
		"go.mod",
	))
	td.Skip = func(name string, isDir bool) bool {
		return (name == "controller" && isDir)
	}
	td.Modules = map[string]string{}
	delete(td.Files, "controller/controller.go")
	delete(td.BFiles, "public/favicon.ico")
	is.NoErr(td.Write(ctx))
	is.NoErr(td.NotExists(
		"public/favicon.ico",
	))
	is.NoErr(td.Exists(
		"controller/controller.go",
		"node_modules/livebud/package.json",
		"node_modules/svelte/package.json",
		"package.json",
		"go.mod",
	))
}
