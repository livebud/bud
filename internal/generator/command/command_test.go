package command_test

import (
	"context"
	"testing"

	"github.com/livebud/bud/internal/cli"
	"github.com/livebud/bud/internal/cli/testcli"
	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/testdir"
)

func TestCommandMigrate(t *testing.T) {
	t.SkipNow()
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["internal/pg/pg.go"] = `
		package pg
		type Database struct {}
	`
	td.Files["command/migrate/migrate.go"] = `
		package migrate

		import "app.com/internal/pg"
		import "context"
		import "fmt"

		func New(db *pg.Database) *Command {
			return &Command{db}
		}

		type Command struct {
			db *pg.Database
		}

		type Flag struct {
			Dir string ` + "`" + `default:"./migrate" short:"D" help:"migrations directory"` + "`" + `
			Table string ` + "`" + `default:"migrate" help:"migration table"` + "`" + `
		}

		type NewFlag struct {
			Table string ` + "`" + `default:"migrate" help:"migration table"` + "`" + `
		}

		func assert(check bool, message string) {
			if !check {
				fmt.Fprintf(os.Stderr, message)
			}
		}

		// New creates a new migration by name
		func (c *Command) New(ctx context.Context, flag *NewFlag, name string) error {
			assert(ctx != nil, "context is not null")
			assert(flag.Table != "", "table is not empty")
			assert(name != "", "name is not empty")
			return nil
		}

		// Up migrates the database at url by n migrations
		func (c *Command) Up(ctx context.Context, flag *Flag, url string, n int) error {
			return nil
		}

		// Info gets information on the current migration
		func (c *Command) Info(flag *Flag, url string) error {
			return nil
		}
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(cli.New(dir))
	stdout, stderr, err := cli.Run(ctx, "-h")
	is.NoErr(err)
	is.Equal(stdout.String(), "")
	is.Equal(stderr.String(), "TODO")
}
