package fsmode

import "io/fs"

type Mode uint8

const (
	Dir Mode = 1 << iota
	Gen
)

const GenDir = Gen | Dir

func (m Mode) IsDir() bool {
	return m&Dir != 0
}

func (m Mode) IsGen() bool {
	return m&Gen != 0
}

func (m Mode) FileMode() fs.FileMode {
	mode := fs.FileMode(0)
	if m.IsDir() {
		mode |= fs.ModeDir
	}
	return mode
}

func (m Mode) String() string {
	var s string
	if m.IsDir() {
		s += "d"
	} else {
		s += "-"
	}
	if m.IsGen() {
		s += "g"
	} else {
		s += "-"
	}
	return s
}
