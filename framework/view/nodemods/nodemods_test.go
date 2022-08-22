package nodemods_test

import (
	"context"
	"io/fs"
	"strings"
	"testing"

	"github.com/livebud/bud/framework/view/nodemods"
	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/testdir"
	"github.com/livebud/bud/internal/versions"
	"github.com/livebud/bud/package/budfs"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log/testlog"
)

func TestNodeModules(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["view/index.svelte"] = `<h1>hi world</h1>`
	td.NodeModules["svelte"] = versions.Svelte
	is.NoErr(td.Write(ctx))
	module, err := gomod.Find(dir)
	is.NoErr(err)
	bfs := budfs.New(module, log)
	bfs.DirGenerator("bud/node_modules", nodemods.New(module))
	// Read from the node_modules directory.
	code, err := fs.ReadFile(bfs, "bud/node_modules/svelte/internal")
	is.NoErr(err)
	is.True(strings.Contains(string(code), `function element(`))
	is.True(strings.Contains(string(code), `function text(`))
}

func TestLiveBud(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.NodeModules["livebud"] = "*"
	is.NoErr(td.Write(ctx))
	module, err := gomod.Find(dir)
	is.NoErr(err)
	bfs := budfs.New(module, log)
	bfs.DirGenerator("bud/node_modules", nodemods.New(module))
	// Read the next runtime file
	code, err := fs.ReadFile(bfs, "bud/node_modules/livebud/runtime/svelte")
	is.NoErr(err)
	is.True(strings.Contains(string(code), `function createView(input)`))
	// Read the first runtime file
	code, err = fs.ReadFile(bfs, "bud/node_modules/livebud/runtime/hot")
	is.NoErr(err)
	is.True(strings.Contains(string(code), `Hot = class`))
}
