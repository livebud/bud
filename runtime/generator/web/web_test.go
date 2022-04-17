package web_test

import (
	"context"
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/internal/budtest"
)

// TODO: generate welcome server when there are no routes
func TestEmpty(t *testing.T) {
	t.SkipNow()
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	bud := budtest.New(dir)
	project, err := bud.Compile(ctx)
	is.NoErr(err)
	app, err := project.Build(ctx)
	is.NoErr(err)
	is.NoErr(app.Exists("bud/.app/web/web.go"))
}

func TestRootAction(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	bud := budtest.New(dir)
	bud.Files["controller/controller.go"] = `
		package controller
		type Controller struct {}
		func (c *Controller) Index() {}
		func (c *Controller) Show() {}
		func (c *Controller) New() {}
		func (c *Controller) Edit() {}
		func (c *Controller) Create() {}
		func (c *Controller) Update() {}
		func (c *Controller) Delete() {}
	`
	project, err := bud.Compile(ctx)
	is.NoErr(err)
	app, err := project.Build(ctx)
	is.NoErr(err)
	is.NoErr(app.Exists("bud/.app/web/web.go"))
	server, err := app.Start(ctx)
	is.NoErr(err)
	defer server.Close()
	res, err := server.Get("/")
	is.NoErr(err)
	is.Equal(res.StatusCode, 204)
	res, err = server.Get("/new")
	is.NoErr(err)
	is.Equal(res.StatusCode, 204)
}
