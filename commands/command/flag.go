package command

import "flag"

type Flag struct {
	Name  string
	Short byte
	Help  string
	Value flag.Value
}

type Flags []*Flag

func (fs Flags) Usage() string {
	return ""
}
