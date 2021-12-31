package router_test

import (
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/internal/test"
)

func TestNoRoutes(t *testing.T) {
	is := is.New(t)
	generator := test.Generator(t)
	app, err := generator.Generate()
	is.NoErr(err)
	is.Equal(false, app.Exists("bud/router/router.go"))
	is.Equal(false, app.Exists("bud/main.go"))
}

func TestRootAction(t *testing.T) {
	is := is.New(t)
	generator := test.Generator(t)
	generator.Files["action/action.go"] = `
		package action
		type Controller struct {}
		func (c *Controller) Index() {}
		func (c *Controller) Show() {}
		func (c *Controller) New() {}
		func (c *Controller) Edit() {}
		func (c *Controller) Create() {}
		func (c *Controller) Update() {}
		func (c *Controller) Delete() {}
	`
	app, err := generator.Generate()
	is.NoErr(err)
	is.Equal(true, app.Exists("bud/router/router.go"))
	is.Equal(true, app.Exists("bud/main.go"))
	server, err := app.Start()
	is.NoErr(err)
	defer server.Close()
	res, err := server.Get("/")
	is.NoErr(err)
	is.Equal(res.StatusCode, 204)
	res, err = server.Get("/new")
	is.NoErr(err)
	is.Equal(res.StatusCode, 204)
}
