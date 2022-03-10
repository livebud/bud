package command_test

import (
	"context"
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/framework2/command"
	"gitlab.com/mnm/bud/internal/testdir"
	"gitlab.com/mnm/bud/package/overlay"
	"gitlab.com/mnm/bud/pkg/gomod"
	goparse "gitlab.com/mnm/bud/pkg/parser"
)

func parse(ctx context.Context, dir string) (*command.State, error) {
	module, err := gomod.Find(dir)
	if err != nil {
		return nil, err
	}
	overlay, err := overlay.Load(module)
	if err != nil {
		return nil, err
	}
	goparser := goparse.New(overlay, module)
	cmd := command.New(overlay, module, goparser)
	return cmd.Parse(ctx)
}

func TestParseEmpty(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	td := testdir.New()
	dir := t.TempDir()
	err := td.Write(dir)
	is.NoErr(err)
	state, err := parse(ctx, dir)
	is.NoErr(err)
	// Empty root command
	is.Equal(len(state.Imports), 0)
	is.Equal(state.Command, nil)
}

const rootCode = `
package command

import "context"

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

func TestParseRoot(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	td := testdir.New()
	td.Files["command/command.go"] = rootCode
	dir := t.TempDir()
	err := td.Write(dir)
	is.NoErr(err)
	state, err := parse(ctx, dir)
	is.NoErr(err)
	// Root command
	rootCmd := state.Command
	is.Equal(rootCmd.Runnable, false)
	is.Equal(rootCmd.Context, false)
	is.Equal(rootCmd.Help, "")
	is.Equal(rootCmd.Name, "")
	is.Equal(string(rootCmd.Full()), "")
	is.Equal(rootCmd.Import.Name, "command")
	is.Equal(rootCmd.Import.Path, "app.com/command")
	is.Equal(len(rootCmd.Subs), 2)
	is.Equal(len(rootCmd.Args), 0)
	is.Equal(len(rootCmd.Flags), 0)
	// Run command
	runCmd := state.Command.Subs[0]
	is.Equal(runCmd.Runnable, true)
	is.Equal(runCmd.Context, true)
	is.Equal(runCmd.Help, "")
	is.Equal(runCmd.Name, "Run")
	is.Equal(string(runCmd.Full()), "Run")
	is.Equal(runCmd.Import.Name, "command")
	is.Equal(runCmd.Import.Path, "app.com/command")
	is.Equal(len(runCmd.Subs), 0)
	// Flags
	is.Equal(len(runCmd.Flags), 2)
	embedFlag := runCmd.Flags[0]
	is.Equal(embedFlag.Name, "Embed")
	is.Equal(embedFlag.Help, "embed assets")
	is.Equal(embedFlag.Type, "*bool")
	is.Equal(embedFlag.Default, nil)
	is.Equal(embedFlag.Optional, true)
	is.True(embedFlag.Short == 0)
	hotFlag := runCmd.Flags[1]
	is.Equal(hotFlag.Name, "Hot")
	is.Equal(hotFlag.Help, "start the hot reload server")
	is.Equal(hotFlag.Type, "*bool")
	is.Equal(hotFlag.Default, nil)
	is.Equal(hotFlag.Optional, true)
	is.True(hotFlag.Short == 0)
	// Args
	is.Equal(len(runCmd.Args), 1)
	arg := runCmd.Args[0]
	is.Equal(arg.Name, "args")
	is.Equal(arg.Type, "...string")

	// Build command
	buildCmd := state.Command.Subs[1]
	is.Equal(buildCmd.Runnable, true)
	is.Equal(buildCmd.Context, true)
	is.Equal(buildCmd.Help, "")
	is.Equal(buildCmd.Name, "Build")
	is.Equal(string(buildCmd.Full()), "Build")
	is.Equal(buildCmd.Import.Name, "command")
	is.Equal(buildCmd.Import.Path, "app.com/command")
	is.Equal(len(buildCmd.Subs), 0)
	// Flags
	is.Equal(len(buildCmd.Flags), 2)
	embedFlag = buildCmd.Flags[0]
	is.Equal(embedFlag.Name, "Embed")
	is.Equal(embedFlag.Help, "embed assets")
	is.Equal(embedFlag.Type, "*bool")
	is.Equal(embedFlag.Default, nil)
	is.Equal(embedFlag.Optional, true)
	is.True(embedFlag.Short == 0)
	hotFlag = buildCmd.Flags[1]
	is.Equal(hotFlag.Name, "Hot")
	is.Equal(hotFlag.Help, "start the hot reload server")
	is.Equal(hotFlag.Type, "*bool")
	is.Equal(hotFlag.Default, nil)
	is.Equal(hotFlag.Optional, true)
	is.True(hotFlag.Short == 0)
	// Args
	is.Equal(len(buildCmd.Args), 0)
}

const migrateCode = `
package migrate

import "app.com/internal/pg"
import "context"

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

// New creates a new migration by name
func (c *Command) New(ctx context.Context, flag *NewFlag, name *string) error {
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

func TestParseMigrate(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	td := testdir.New()
	td.Files["internal/pg/pg.go"] = `
		package pg
		type Database struct {}
	`
	td.Files["command/migrate/migrate.go"] = migrateCode
	dir := t.TempDir()
	err := td.Write(dir)
	is.NoErr(err)
	state, err := parse(ctx, dir)
	is.NoErr(err)
	_ = state
	// // Root command
	// is.Equal(state.Command.Context, false)
	// is.Equal(len(state.Command.Flags), 0)
	// is.Equal(state.Command.Help, "")
	// is.Equal(state.Command.Name, "")
	// is.Equal(state.string(Command.Full()), "")
	// // is.Equal(state.Command.Path, "command")
	// is.Equal(state.Command.Import, nil)
}

func TestParseRootMigrate(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	td := testdir.New()
	td.Files["internal/pg/pg.go"] = `
		package pg
		type Database struct {}
	`
	td.Files["command/command.go"] = rootCode
	td.Files["command/migrate/migrate.go"] = migrateCode
	dir := t.TempDir()
	err := td.Write(dir)
	is.NoErr(err)
	state, err := parse(ctx, dir)
	is.NoErr(err)

	// Root command
	rootCmd := state.Command
	is.Equal(len(rootCmd.Subs), 3)
	is.Equal(len(rootCmd.Args), 0)
	is.Equal(len(rootCmd.Flags), 0)

	// Migrate command
	migrateCmd := state.Command.Subs[2]
	is.Equal(migrateCmd.Name, "migrate")
	is.Equal(migrateCmd.Runnable, false)
	is.Equal(migrateCmd.Context, false)
	is.Equal(migrateCmd.Help, "")
	is.Equal(string(migrateCmd.Full()), "migrate")
	is.Equal(migrateCmd.Import.Name, "migrate")
	is.Equal(migrateCmd.Import.Path, "app.com/command/migrate")
	is.Equal(len(migrateCmd.Flags), 0)
	is.Equal(len(migrateCmd.Subs), 3)

	// New command
	newCmd := migrateCmd.Subs[0]
	is.Equal(newCmd.Name, "New")
	is.Equal(newCmd.Runnable, true)
	is.Equal(newCmd.Context, true)
	is.Equal(newCmd.Help, "")
	is.Equal(string(newCmd.Full()), "migrate New")
	is.Equal(newCmd.Import.Name, "migrate")
	is.Equal(newCmd.Import.Path, "app.com/command/migrate")
	is.Equal(len(newCmd.Subs), 0)
	// Flags
	is.Equal(len(newCmd.Flags), 1)
	tableFlag := newCmd.Flags[0]
	is.Equal(tableFlag.Name, "Table")
	is.Equal(tableFlag.Help, "migration table")
	is.Equal(tableFlag.Type, "string")
	is.Equal(*tableFlag.Default, `"migrate"`)
	is.Equal(tableFlag.Optional, false)
	is.True(tableFlag.Short == 0)
	// Args
	is.Equal(len(newCmd.Args), 1)
	nameArg := newCmd.Args[0]
	is.Equal(nameArg.Name, "name")
	is.Equal(nameArg.Type, "*string")
	is.Equal(nameArg.Optional, true)

	// Up command
	upCmd := migrateCmd.Subs[1]
	is.Equal(upCmd.Name, "Up")
	is.Equal(upCmd.Runnable, true)
	is.Equal(upCmd.Context, true)
	is.Equal(upCmd.Help, "")
	is.Equal(string(upCmd.Full()), "migrate Up")
	is.Equal(upCmd.Import.Name, "migrate")
	is.Equal(upCmd.Import.Path, "app.com/command/migrate")
	is.Equal(len(upCmd.Subs), 0)
	// Flags
	is.Equal(len(upCmd.Flags), 2)
	dirFlag := upCmd.Flags[0]
	is.Equal(dirFlag.Name, "Dir")
	is.Equal(dirFlag.Help, "migrations directory")
	is.Equal(dirFlag.Type, "string")
	is.Equal(*dirFlag.Default, `"./migrate"`)
	is.Equal(dirFlag.Optional, false)
	is.True(dirFlag.Short == 'D')
	tableFlag = upCmd.Flags[1]
	is.Equal(tableFlag.Name, "Table")
	is.Equal(tableFlag.Help, "migration table")
	is.Equal(tableFlag.Type, "string")
	is.Equal(*tableFlag.Default, `"migrate"`)
	is.Equal(tableFlag.Optional, false)
	is.True(tableFlag.Short == 0)
	// Args
	is.Equal(len(upCmd.Args), 2)
	urlArg := upCmd.Args[0]
	is.Equal(urlArg.Name, "url")
	is.Equal(urlArg.Type, "string")
	is.Equal(urlArg.Optional, false)
	nArg := upCmd.Args[1]
	is.Equal(nArg.Name, "n")
	is.Equal(nArg.Type, "*int")
	is.Equal(nArg.Optional, true)

	// Info command
	infoCmd := migrateCmd.Subs[2]
	is.Equal(infoCmd.Name, "Info")
	is.Equal(infoCmd.Runnable, true)
	is.Equal(infoCmd.Context, true)
	is.Equal(infoCmd.Help, "")
	is.Equal(string(infoCmd.Full()), "migrate Info")
	is.Equal(infoCmd.Import.Name, "migrate")
	is.Equal(infoCmd.Import.Path, "app.com/command/migrate")
	is.Equal(len(infoCmd.Subs), 0)
	// Flags
	is.Equal(len(infoCmd.Flags), 2)
	dirFlag = infoCmd.Flags[0]
	is.Equal(dirFlag.Name, "Dir")
	is.Equal(dirFlag.Help, "migrations directory")
	is.Equal(dirFlag.Type, "string")
	is.Equal(*dirFlag.Default, `"./migrate"`)
	is.Equal(dirFlag.Optional, false)
	is.True(dirFlag.Short == 'D')
	tableFlag = infoCmd.Flags[1]
	is.Equal(tableFlag.Name, "Table")
	is.Equal(tableFlag.Help, "migration table")
	is.Equal(tableFlag.Type, "string")
	is.Equal(*tableFlag.Default, `"migrate"`)
	is.Equal(tableFlag.Optional, false)
	is.True(tableFlag.Short == 0)
	// Args
	is.Equal(len(infoCmd.Args), 1)
	urlArg = infoCmd.Args[0]
	is.Equal(urlArg.Name, "url")
	is.Equal(urlArg.Type, "string")
	is.Equal(urlArg.Optional, false)
}

func TestWebService(t *testing.T) {

}

// TODO: test deep path
