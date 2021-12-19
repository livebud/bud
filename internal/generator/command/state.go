package command

import (
	"fmt"
	"strings"

	"github.com/matthewmueller/gotext"
	"github.com/matthewmueller/text"
	"gitlab.com/mnm/bud/internal/di"
	"gitlab.com/mnm/bud/internal/imports"
)

type State struct {
	Imports []*imports.Import
	Command *Command
}

type Command struct {
	Name    string
	Usage   string
	Flags   []*Flag
	Args    []*Arg
	Subs    []*Command
	Deps    []di.Dependency
	Context bool
}

func (c *Command) Slim() string {
	return gotext.Slim(c.Name)
}

func (c *Command) Pascal() string {
	return gotext.Pascal(c.Name)
}

func (c *Command) Structs() (structs []*commandStruct) {
	stct := &commandStruct{
		Name: c.Name,
	}
	structs = append(structs, stct)
	for _, sub := range c.Subs {
		stct.Fields = append(stct.Fields, &commandStructField{
			Name: sub.Name,
		})
		structs = append(structs, sub.Structs()...)
	}
	return structs
}

type commandStruct struct {
	Name   string
	Fields []*commandStructField
}

func (c *commandStruct) Camel(suffix ...string) string {
	return gotext.Camel(append([]string{c.Name}, suffix...)...)
}

func (c *commandStruct) Pascal(suffix ...string) string {
	return gotext.Pascal(append([]string{c.Name}, suffix...)...)
}

type commandStructField struct {
	Name string
	Type string
}

func (c *commandStructField) Camel(suffix ...string) string {
	return gotext.Camel(append([]string{c.Name}, suffix...)...)
}

func (c *commandStructField) Pascal(suffix ...string) string {
	return gotext.Pascal(append([]string{c.Name}, suffix...)...)
}

type Flag struct {
	Name    string
	Usage   string
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
	Usage   string
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

type Provider = di.Provider
