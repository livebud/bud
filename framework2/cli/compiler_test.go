package cli_test

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/internal/testdir"
	"gitlab.com/mnm/bud/internal/tester"
)

// TODO: show help on empty
func TestEmpty(t *testing.T) {
	t.SkipNow()
	is := is.New(t)
	ctx := context.Background()
	dir := "_tmp"
	td := testdir.New()
	err := td.Write(dir)
	is.NoErr(err)
	cli, err := tester.Compile(ctx, dir)
	is.NoErr(err)
	stdout, stderr, err := cli.Command("-h")
	is.NoErr(err)
	is.Equal(stdout, "Usage")
	is.Equal(stderr, "")
}

func TestHelp(t *testing.T) {
	t.SkipNow()
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New()
	err := td.Write(dir)
	is.NoErr(err)
	cli, err := tester.Compile(ctx, dir)
	is.NoErr(err)
	stdout, stderr, err := cli.Command("-h")
	is.NoErr(err)
	is.Equal(stderr, "")
	is.True(strings.Contains(stdout, "Usage")) // should contain Usage
	is.True(strings.Contains(stdout, "build")) // should contain build
	is.True(strings.Contains(stdout, "run"))   // should contain run
}

func TestCommandMigrate(t *testing.T) {
	t.SkipNow()
	is := is.New(t)
	ctx := context.Background()
	dir := "_tmp"
	os.RemoveAll(dir)
	td := testdir.New()
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
	err := td.Write(dir)
	is.NoErr(err)
	cli, err := tester.Compile(ctx, dir)
	is.NoErr(err)
	stdout, stderr, err := cli.Command("-h")
	is.NoErr(err)
	is.Equal(stdout, "")
	is.Equal(stderr, "")
}

func TestBuild(t *testing.T) {
	t.SkipNow()
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New()
	err := td.Write(dir)
	is.NoErr(err)
	cli, err := tester.Compile(ctx, dir)
	is.NoErr(err)
	stdout, stderr, err := cli.Command("build")
	is.NoErr(err)
	is.Equal(stderr, "")
	is.Equal(stdout, "")
	// // Run the app
	// app := testapp.New(dir)
	// stdout, stderr, err = app.Run("-h")
	// is.NoErr(err)
	// is.Equal(stderr, "")
	// is.True(strings.Contains(stdout, "Usage:")) // should contain Usage
	// is.True(strings.Contains(stdout, "app"))    // should contain app
}
