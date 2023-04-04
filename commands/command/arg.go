package command

import "flag"

type Arg struct {
	Name  string
	Help  string
	Value flag.Value
}

type Args []*Arg

func (as Args) Usage() string {
	return ""
}
