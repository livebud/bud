package commandp_test

import (
	"context"
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/framework/commandp"
	"gitlab.com/mnm/bud/internal/testdir"
	"gitlab.com/mnm/bud/package/overlay"
	"gitlab.com/mnm/bud/pkg/gomod"
	goparse "gitlab.com/mnm/bud/pkg/parser"
)

func parse(ctx context.Context, dir string) (*commandp.State, error) {
	module, err := gomod.Find(dir)
	if err != nil {
		return nil, err
	}
	overlay, err := overlay.Load(module)
	if err != nil {
		return nil, err
	}
	goparser := goparse.New(overlay, module)
	parser := commandp.New(overlay, module, goparser)
	return parser.Parse(ctx)
}

func TestEmpty(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	td := testdir.New()
	dir := t.TempDir()
	err := td.Write(dir)
	is.NoErr(err)
	state, err := parse(ctx, dir)
	is.NoErr(err)
	// Empty root command
	is.Equal(state.Command.Runnable, false)
	is.Equal(state.Command.Context, false)
	is.Equal(len(state.Command.Flags), 0)
	is.Equal(state.Command.Help, "")
	is.Equal(state.Command.Name, "")
	is.Equal(len(state.Command.Parents), 0)
	is.Equal(len(state.Command.Subs), 0)
	is.Equal(state.Command.Path, "command")
	is.Equal(state.Command.Import, nil)
}

func TestRootCommand(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	td := testdir.New()
	td.Files["command/command.go"] = `
		package command

		// Command manages the project
		type Command struct {
		}

		type Flag struct {
			Embed *bool ` + "`" + `help:"embed assets"` + "`" + `
			Hot *bool ` + "`" + `help:"start the hot reload server"` + "`" + `
		}

		// Run runs the development server
		func (c *Command) Run(ctx context.Context, flag *Flag, args ...string) error {
			return nil
		}

		// Build builds the production server
		func (c *Command) Build(ctx context.Context, flag *Flag) error {
			return nil
		}
	`
	dir := t.TempDir()
	err := td.Write(dir)
	is.NoErr(err)
	state, err := parse(ctx, dir)
	is.NoErr(err)
	// Root command
	is.Equal(state.Command.Runnable, false)
	is.Equal(state.Command.Context, false)
	is.Equal(len(state.Command.Flags), 0)
	is.Equal(state.Command.Help, "")
	is.Equal(state.Command.Name, "")
	is.Equal(len(state.Command.Parents), 0)
	is.Equal(len(state.Command.Subs), 3)
	is.Equal(state.Command.Path, "command")
	is.Equal(state.Command.Import, nil)
}

func TestMigrateCommand(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	td := testdir.New()
	td.Files["internal/pg/pg.go"] = `
		package pg
		type Database struct {}
	`
	td.Files["command/migrate/migrate.go"] = `
		package migrate

		func New(db *pg.Database) *Command {
			return &Command{db}
		}

		type Command struct {
			db *pg.Database
		}

		type Flag struct {
			Dir string ` + "`" + `default:"./migrate" help:"migrations directory"` + "`" + `
			Table string ` + "`" + `default:"migrate" help:"migration table"` + "`" + `
		}

		// New creates a new migration by name
		func (c *Command) New(ctx context.Context, flag *Flag, name *string) error {
			return nil
		}

		// Up migrates the database at url by n migrations
		func (c *Command) Up(ctx context.Context, flag *Flag, url string, n *int) error {
			return nil
		}

		// Info gets information on the current migration
		func (c *Command) Info(ctx context.Context, flag *Flag, url string) error {
			return nil
		}
	`
	dir := t.TempDir()
	err := td.Write(dir)
	is.NoErr(err)
	state, err := parse(ctx, dir)
	is.NoErr(err)
	// Root command
	is.Equal(state.Command.Context, false)
	is.Equal(len(state.Command.Flags), 0)
	is.Equal(state.Command.Help, "")
	is.Equal(state.Command.Name, "")
	is.Equal(len(state.Command.Parents), 0)
	is.Equal(state.Command.Path, "command")
	is.Equal(state.Command.Import, nil)
}

func TestWebService(t *testing.T) {

}
