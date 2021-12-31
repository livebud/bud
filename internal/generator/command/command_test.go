package command_test

import (
	"strings"
	"testing"

	"gitlab.com/mnm/bud/internal/test"

	"github.com/lithammer/dedent"
	"github.com/matryer/is"
	"github.com/matthewmueller/diff"
)

func redent(s string) string {
	return strings.TrimSpace(dedent.Dedent(s)) + "\n"
}

func isEqual(t testing.TB, actual, expect string) {
	diff.TestString(t, redent(expect), redent(actual))
}

func TestEmpty(t *testing.T) {
	is := is.New(t)
	generator := test.Generator(t)
	app, err := generator.Generate()
	is.NoErr(err)
	is.True(!app.Exists("bud/command/command.go"))
}

func TestCommand(t *testing.T) {
	is := is.New(t)
	generator := test.Generator(t)
	generator.Files["internal/hn/hn.go"] = `
		package hn

		func New() *Client {
			return &Client{"https://news.ycombinator.com"}
		}

		type Client struct {
			base string
		}

		func (c *Client) String() string {
			return c.base
		}
	`
	generator.Files["command/deploy/deploy.go"] = `
		package deploy

		import (
			"context"
			"fmt"

			"app.com/internal/hn"
		)

		type Command struct {
			HN        *hn.Client
			AccessKey string ` + "`" + `flag:"access-key" help:"aws access key"` + "`" + `
			SecretKey string ` + "`" + `flag:"secret-key" help:"aws secret key"` + "`" + `
		}

		func (c *Command) Run(ctx context.Context) error {
			fmt.Println(c.HN, c.AccessKey, c.SecretKey)
			return nil
		}
	`
	app, err := generator.Generate()
	is.NoErr(err)
	is.True(app.Exists("bud/main.go"))
	isEqual(t, app.Run("-h"), `
		Usage:
		  app [command]

		Commands:
		  deploy
	`)
}

func TestNested(t *testing.T) {
	is := is.New(t)
	generator := test.Generator(t)
	generator.Files["command/deploy/deploy.go"] = `
		package deploy

		import (
			"context"
			"fmt"

			router "gitlab.com/mnm/bud/router"
		)

		type Command struct {
			Router    *router.Router
			AccessKey string ` + "`" + `flag:"access-key" help:"aws access key"` + "`" + `
			SecretKey string ` + "`" + `flag:"secret-key" help:"aws secret key"` + "`" + `
		}

		func (c *Command) Run(ctx context.Context) error {
			fmt.Println(c.Router, c.AccessKey, c.SecretKey)
			return nil
		}
	`
	generator.Files["command/new/new.go"] = `
		package new

		import (
			"context"
			"fmt"

			router "gitlab.com/mnm/bud/router"
		)

		type Command struct {
			Router *router.Router
			DryRun bool ` + "`" + `flag:"dry-run" help:"run but don't write" default:"false"` + "`" + `
		}

		func (c *Command) Run(ctx context.Context) error {
			fmt.Println("creating new", c.DryRun)
			return nil
		}
	`
	generator.Files["command/new/view/view.go"] = `
		package view

		import (
			"context"
			"fmt"
		)

		type Command struct {
			Name     string ` + "`" + `arg:"name" help:"name of the view"` + "`" + `
			WithTest bool   ` + "`" + `flag:"with-test" help:"include a view test" default:"true"` + "`" + `
		}

		func (c *Command) Run(ctx context.Context) error {
			fmt.Println("creating new view", c.Name, c.WithTest)
			return nil
		}
	`
	app, err := generator.Generate()
	is.NoErr(err)
	is.True(app.Exists("bud/main.go"))
	isEqual(t, app.Run("-h"), `
		Usage:
		  app [command]

		Commands:
		  deploy
		  new
	`)
	isEqual(t, app.Run("new", "-h"), `
		Usage:
		  new [flags] [command]

		Flags:
		  --dry-run  run but don't write

		Commands:
		  view
	`)
	isEqual(t, app.Run("new", "view", "-h"), `
		Usage:
		  view [flags]

		Flags:
		  --with-test  include a view test
	`)
}

// TODO: test deeply nested
