package command

import (
	"fmt"
	"strings"

	"github.com/matthewmueller/gotext"
	"github.com/matthewmueller/text"
	"gitlab.com/mnm/bud/internal/imports"
)

func methodName(dataType string) (string, error) {
	switch strings.TrimLeft(dataType, "*") {
	case "bool":
		return "Bool", nil
	case "string":
		return "String", nil
	default:
		return "", fmt.Errorf("command: unhandled type for method %q", dataType)
	}
}

// Flatten out the commands
func flatten(commands []*Command) (results []*Command) {
	for _, cmd := range commands {
		results = append(results, cmd)
		results = append(results, flatten(cmd.Subs)...)
	}
	return results
}

type State struct {
	Imports []*imports.Import
	// Functions []*Function
	// Structs   []*Struct
	Command *Command
}

// Flatten out the commands, intentionally ignoring the root command
// because that is custom generated.
func (s *State) Commands() []*Command {
	return flatten(s.Command.Subs)
}

type Command struct {
	Parents  []string
	Import   *imports.Import
	Name     string
	Slug     string
	Help     string
	Flags    []*Flag
	Args     []*Arg
	Subs     []*Command
	Deps     []*Dep
	Context  bool
	Runnable bool
}

func (c *Command) Pascal() string {
	return gotext.Pascal(c.Name)
}

func (c *Command) Full() Full {
	// Make a copy
	parents := make([]string, len(c.Parents))
	for i, parent := range c.Parents {
		parents[i] = parent
	}
	return Full(strings.Join(append(parents, c.Name), " "))
}

type Full string

func (f Full) Pascal() string {
	return gotext.Pascal(string(f))
}

type Flag struct {
	Name    string
	Help    string
	Type    string
	Default string
	Short   byte
}

func (f *Flag) Pascal() string {
	return gotext.Pascal(f.Name)
}

func (f *Flag) Slug() string {
	return text.Slug(f.Name)
}

func (f *Flag) Method() (string, error) {
	return methodName(f.Type)
}

type Arg struct {
	Name    string
	Help    string
	Type    string
	Default string
}

func (a *Arg) Pascal() string {
	return gotext.Pascal(a.Name)
}

func (a *Arg) Slug() string {
	return text.Slug(a.Name)
}

func (a *Arg) Method() (string, error) {
	return methodName(a.Type)
}

type Dep struct {
	Import *imports.Import
	Name   string
	Type   string
}

func (d *Dep) Camel() string {
	return gotext.Camel(d.Type)
}
