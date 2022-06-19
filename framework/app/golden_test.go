package app_test

import (
	"context"
	"errors"
	"io/fs"
	"testing"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/framework/app"
	"github.com/livebud/bud/framework/controller"
	"github.com/livebud/bud/framework/web"
	"github.com/livebud/bud/internal/golden"
	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/testdir"
	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/log/testlog"
	"github.com/livebud/bud/package/parser"
)

func load(fsys fs.FS, log log.Interface, module *gomod.Module, flag *framework.Flag) (*app.State, error) {
	parser := parser.New(fsys, module)
	injector := di.New(fsys, log, module, parser)
	return app.Load(fsys, injector, module, flag)
}

func generateWeb(fsys fs.FS, module *gomod.Module) ([]byte, error) {
	parser := parser.New(fsys, module)
	state, err := web.Load(fsys, module, parser)
	if err != nil {
		return nil, err
	}
	return web.Generate(state)
}

func generateController(fsys fs.FS, log log.Interface, module *gomod.Module) ([]byte, error) {
	parser := parser.New(fsys, module)
	injector := di.New(fsys, log, module, parser)
	state, err := controller.Load(fsys, injector, module, parser)
	if err != nil {
		return nil, err
	}
	return controller.Generate(state)
}

func TestEmpty(t *testing.T) {
	is := is.New(t)
	log := testlog.Log()
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	is.NoErr(td.Write(ctx))
	module, err := gomod.Find(dir)
	is.NoErr(err)
	state, err := load(module, log, module, &framework.Flag{Embed: false})
	is.True(err != nil)
	is.True(errors.Is(err, fs.ErrNotExist))
	is.Equal(state, nil)
}

func TestWelcome(t *testing.T) {
	is := is.New(t)
	log := testlog.Log()
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	is.NoErr(td.Write(ctx))
	module, err := gomod.Find(dir)
	is.NoErr(err)
	web, err := generateWeb(module, module)
	is.NoErr(err)
	td.BFiles["bud/internal/app/web/web.go"] = web
	is.NoErr(td.Write(ctx))
	state, err := load(module, log, module, &framework.Flag{Embed: false})
	is.NoErr(err)
	code, err := app.Generate(state)
	is.NoErr(err)
	is.NoErr(parser.Check(code))
	golden.TestGenerator(t, state, code)
}

func TestControllerWeb(t *testing.T) {
	is := is.New(t)
	log := testlog.Log()
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["controller/controller.go"] = `
		package controller
		type Controller struct {}
		func (c *Controller) Index() string { return "hello" }
	`
	is.NoErr(td.Write(ctx))
	module, err := gomod.Find(dir)
	is.NoErr(err)
	controller, err := generateController(module, log, module)
	is.NoErr(err)
	td.BFiles["bud/internal/app/controller/controller.go"] = controller
	is.NoErr(td.Write(ctx))
	web, err := generateWeb(module, module)
	is.NoErr(err)
	td.BFiles["bud/internal/app/web/web.go"] = web
	is.NoErr(td.Write(ctx))
	state, err := load(module, log, module, &framework.Flag{Embed: false})
	is.NoErr(err)
	code, err := app.Generate(state)
	is.NoErr(err)
	is.NoErr(parser.Check(code))
	golden.TestGenerator(t, state, code)
}
