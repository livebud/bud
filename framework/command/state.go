package command

import (
	"strings"

	"github.com/livebud/bud/internal/imports"
)

type State struct {
	Imports  []*imports.Import
	Commands []*Command
}

func (s *State) Packages() []*Command {
	return packages(s.Commands)
}

func packages(cmds []*Command) (pkgs []*Command) {
	for _, cmd := range cmds {
		if cmd.Type != nil {
			pkgs = append(pkgs, cmd)
		}
		pkgs = append(pkgs, packages(cmd.Commands)...)
	}
	return pkgs
}

type Type struct {
	Path    string
	Package string
	Name    string
}

func (t *Type) String() string {
	out := new(strings.Builder)
	if strings.HasPrefix(t.Name, "*") {
		out.WriteString("*")
	}
	if t.Package != "" {
		out.WriteString(t.Package)
		out.WriteString(".")
	}
	out.WriteString(strings.TrimPrefix(t.Name, "*"))
	return out.String()
}

type Command struct {
	Name     string
	FullName string
	Pascal   string
	Desc     string

	Type *Type

	Runnable   bool
	Package    string
	Input      string
	HasContext bool

	Flags    []*Flag
	Args     []*Arg
	Commands []*Command
}

// type Commands []*Command

// type Command struct {
// 	Name       string
// 	Pascal     string
// 	Desc       string
// 	Input      *Type
// 	Flags      []*Flag
// 	Args       []*Arg
// 	Commands   Commands
// 	Runnable   bool
// 	HasContext bool
// }

// func (cmds Commands) All() Commands {
// 	var out Commands
// 	for _, cmd := range cmds {
// 		out = append(out, cmd)
// 		out = append(out, cmd.Commands.All()...)
// 	}
// 	return out
// }

type Flag struct {
	Name    string
	Pascal  string
	Type    string
	Desc    string
	Default string
}

type Arg struct {
	Name    string
	Pascal  string
	Type    string
	Desc    string
	Default string
}
