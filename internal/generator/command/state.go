package command

import (
	"github.com/matthewmueller/gotext"
	"gitlab.com/mnm/bud/internal/di"
	"gitlab.com/mnm/bud/internal/imports"
)

type State struct {
	Imports  []*imports.Import
	Provider *di.Provider
	Command  *Command
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

type Arg struct {
	Name    string
	Usage   string
	Type    string
	Default string
}
