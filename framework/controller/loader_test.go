package controller_test

import (
	"context"
	"errors"
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/parser"

	"github.com/livebud/bud/framework/controller"
	"github.com/livebud/bud/internal/golden"
	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/testdir"
	"github.com/livebud/bud/package/gomod"
)

func load(fsys fs.FS, module *gomod.Module) (*controller.State, error) {
	parser := parser.New(fsys, module)
	injector := di.New(fsys, module, parser)
	return controller.Load(fsys, injector, module, parser)
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
	is.True(err != nil)
	is.True(errors.Is(err, fs.ErrNotExist))
	is.Equal(state, nil)
}

func TestEmptyControllerDir(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.MapFiles["controller"] = &fstest.MapFile{Mode: fs.ModeDir}
	is.NoErr(td.Write(ctx))
	module, err := gomod.Find(dir)
	is.NoErr(err)
	state, err := load(module, module)
	is.True(err != nil)
	is.True(errors.Is(err, fs.ErrNotExist))
	// In this case, state has already been initialized.
	is.True(state != nil)
}

func TestHelloString(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["controller/controller.go"] = `
		package controller
		type Controller struct {}
		func (c *Controller) Index() string { return "Root" }
	`
	is.NoErr(td.Write(ctx))
	module, err := gomod.Find(dir)
	is.NoErr(err)
	state, err := load(module, module)
	is.NoErr(err)
	golden.State(t, state)
}
