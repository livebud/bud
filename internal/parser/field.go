package parser

// Fielder is an interface containing the common operations on fields
type Fielder interface {
	File() *File
	Name() string
	Type() Type
}

func fieldString(f Fielder) string {
	return f.Name() + " " + f.Type().String()
}
