package process

import (
	"fmt"
	"strings"

	"github.com/matthewmueller/gotext"
	"github.com/matthewmueller/text"
	"gitlab.com/mnm/bud/internal/imports"
)

type State struct {
	Imports []*imports.Import
	Process *Process
}

type Process struct {
	Parents  []string
	Import   *imports.Import
	Name     string
	Slug     string
	Help     string
	Flags    []*Flag
	Args     []*Arg
	Subs     []*Process
	Deps     []*Dep
	Context  bool
	Runnable bool
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
