package command_test

import (
	"fmt"
	"go/format"
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/framework2/command"
)

func rootState() *command.Cmd {
	return &command.Cmd{
		Name:     "",
		Runnable: false,
		Subs: []*command.Cmd{
			{
				Name:     "Run",
				Context:  true,
				Runnable: true,
				Help:     "Run runs the development server",
				Flags: []*command.Flag{
					{
						Name: "Embed",
						Type: "*bool",
						Help: "embed assets",
					},
					{
						Name: "Hot",
						Type: "*bool",
						Help: "start the hot reload server",
					},
				},
				Args: []*command.Arg{
					{
						Name: "args",
						Type: "...string",
					},
				},
			},
			{
				Name:     "Build",
				Context:  true,
				Runnable: true,
				Help:     "Build builds the production server",
				Flags: []*command.Flag{
					{
						Name: "Embed",
						Type: "*bool",
						Help: "embed assets",
					},
					{
						Name: "Hot",
						Type: "*bool",
						Help: "start the hot reload server",
					},
				},
			},
		},
	}
}

func migrateState() *command.Cmd {
	return &command.Cmd{
		Name: "migrate",
		Help: "migrate your database",
		Subs: []*command.Cmd{
			{
				Name:     "New",
				Context:  true,
				Runnable: true,
				Help:     "create a new migration by name",
				Flags: []*command.Flag{
					{
						Name:    "Dir",
						Type:    "string",
						Help:    "migrations directory",
						Default: ptrString(`"./migrate"`),
					},
					{
						Name:    "Table",
						Type:    "string",
						Help:    "migration table",
						Default: ptrString(`"migrate"`),
					},
				},
				Args: []*command.Arg{
					{
						Name: "name",
						Type: "*string",
					},
				},
			},
			{
				Name:     "Up",
				Context:  true,
				Runnable: true,
				Help:     "migrate the database at url by n migrations",
				Flags: []*command.Flag{
					{
						Name:    "Dir",
						Type:    "string",
						Help:    "migrations directory",
						Default: ptrString(`"./migrate"`),
					},
					{
						Name:    "Table",
						Type:    "string",
						Help:    "migration table",
						Default: ptrString(`"migrate"`),
					},
				},
				Args: []*command.Arg{
					{
						Name: "url",
						Type: "string",
					},
					{
						Name: "n",
						Type: "*int",
					},
				},
			},
			{
				Name:     "Info",
				Context:  true,
				Runnable: true,
				Help:     "gets information on the current migration",
				Flags: []*command.Flag{
					{
						Name:    "Dir",
						Type:    "string",
						Help:    "migrations directory",
						Default: ptrString(`"./migrate"`),
					},
					{
						Name:    "Table",
						Type:    "string",
						Help:    "migration table",
						Default: ptrString(`"migrate"`),
					},
				},
				Args: []*command.Arg{
					{
						Name: "url",
						Type: "string",
					},
				},
			},
		},
	}
}

func ptrString(s string) *string {
	return &s
}

func TestGenerateEmpty(t *testing.T) {
	is := is.New(t)
	code, err := command.Generate(&command.State{})
	is.NoErr(err)
	code, err = format.Source(code)
	is.NoErr(err)
}

func TestGenerateRoot(t *testing.T) {
	is := is.New(t)
	code, err := command.Generate(&command.State{
		Command: rootState(),
	})
	is.NoErr(err)
	code, err = format.Source(code)
	is.NoErr(err)
}

func TestGenerateMigrate(t *testing.T) {
	is := is.New(t)
	rootCmd := rootState()
	rootCmd.Subs = append(rootCmd.Subs, migrateState())
	code, err := command.Generate(&command.State{
		Command: rootCmd,
	})
	is.NoErr(err)
	code, err = format.Source(code)
	is.NoErr(err)
	fmt.Println(string(code))
}
