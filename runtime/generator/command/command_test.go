package command_test

import (
	"context"
	"strings"
	"testing"

	"github.com/livebud/bud/internal/budtest"

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
	ctx := context.Background()
	dir := t.TempDir()
	bud := budtest.New(dir)
	project, err := bud.Compile(ctx)
	is.NoErr(err)
	app, err := project.Build(ctx)
	is.NoErr(err)
	is.NoErr(app.Exists("bud/.app/command/command.go"))
}

func TestCommand(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	bud := budtest.New(dir)
	bud.Files["internal/hn/hn.go"] = `
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
	bud.Files["command/deploy/deploy.go"] = `
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
	project, err := bud.Compile(ctx)
	is.NoErr(err)
	app, err := project.Build(ctx)
	is.NoErr(err)
	stdout, stderr, err := app.Execute(ctx, "-h")
	is.NoErr(err)
	is.NoErr(stderr.Expect(""))
	isEqual(t, stdout.String(), `
		Usage:
		  app [command]

		Commands:
		  deploy
	`)
}

func TestNested(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	bud := budtest.New(dir)
	bud.Files["command/deploy/deploy.go"] = `
		package deploy

		import (
			"context"
			"fmt"

			router "github.com/livebud/bud/package/router"
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
	bud.Files["command/new/new.go"] = `
		package new

		import (
			"context"
			"fmt"

			router "github.com/livebud/bud/package/router"
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
	bud.Files["command/new/view/view.go"] = `
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
	project, err := bud.Compile(ctx)
	is.NoErr(err)
	app, err := project.Build(ctx)
	is.NoErr(err)
	stdout, stderr, err := app.Execute(ctx, "-h")
	is.NoErr(err)
	is.NoErr(stderr.Expect(""))
	isEqual(t, stdout.String(), `
		Usage:
		  app [command]

		Commands:
		  deploy
		  new
	`)
	stdout, stderr, err = app.Execute(ctx, "new", "-h")
	is.NoErr(err)
	is.NoErr(stderr.Expect(""))
	isEqual(t, stdout.String(), `
		Usage:
		  new [flags] [command]

		Flags:
		  --dry-run  run but don't write

		Commands:
		  view
	`)
	stdout, stderr, err = app.Execute(ctx, "new", "view", "-h")
	is.NoErr(err)
	is.NoErr(stderr.Expect(""))
	isEqual(t, stdout.String(), `
		Usage:
		  view [flags] <name>

		Flags:
		  --with-test  include a view test
	`)
}

// TODO: test deeply nested
