package parser

// builtin declaration
type builtin string

// Name is the built-in type
func (b builtin) Name() string {
	return string(b)
}

func (b builtin) Kind() Kind {
	return KindBuiltin
}

// Directory for builtin is blank
func (b builtin) Directory() string {
	return ""
}

// Package for builtin is blank
// TODO: there should probably be a built-in package to avoid panics
func (b builtin) Package() *Package {
	return nil
}
