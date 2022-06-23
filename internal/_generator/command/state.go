package command

import (
	"fmt"
	"strings"

	"github.com/livebud/bud/internal/imports"
	"github.com/livebud/bud/package/di"
	"github.com/matthewmueller/gotext"
	"github.com/matthewmueller/text"
)

type State struct {
	Imports  []*imports.Import
	Command  *Cmd
	Provider *di.Provider
}

// Flatten out the commands, intentionally ignoring the root command
// because that is custom generated.
func (s *State) Commands() []*Cmd {
	return flatten(s.Command.Subs)
}

type Cmd struct {
	Parent   *Cmd
	Import   *imports.Import
	Name     string
	Help     string
	Flags    []*Flag
	Args     []*Arg
	Subs     []*Cmd
	Context  bool
	Runnable bool
}

func (c *Cmd) Pascal() string {
	return gotext.Pascal(c.Name)
}

func (c *Cmd) Slug() string {
	return text.Slug(c.Name)
}

func (c *Cmd) Full() Full {
	// Make a copy
	if c.Parent == nil {
		return Full(c.Name)
	}
	return Full(strings.TrimSpace(string(c.Parent.Full()) + " " + c.Name))
}

type Full string

func (f Full) Pascal() string {
	return gotext.Pascal(string(f))
}

func (f Full) Camel() string {
	return gotext.Camel(string(f))
}

type Flag struct {
	Name     string
	Help     string
	Type     string
	Default  *string
	Optional bool
	Short    byte
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
	Name     string
	Type     string
	Optional bool
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

func methodName(dataType string) (string, error) {
	switch strings.TrimLeft(dataType, "*") {
	case "bool":
		return "Bool", nil
	case "string":
		return "String", nil
	case "...string":
		return "Strings", nil
	case "int":
		return "Int", nil
	default:
		return "", fmt.Errorf("command: unhandled type for method %q", dataType)
	}
}

// Flatten out the commands
func flatten(commands []*Cmd) (results []*Cmd) {
	for _, cmd := range commands {
		results = append(results, cmd)
		results = append(results, flatten(cmd.Subs)...)
	}
	return results
}
