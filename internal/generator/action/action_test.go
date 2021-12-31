package action_test

import (
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/internal/test"
)

func TestIndexString(t *testing.T) {
	is := is.New(t)
	generator := test.Generator(t)
	generator.Files["action/action.go"] = `
		package action

		type Controller struct {
		}

		func (c *Controller) Index() string {
			return "Hello Users!"
		}
	`
	app, err := generator.Generate()
	is.NoErr(err)
	is.True(app.Exists("bud/action/action.go"))
	is.True(app.Exists("bud/main.go"))
	server, err := app.Start()
	is.NoErr(err)
	defer server.Close()
	res, err := server.Get("/")
	is.NoErr(err)
	is.NoErr(err)
	res.Expect(`
		HTTP/1.1 200 OK
		Content-Type: application/json
		Date: Fri, 31 Dec 2021 00:00:00 GMT

		"Hello Users!"
	`)
}

func TestAboutIndexString(t *testing.T) {
	is := is.New(t)
	generator := test.Generator(t)
	generator.Files["action/about/action.go"] = `
		package action
		type Controller struct {}
		func (c *Controller) Index() string { return "About" }
	`
	app, err := generator.Generate()
	is.NoErr(err)
	is.True(app.Exists("bud/action/action.go"))
	is.True(app.Exists("bud/main.go"))
	server, err := app.Start()
	is.NoErr(err)
	defer server.Close()
	res, err := server.Get("/about")
	is.NoErr(err)
	res.Expect(`
		HTTP/1.1 200 OK
		Content-Type: application/json
		Date: Fri, 31 Dec 2021 00:00:00 GMT

		"About"
	`)
}

func TestCreateEmpty(t *testing.T) {
	is := is.New(t)
	generator := test.Generator(t)
	generator.Files["action/action.go"] = `
		package action
		type Controller struct {}
		func (c *Controller) Create() {
		}
	`
	app, err := generator.Generate()
	is.NoErr(err)
	is.True(app.Exists("bud/action/action.go"))
	is.True(app.Exists("bud/main.go"))
	server, err := app.Start()
	is.NoErr(err)
	defer server.Close()
	res, err := server.Post("/", nil)
	is.NoErr(err)
	res.Expect(`
		HTTP/1.1 302 Found
		Date: Fri, 31 Dec 2021 00:00:00 GMT
		Location: /
	`)
}
