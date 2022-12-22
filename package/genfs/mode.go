package genfs

import "io/fs"

type mode uint8

const (
	modeDir mode = 1 << iota
	modeGen
)

func (m mode) IsDir() bool {
	return m&modeDir != 0
}

func (m mode) IsGen() bool {
	return m&modeGen != 0
}

func (m mode) FileMode() fs.FileMode {
	mode := fs.FileMode(0)
	if m.IsDir() {
		mode |= fs.ModeDir
	}
	return mode
}

func (m mode) String() string {
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
