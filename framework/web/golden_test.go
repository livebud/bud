package web_test

import (
	"context"
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/livebud/bud/framework/web"
	"github.com/livebud/bud/internal/golden"
	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/testdir"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/parser"
)

func load(fsys fs.FS, module *gomod.Module) (*web.State, error) {
	parser := parser.New(fsys, module)
	return web.Load(fsys, module, parser)
}

func TestEmpty(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	is.NoErr(td.Write(ctx))
	module, err := gomod.Find(dir)
	is.NoErr(err)
	state, err := load(module, module)
	is.NoErr(err)
	golden.State(t, state)
	code, err := web.Generate(state)
	is.NoErr(err)
	is.NoErr(parser.Check(code))
	golden.Code(t, code)
}

func TestEmptyDirs(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.MapFiles["view"] = &fstest.MapFile{Mode: fs.ModeDir}
	td.MapFiles["controller"] = &fstest.MapFile{Mode: fs.ModeDir}
	td.MapFiles["public"] = &fstest.MapFile{Mode: fs.ModeDir}
	is.NoErr(td.Write(ctx))
	module, err := gomod.Find(dir)
	is.NoErr(err)
	state, err := load(module, module)
	is.NoErr(err)
	golden.State(t, state)
	code, err := web.Generate(state)
	is.NoErr(err)
	is.NoErr(parser.Check(code))
	golden.Code(t, code)
}
