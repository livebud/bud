package commandrt

import (
	"context"
	"io"

	"github.com/livebud/bud/package/commander"
)

func New(s *Schema) *commander.CLI {
	cli := commander.New(s.Root.Name)
	return cli
}

type Schema struct {
	Dir    string
	Env    []string
	Stderr io.Writer
	Stdout io.Writer
	Stdin  io.Reader
	Root   *Command
}

type Command struct {
	Name     string
	Help     string
	Commands []*Command
	Examples []*Example
	Flags    []*Flag
	Args     []*Arg
	Run      func(context.Context) error
	Start    func(context.Context) error
}

type Flag struct {
	Name  string
	Short string
	Help  string
}

type Arg struct {
	Name  string
	Short string
	Help  string
}

type Example struct {
	Help  string
	Usage string
}
