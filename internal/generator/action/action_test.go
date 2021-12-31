package action_test

import (
	"io"
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
	is.Equal(res.StatusCode, 200)
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	is.NoErr(err)
	is.Equal(string(body), `"Hello Users!"`)
}

func TestAboutString(t *testing.T) {
	is := is.New(t)
	generator := test.Generator(t)
	generator.Files["action/about/action.go"] = `
		package action

		type Controller struct {
		}

		func (c *Controller) Index() string {
			return "About"
		}
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
	is.Equal(res.StatusCode, 200)
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	is.NoErr(err)
	is.Equal(string(body), `"About"`)
}
