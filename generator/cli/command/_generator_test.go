package command_test

import (
	"fmt"
	"go/format"
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/generator/command"
	"gitlab.com/mnm/bud/internal/imports"
)

func rootState() *command.State {
	return &command.State{
		Imports: []*imports.Import{
			{
				Name: "command",
				Path: "app.com/command",
			},
		},
		Command: &command.Cmd{
			Import: &imports.Import{
				Name: "command",
				Path: "app.com/command",
			},
			Subs: []*command.Cmd{
				{
					Parent: &command.Cmd{
						Import: &imports.Import{
							Name: "command",
							Path: "app.com/command",
						},
						Subs: []*command.Cmd{
							{
								Parent: nil,
								Import: &imports.Import{
									Name: "command",
									Path: "app.com/command",
								},
								Name: "Run",
								Flags: []*command.Flag{
									{
										Name:     "Embed",
										Help:     "embed assets",
										Type:     "*bool",
										Optional: true,
									},
									{
										Name:     "Hot",
										Help:     "start the hot reload server",
										Type:     "*bool",
										Optional: true,
									},
								},
								Args: []*command.Arg{{
									Name: "args",
									Type: "...string",
								}},
								Context:  true,
								Runnable: true,
							},
							{
								Parent: nil,
								Import: &imports.Import{
									Name: "command",
									Path: "app.com/command",
								},
								Name: "Build",
								Flags: []*command.Flag{
									{
										Name:     "Embed",
										Help:     "embed assets",
										Type:     "*bool",
										Optional: true,
									},
									{
										Name:     "Hot",
										Help:     "start the hot reload server",
										Type:     "*bool",
										Optional: true,
									},
								},
								Context:  true,
								Runnable: true,
							},
						},
					},
					Import: &imports.Import{
						Name: "command",
						Path: "app.com/command",
					},
					Name: "Run",
					Flags: []*command.Flag{
						{
							Name:     "Embed",
							Help:     "embed assets",
							Type:     "*bool",
							Optional: true,
						},
						{
							Name:     "Hot",
							Help:     "start the hot reload server",
							Type:     "*bool",
							Optional: true,
						},
					},
					Args: []*command.Arg{{
						Name: "args",
						Type: "...string",
					}},
					Context:  true,
					Runnable: true,
				},
				{
					Parent: &command.Cmd{
						Import: &imports.Import{
							Name: "command",
							Path: "app.com/command",
						},
						Subs: []*command.Cmd{
							{
								Parent: nil,
								Import: &imports.Import{
									Name: "command",
									Path: "app.com/command",
								},
								Name: "Run",
								Flags: []*command.Flag{
									{
										Name:     "Embed",
										Help:     "embed assets",
										Type:     "*bool",
										Optional: true,
									},
									{
										Name:     "Hot",
										Help:     "start the hot reload server",
										Type:     "*bool",
										Optional: true,
									},
								},
								Args: []*command.Arg{{
									Name: "args",
									Type: "...string",
								}},
								Context:  true,
								Runnable: true,
							},
							{
								Parent: nil,
								Import: &imports.Import{
									Name: "command",
									Path: "app.com/command",
								},
								Name: "Build",
								Flags: []*command.Flag{
									{
										Name:     "Embed",
										Help:     "embed assets",
										Type:     "*bool",
										Optional: true,
									},
									{
										Name:     "Hot",
										Help:     "start the hot reload server",
										Type:     "*bool",
										Optional: true,
									},
								},
								Context:  true,
								Runnable: true,
							},
						},
					},
					Import: &imports.Import{
						Name: "command",
						Path: "app.com/command",
					},
					Name: "Build",
					Flags: []*command.Flag{
						{
							Name:     "Embed",
							Help:     "embed assets",
							Type:     "*bool",
							Optional: true,
						},
						{
							Name:     "Hot",
							Help:     "start the hot reload server",
							Type:     "*bool",
							Optional: true,
						},
					},
					Context:  true,
					Runnable: true,
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
	is.Equal(err.Error(), `command: generator must have a root command`)
	is.Equal(code, nil)
}

func TestGenerateRoot(t *testing.T) {
	is := is.New(t)
	code, err := command.Generate(rootState())
	is.NoErr(err)
	fmt.Println(string(code))
	_, err = format.Source(code)
	is.NoErr(err)
}

func TestGenerateMigrate(t *testing.T) {
	is := is.New(t)
	rootState := rootState()
	rootState.Command.Subs = append(rootState.Command.Subs, migrateState())
	code, err := command.Generate(rootState)
	is.NoErr(err)
	_, err = format.Source(code)
	is.NoErr(err)
}
