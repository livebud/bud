package testdir_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/testdir"
	"github.com/livebud/bud/internal/version"
)

func TestDir(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
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
	// Ensure livebud doesn't leak into dir
	is.NoErr(td.NotExists(
		"qs",
		"url",
		"tsconfig.json",
	))
}

func TestRefresh(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
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
	favicon := []byte{0x01}
	td.BFiles["public/favicon.ico"] = favicon
	td.NodeModules["uid"] = "2.0.0"
	is.NoErr(td.Write(ctx))
	is.NoErr(td.Exists(
		"controller/controller.go",
		"public/favicon.ico",
		"node_modules/livebud/package.json",
		"node_modules/svelte/package.json",
		"node_modules/uid/package.json",
		"package.json",
		"go.mod",
	))
	fav, err := os.ReadFile(filepath.Join(dir, "public/favicon.ico"))
	is.NoErr(err)
	is.Equal(favicon, fav)
}
